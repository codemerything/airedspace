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

func (r *Repository) CreateFilm(ctx context.Context, f *Films) (*Films, error) {
	query := "INSERT INTO films (tmdb_id, title, tag_line, poster, foreground_poster, description, release_year, runtime) VALUES (?,?,?,?,?,?,?,?)"

	res, err := r.db.ExecContext(ctx, query, f.TMDB_ID, f.Title, f.TagLine, f.Poster, f.BackdropPoster, f.Description, f.Year, f.Time)
	if err != nil {
		return nil, fmt.Errorf("failed to insert a film: %w", err)
	}

	f.FilmID, err = res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get filmid: %w", err)
	}

	return f, nil

}

func (r *Repository) SearchFilm(ctx context.Context, f *Films) ([]Films, error) {

	var films []Films
	query := `SELECT  tmdb_id, title, tag_line, poster, foreground_poster, description, release_year FROM films WHERE title LIKE ?`

	rows, err := r.db.QueryContext(ctx, query, (f.Title + "%"))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch row: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var tempFilms Films
		if err := rows.Scan(&tempFilms.TMDB_ID, &tempFilms.Title, &tempFilms.TagLine, &tempFilms.Poster, &tempFilms.BackdropPoster, &tempFilms.Description, &tempFilms.Year); err != nil {
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
	query := "INSERT INTO reviews (film_id, user_id, audio_url, stars, tmdb_id, is_rewatch) VALUES (?,?,?,?,?,?)"

	res, err := r.db.ExecContext(ctx, query, rev.FilmID, rev.UserID, rev.AudioURL, rev.Stars, rev.TMDB_ID, rev.IsRewatch)

	if err != nil {
		return nil, fmt.Errorf("failed to insert review %w", err)
	}

	rev.ID, err = res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get review id: %w", err)
	}

	return rev, nil
}
