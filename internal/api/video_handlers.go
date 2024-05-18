package api

import (
	"backend/internal/db"
	"backend/internal/models"
	"backend/internal/service"
	"backend/internal/utils"
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
)

func GetFilmDetails(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid film ID"})
		return
	}

	conn, err := db.DB.Acquire(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to acquire database connection"})
		return
	}
	defer conn.Release()

	film, err := service.GetFilmDetails(conn, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, film)
}

func GetFilmContentURL(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid film ID"})
		return
	}

	conn, err := db.DB.Acquire(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to acquire database connection"})
		return
	}
	defer conn.Release()

	url, err := service.GetFilmContentURL(conn, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": url})
}

func GetEpisodeContentURL(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid film ID"})
		return
	}

	season, err := strconv.Atoi(c.Query("season"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid season number"})
		return
	}

	episode, err := strconv.Atoi(c.Query("episode"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid episode number"})
		return
	}

	conn, err := db.DB.Acquire(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to acquire database connection"})
		return
	}
	defer conn.Release()

	url, err := service.GetEpisodeContentURL(conn, id, season, episode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": url})
}

func UploadFilm(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	fileContent, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer func(fileContent multipart.File) {
		err := fileContent.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(fileContent)

	fileBuffer := make([]byte, file.Size)
	if _, err := fileContent.Read(fileBuffer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}

	title := file.Filename
	s3Key := utils.GenerateS3Key(title)
	log.Printf("Uploading file to S3 with key: %s", s3Key)
	_, err = utils.UploadFile(s3Key, fileBuffer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
		return
	}

	conn, err := db.DB.Acquire(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to acquire database connection"})
		return
	}
	defer conn.Release()

	film := models.Film{
		Title:       title,
		Description: "Description based on the title or some default text.",
		IsSeries:    false, // Adjust this based on your needs or another form field
		S3Key:       s3Key,
	}

	if err := service.CreateFilm(conn, &film); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, film)
}

func UploadEpisode(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	filmID, err := strconv.Atoi(c.PostForm("film_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid film ID"})
		return
	}

	seasonNumber, err := strconv.Atoi(c.PostForm("season_number"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid season number"})
		return
	}

	episodeNumber, err := strconv.Atoi(c.PostForm("episode_number"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid episode number"})
		return
	}

	fileContent, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer fileContent.Close()

	fileBuffer := make([]byte, file.Size)
	if _, err := fileContent.Read(fileBuffer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}

	title := file.Filename
	s3Key := utils.GenerateS3Key(title)
	log.Printf("Uploading file to S3 with key: %s", s3Key)
	_, err = utils.UploadFile(s3Key, fileBuffer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
		return
	}

	conn, err := db.DB.Acquire(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to acquire database connection"})
		return
	}
	defer conn.Release()

	episode := models.Episode{
		FilmID:        filmID,
		SeasonNumber:  seasonNumber,
		EpisodeNumber: episodeNumber,
		Title:         title,
		Description:   "Description based on the title or some default text.",
		S3Key:         s3Key,
	}

	if err := service.CreateEpisode(conn, &episode); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, episode)
}
