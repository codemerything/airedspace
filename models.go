package main

import "time"

type User struct {
	UserID    int64     `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Name      string    `json:"name" db:"name"`
	LastName  string    `json:"lastname" db:"lastname"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"password" db:"password"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Films struct {
	FilmID         int64  `json:"id" db:"id"`
	Title          string `json:"title" db:"title"`
	TMDB_ID        int    `json:"tmdb_id" db:"tmdb_id"`
	Year           string `json:"release_year" db:"release_year"`
	Poster         string `json:"poster" db:"poster"`
	BackdropPoster string `json:"foreground_poster" db:"foreground_poster"`
	Description    string `json:"description" db:"description"`
	TagLine        string `json:"tag_line" db:"tag_line"`
	Time           int    `json:"runtime" db:"runtime"`
}

type Review struct {
	ID        int64   `json:"id" db:"id"`
	FilmID    int64   `json:"film_id" db:"film_id"`
	UserID    int64   `json:"user_id" db:"user_id"`
	AudioURL  string  `json:"audio_url" db:"audio_url"`
	Stars     float64 `json:"stars" db:"stars"`
	TMDB_ID   string  `json:"tmdb_id" db:"tmdb_id"`
	IsRewatch int64   `json:"is_rewatch" db:"is_rewatch"`
}
