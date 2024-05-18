package models

type Film struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	IsSeries    bool   `json:"is_series"`
	S3Key       string `json:"s3_key"`
}

type Episode struct {
	ID            int    `json:"id"`
	FilmID        int    `json:"film_id"`
	SeasonNumber  int    `json:"season_number"`
	EpisodeNumber int    `json:"episode_number"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	S3Key         string `json:"s3_key"`
}
