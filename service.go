// the chef: handles business logic

package main

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
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
type LoginInput struct {
	Username string
	Password string
}

type Service struct {
	repo *Repository
}

type Claims struct {
	Username string `json:"username"`
}

// how to generate a jwt token

var mySecretKey = []byte("secretsecret")

func generateToken(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(mySecretKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// jwt token ends

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

func (s *Service) Login(input LoginInput) (string, error) {
	// check if user exists
	//
	user := &User{
		Username: input.Username,
	}

	if input.Username == "" {
		return "", errors.New("Insert email")
	}

	if input.Password == "" {
		return "", errors.New("Insert password")
	}

	fetchedUser, err := s.repo.FetchUserByUsername(user)
	if err != nil {
		return "", errors.New("Fetch failed")
	}

	if err = bcrypt.CompareHashAndPassword([]byte(fetchedUser.Password), []byte(input.Password)); err != nil {
		return "", errors.New("Invalid credentials")
	}

	jwt, err := generateToken(user.Username)
	if err != nil {
		return "", errors.New("Failed to generate JWT token")
	}

	return jwt, nil
}
