package api

import (
	"backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func Router(router *gin.Engine) {
	router.POST("/signup", Signup)
	router.GET("/users/:id", GetUser)
	router.POST("/login", Login)
	router.PUT("/users/update", middleware.TokenAuthMiddleware(), ManageAccount)
	router.GET("/confirm_email", ConfirmEmailHandler)
	router.POST("/send-confirmation-email", SendConfirmationEmailHandler)

}
