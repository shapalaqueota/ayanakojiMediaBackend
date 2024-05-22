package api

import (
	"backend/internal/db"
	"backend/internal/service"
	"backend/internal/utils"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ConfirmEmailHandler(c *gin.Context) {
	token := c.Query("token")

	tokenClaims, err := utils.VerifyEmailToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	userID, ok := tokenClaims["userID"].(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to extract userID from token"})
		return
	}

	conn, err := db.DB.Acquire(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection error"})
		return
	}
	defer conn.Release()

	if err := service.ConfirmUserEmail(conn, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email confirmed successfully"})
}

func SendConfirmationEmailHandler(c *gin.Context) {
	userEmail := c.Param("email")
	userID := c.GetString("userID")

	// Вызываем функцию отправки письма
	if err := utils.SendEmailConfirmation(userEmail, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send confirmation email", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Confirmation email sent successfully"})
}

func GenerateTokenHandler(c *gin.Context) {
	userID := c.Param("userID")
	email := c.Param("email")

	token, err := utils.GenerateEmailConfirmationToken(email, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
