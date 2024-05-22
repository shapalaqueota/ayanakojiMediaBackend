package api

import (
	"backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func Router(router *gin.Engine) {
	router.POST("/signup", Signup)
	router.POST("/login", Login)

	router.POST("/send-confirmation-email", SendConfirmationEmailHandler)
	router.GET("/search", SearchFilmsHandler)

	authorized := router.Group("/")
	authorized.Use(middleware.TokenAuthMiddleware())
	{
		authorized.PUT("/users/update", ManageAccount)
		authorized.GET("/users/:id", GetUser)
		authorized.GET("/film/:id", GetFilmDetails)
		authorized.GET("/film/:id/content", GetFilmContentURL)
		authorized.GET("/film/:id/episode", GetEpisodeContentURL)
		authorized.GET("/films", GetAllFilms)
		authorized.POST("/upload/film", UploadFilm)
		authorized.POST("/upload/episode", UploadEpisode)
		authorized.GET("/confirm_email", ConfirmEmailHandler)
		authorized.GET("/generate_token/:userID/:email", GenerateTokenHandler)
	}
}
