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
	query := "INSERT INTO users (username,email,password) VALUES (?,?,?)"

	res, err := r.db.ExecContext(ctx, query, u.Username, u.Email, u.Password)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	u.UserID, err = res.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get userid: %w", err)
	}
	return nil
}
