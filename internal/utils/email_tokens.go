package utils

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"gopkg.in/gomail.v2"
	"time"
)

const (
	smtpHost = "smtp.gmail.com"
	smtpPort = 587
	smtpUser = "neqres@gmail.com"
	smtpPass = "Azamat1212403"
)

func GenerateEmailConfirmationToken(email, userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  email,
		"userID": userID,
		"action": "verifyEmail",
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	})
	return token.SignedString([]byte(secretKey))
}

func VerifyEmailToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func SendEmailConfirmation(userEmail, userID string) error {
	token, err := GenerateEmailConfirmationToken(userEmail, userID)
	if err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", smtpUser)
	m.SetHeader("To", userEmail)
	m.SetHeader("Subject", "Подтверждение вашей почты")
	m.SetBody("text/html", fmt.Sprintf("Пожалуйста, подтвердите вашу почту, перейдя по <a href=\"http://localhost:8080/confirm_email?token=%s\">ссылке</a>.", token))

	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)

	// Отправка письма
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
