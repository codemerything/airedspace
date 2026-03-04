// handles storage of data to database explicitly

package main

import (
	"context"
	"database/sql"
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
