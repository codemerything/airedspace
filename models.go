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
	Slug           string `json:"slug" db:"slug"`
	Director       string `json:"director" db:"job"`
}

type Persons struct {
	ID      int64  `json:"id" db:"id"`
	TMDB_ID int    `json:"tmdb_id" db:"tmdb_id"`
	Name    string `json:"name" db:"name"`
}

type Genres struct {
	Id   int64  `json:"id" db:"db"`
	Name string `json:"name" db:"name"`
}

type FilmsGenre struct {
	FilmID  int64 `json:"film_id" db:"film_id"`
	GenreID int64 `json:"genre_id" db:"genre_id"`
}

type Languages struct {
	ID       int64  `json:"id" db:"id"`
	Language string `json:"language" db:"language"`
}

type FilmLanguage struct {
	FilmID     int `json:"film_id" db:"film_id"`
	LanguageID int `json:"language_id" db:"language_id"`
}

type FilmsCast struct {
	FilmID        int64  `json:"film_id" db:"film_id"`
	CastID        int64  `json:"cast_id" db:"cast_id"`
	CharacterName string `json:"character_name" db:"character_name"`
}

type Review struct {
	ID          string  `json:"id" db:"id"`
	FilmID      int64   `json:"film_id" db:"film_id"`
	UserID      int64   `json:"user_id" db:"user_id"`
	AudioURL    string  `json:"audio_url" db:"audio_url"`
	Stars       float64 `json:"stars" db:"stars"`
	TMDB_ID     string  `json:"tmdb_id" db:"tmdb_id"`
	IsRewatch   int64   `json:"is_rewatch" db:"is_rewatch"`
	ReviewLikes int64   `json:"review_likes" db:"review_likes"`
}
