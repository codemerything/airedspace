// handles storage of data to database explicitly

package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type Repository struct {
	db *sql.DB
}

type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

func (r *Repository) CreateUser(ctx context.Context, u *User) error {
	query := "INSERT INTO users (username,name,lastname, email,password,created_at, updated_at) VALUES (?,?,?,?,?,?,?)"

	res, err := r.db.ExecContext(ctx, query, u.Username, u.Name, u.LastName, u.Email, u.Password, u.CreatedAt, u.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	u.UserID, err = res.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get userid: %w", err)
	}
	return nil
}

func (r *Repository) FetchUserByUsername(u *User) (User, error) {
	var user User
	query := "SELECT username, password, email FROM users WHERE username = ?"

	err := r.db.QueryRow(query, u.Username).Scan(&user.Username, &user.Password, &user.Email)

	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, errors.New("User not found")
		}
		return User{}, err
	}

	return user, err
}

func (r *Repository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, nil)
}

func (r *Repository) CreateFilm(ctx context.Context, db DBTX, f *Films) (*Films, error) {
	query := "INSERT INTO films (tmdb_id, title, tag_line, poster, foreground_poster, description, release_year, runtime) VALUES (?,?,?,?,?,?,?,?)"

	res, err := db.ExecContext(ctx, query, f.TMDB_ID, f.Title, f.TagLine, f.Poster, f.BackdropPoster, f.Description, f.Year, f.Time)
	if err != nil {
		return nil, fmt.Errorf("failed to insert a film: %w", err)
	}

	f.FilmID, err = res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get filmid: %w", err)
	}

	return f, nil

}

func (r *Repository) CreatePerson(ctx context.Context, db DBTX, p *Persons) (*Persons, error) {
	query := "INSERT INTO persons (id,tmdb_id, name) VALUES (?,?,?) ON DUPLICATE KEY UPDATE tmdb_id = tmdb_id, id = LAST_INSERT_ID(id) "

	res, err := db.ExecContext(ctx, query, p.ID, p.TMDB_ID, p.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to insert person to table: %w", err)
	}
	p.ID, _ = res.LastInsertId()

	return p, nil
}

func (r *Repository) CreateFilmCast(ctx context.Context, db DBTX, fc *FilmsCast) error {
	query := "INSERT INTO films_cast (film_id, cast_id, character_name) VALUES (?,?,?)"

	_, err := db.ExecContext(ctx, query, fc.FilmID, fc.CastID, fc.CharacterName)
	if err != nil {
		return fmt.Errorf("failed to insert into table: %w", err)
	}
	return nil
}

func (r *Repository) CreateGenre(ctx context.Context, db DBTX, g *Genres) (*Genres, error) {
	query := "INSERT INTO genres (name) VALUES ?"

	res, err := db.ExecContext(ctx, query, g.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to insert genre into table: %w", err)
	}

	g.Id, _ = res.LastInsertId()
	return g, nil
}

func (r *Repository) CreateFilmGenre(ctx context.Context, db DBTX, g *FilmsGenre) error {
	query := "INSERT INTO films_genre (film_id, genre_id) VALUES (?,?)"

	_, err := db.ExecContext(ctx, query, g.FilmID, g.GenreID)
	if err != nil {
		return fmt.Errorf("failed to insert into films genre table: %w", err)
	}
	return nil
}

func (r *Repository) SearchFilm(ctx context.Context, f *Films) ([]Films, error) {

	var films []Films
	query := `SELECT  tmdb_id, title, tag_line, poster, foreground_poster, description, release_year, runtime FROM films WHERE title LIKE ?`

	rows, err := r.db.QueryContext(ctx, query, (f.Title + "%"))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch row: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var tempFilms Films
		if err := rows.Scan(&tempFilms.TMDB_ID, &tempFilms.Title, &tempFilms.TagLine, &tempFilms.Poster, &tempFilms.BackdropPoster, &tempFilms.Description, &tempFilms.Year, &tempFilms.Time); err != nil {
			return nil, fmt.Errorf("failed to scan films: %w", err)
		}
		films = append(films, tempFilms)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error fetching row: %w", err)
	}

	return films, nil
}

func (r *Repository) AddReview(ctx context.Context, rev *Review) (*Review, error) {
	// taking the filmid userid audiourl stars tmdb_id rewatch and saving it into this table
	// where are we getting the data from? doesnt matter all we are doing is inserting
	//
	query := "INSERT INTO reviews (id,film_id, user_id, audio_url, stars, tmdb_id, is_rewatch) VALUES (?,?,?,?,?,?,?)"

	rev.ID = randomHex(8)
	_, err := r.db.ExecContext(ctx, query, rev.ID, rev.FilmID, rev.UserID, rev.AudioURL, rev.Stars, rev.TMDB_ID, rev.IsRewatch)
	if err != nil {
		return nil, fmt.Errorf("failed to insert review %w", err)
	}

	return rev, nil
}
