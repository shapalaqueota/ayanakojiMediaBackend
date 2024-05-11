package api

import (
	"github.com/gin-gonic/gin"
)

func Router(router *gin.Engine) {
	router.POST("/users", Signup)
	router.GET("/users/:id", GetUser)
	router.POST("/login", Login)
	//router.GET("/", func(c *gin.Context) {})
	//router.GET("/movies", getMovies)
	//router.GET("/search", searchMovies)
	//router.POST("/subscriptions", createSubscription)
	//router.GET("/subscriptions/:userId", getSubscription)
	//router.GET("/users/:userId", getUser)
	//router.PUT("/users/:userId", updateUser)
}
