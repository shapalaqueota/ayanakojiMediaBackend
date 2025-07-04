package main

import (
	"backend/config"
	"backend/internal/api"
	"backend/internal/db"
	"backend/internal/middleware"
	"backend/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	db.ConnectDB()
	defer db.DB.Close()

	utils.InitVKCloudService()

	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORSMiddleware())
	api.Router(router)
	config.LoadConfig()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}

	err := router.Run(":" + port)
	if err != nil {
		log.Fatal("Error starting server", err)
	}
}
