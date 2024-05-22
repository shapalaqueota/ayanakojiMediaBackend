package utils

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"os"
)

var s3Client *s3.S3
var bucket string

func InitVKCloudService() {
	region := os.Getenv("VK_CLOUD_REGION")
	accessKey := os.Getenv("VK_CLOUD_ACCESS_KEY")
	secretKey := os.Getenv("VK_CLOUD_SECRET_KEY")
	bucket = os.Getenv("VK_CLOUD_BUCKET")

	log.Printf("Initializing VK Cloud Service with region: %s, bucket: %s", region, bucket)
	log.Printf("Access Key: %s", accessKey)
	log.Printf("Secret Key: %s", secretKey)

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
		Credentials: credentials.NewStaticCredentials(
			accessKey, secretKey, ""),
		Endpoint:         aws.String("https://hb.kz-ast.vkcs.cloud"),
		S3ForcePathStyle: aws.Bool(true), // Required for VK Cloud
	})
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
	}

	s3Client = s3.New(sess)

	// Test DNS resolution
	testDNSResolution(bucket, region)
}

func testDNSResolution(bucket, region string) {
	host := "hb.kz-ast.vkcs.cloud"
	resp, err := http.Get(fmt.Sprintf("https://%s", host))
	if err != nil {
		log.Printf("DNS resolution failed for host %s: %v", host, err)
	} else {
		log.Printf("DNS resolution successful for host %s, status: %s", host, resp.Status)
	}
}

func GenerateS3Key(title string) string {
	ext := filepath.Ext(title)          // Получаем расширение файла
	name := title[:len(title)-len(ext)] // Убираем расширение для генерации уникального имени
	h := sha1.New()
	h.Write([]byte(name + time.Now().String()))
	return fmt.Sprintf("%x%s", h.Sum(nil), ext) // Добавляем обратно расширение
}

func UploadFile(key string, file []byte) (string, error) {
	log.Printf("Uploading file with key: %s", key)
	log.Printf("Bucket: %s", bucket)
	log.Printf("File size: %d bytes", len(file))
	_, err := s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(file),
		ACL:    aws.String("public-read"),
	})
	if err != nil {
		log.Printf("Failed to upload file to bucket: %v", err)
		return "", fmt.Errorf("failed to upload file: %v", err)
	}
	return key, nil
}

func GetPresignedURL(key string) (string, error) {
	req, _ := s3Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	urlStr, err := req.Presign(15 * time.Minute)
	if err != nil {
		return "", fmt.Errorf("Failed to sign request: %v", err)
	}

	return urlStr, nil
}
