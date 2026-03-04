// the chef: handles business logic

package main

import (
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type CreateUserInput struct {
	Username  string
	Email     string
	Password  string
	Name      string
	LastName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Service struct {
	repo *Repository
}

func (s *Service) CreateUser(input CreateUserInput) error {
	if input.Username == "" {
		return errors.New("username is required")
	}

	if input.Email == "" {
		return errors.New("email is required")
	}

	if input.Password == "" {
		return errors.New("password is required")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	if err != nil {
		return errors.New("Hashing failed")
	}

	user := &User{
		Username:  input.Username,
		Email:     input.Email,
		Name:      input.Name,
		LastName:  input.LastName,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = s.repo.CreateUser(context.TODO(), user)
	return err

}
