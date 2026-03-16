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
	FilmID           int64  `json:"id" db:"id"`
	Title            string `json:"title" db:"title"`
	TMDB_ID          string `json:"tmdb_id" db:"tmdb_id"`
	Year             string `json:"year" db:"year"`
	Poster           string `json:"image" db:"image"`
	ForegroundPoster string `json:"foreground_poster" db:"foreground_poster"`
	Description      string `json:"desc" db:"desc"`
	TagLine          string `json:"tag_line" db:"tag_line"`
	Time             string `json:"movie_time" db:"movie_time"`
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
