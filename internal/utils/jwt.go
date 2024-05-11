package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

const secretKey = "bomboclat"

func GenerateJWT(email, userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  email,
		"userID": userID,
		"exp":    time.Now().Add(time.Hour * 72).Unix(),
	})

	return token.SignedString([]byte(secretKey))
}
