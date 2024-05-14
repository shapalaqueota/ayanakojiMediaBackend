package utils

import (
	"crypto/sha1"
	"fmt"
	"github.com/google/uuid"
	"io"
	"time"
)

func GenerateID(title string) string {
	hash := sha1.New()
	io.WriteString(hash, title)
	io.WriteString(hash, time.Now().String())
	uuidPart := uuid.New().String()
	return fmt.Sprintf("%x-%s", hash.Sum(nil), uuidPart)
}
