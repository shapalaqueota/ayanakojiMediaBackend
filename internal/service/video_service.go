package service

import (
	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/utils"
	"github.com/jackc/pgx/v4/pgxpool"
)

func GetFilmDetails(conn *pgxpool.Conn, id int) (*models.Film, error) {
	return repository.GetFilmByID(conn, id)
}

func GetFilmContentURL(conn *pgxpool.Conn, id int) (string, error) {
	film, err := repository.GetFilmByID(conn, id)
	if err != nil {
		return "", err
	}

	url, err := utils.GetPresignedURL(film.S3Key)
	if err != nil {
		return "", err
	}

	return url, nil
}

func GetEpisodeContentURL(conn *pgxpool.Conn, filmID, season, episode int) (string, error) {
	ep, err := repository.GetEpisode(conn, filmID, season, episode)
	if err != nil {
		return "", err
	}

	url, err := utils.GetPresignedURL(ep.S3Key)
	if err != nil {
		return "", err
	}

	return url, nil
}

func CreateFilm(conn *pgxpool.Conn, film *models.Film) error {
	return repository.CreateFilm(conn, film)
}

func CreateEpisode(conn *pgxpool.Conn, episode *models.Episode) error {
	return repository.CreateEpisode(conn, episode)
}
