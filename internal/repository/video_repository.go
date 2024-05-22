package repository

import (
	"backend/internal/db"
	"backend/internal/models"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

func GetFilmByID(conn *pgxpool.Conn, id int) (*models.Film, error) {
	var film models.Film
	err := conn.QueryRow(context.Background(), "SELECT id, title, description, is_series, s3_key FROM movies WHERE id=$1", id).Scan(
		&film.ID, &film.Title, &film.Description, &film.IsSeries, &film.S3Key,
	)
	if err != nil {
		return nil, err
	}
	return &film, nil
}

func GetEpisode(conn *pgxpool.Conn, filmID, season, episode int) (*models.Episode, error) {
	var ep models.Episode
	err := conn.QueryRow(context.Background(), "SELECT id, film_id, season_number, episode_number, title, description, s3_key FROM episodes WHERE film_id=$1 AND season_number=$2 AND episode_number=$3",
		filmID, season, episode).Scan(
		&ep.ID, &ep.FilmID, &ep.SeasonNumber, &ep.EpisodeNumber, &ep.Title, &ep.Description, &ep.S3Key,
	)
	if err != nil {
		return nil, err
	}
	return &ep, nil
}

func CreateFilm(conn *pgxpool.Conn, film *models.Film) error {
	_, err := conn.Exec(context.Background(), "INSERT INTO movies (title, description, s3_key, is_series) VALUES ($1, $2, $3, $4)",
		film.Title, film.Description, film.S3Key, film.IsSeries)
	return err
}

func CreateEpisode(conn *pgxpool.Conn, episode *models.Episode) error {
	_, err := conn.Exec(context.Background(), "INSERT INTO episodes (film_id, season_number, episode_number, title, description, s3_key) VALUES ($1, $2, $3, $4, $5, $6)",
		episode.FilmID, episode.SeasonNumber, episode.EpisodeNumber, episode.Title, episode.Description, episode.S3Key)
	return err
}

func SearchFilms(conn *pgxpool.Conn, query string) ([]models.Film, error) {
	rows, err := conn.Query(context.Background(), "SELECT id, title, description, is_series, s3_key FROM movies WHERE title ILIKE $1", "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var films []models.Film
	for rows.Next() {
		var film models.Film
		err := rows.Scan(&film.ID, &film.Title, &film.Description, &film.IsSeries, &film.S3Key)
		if err != nil {
			return nil, err
		}
		films = append(films, film)
	}

	return films, nil
}

func GetAllFilms() ([]models.Film, error) {
	var films []models.Film
	rows, err := db.DB.Query(context.Background(), "SELECT id, title, description, is_series, s3_key FROM movies")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var film models.Film
		err := rows.Scan(&film.ID, &film.Title, &film.Description, &film.IsSeries, &film.S3Key)
		if err != nil {
			return nil, err
		}
		films = append(films, film)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return films, nil
}
