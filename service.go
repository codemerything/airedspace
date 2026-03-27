// the chef: handles business logic

package main

import (
	"context"
	"errors"
	"fmt"
	"math"
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

type Movie struct {
	Title string
}

type Service struct {
	repo *Repository
	tmdb *TMDB
	cfg  *Config
}

type Claims struct {
	Username string `json:"username"`
}

// how to generate a jwt token

func generateToken(username string, userID float64, cfg *Config) (string, error) {

	claims := jwt.MapClaims{
		"username": username,
		"user_id":  userID,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	fmt.Println("signing with secret:", string(cfg.JWTSecret))
	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	fmt.Println("sign error:", err)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// jwt code ends

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

	if input.Username == "" {
		return "", errors.New("Insert email")
	}

	user := &User{
		Username: input.Username,
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

	jwt, err := generateToken(fetchedUser.Username, float64(fetchedUser.UserID), s.cfg)
	if err != nil {
		return "", errors.New("Failed to generate JWT token")
	}

	return jwt, nil
}

func (s *Service) Search(input Movie) ([]Films, error) {

	// user searches for a name of movie
	// check the database if that search already exists and give it to the user
	// if not search the api and add that search to the films
	// save that into the database (pessimistic approach)
	// display to user
	//
	// if input is empty then return error
	//

	if input.Title == "" {
		return nil, errors.New("Empty input field")
	}

	title := &Films{
		Title: input.Title,
	}

	fetchedFilm, err := s.repo.SearchFilm(context.Background(), title)
	if err != nil {
		return nil, err
	}

	if len(fetchedFilm) == 0 {
		id, err := s.tmdb.FetchMovieID(input.Title)
		if err != nil {
			return nil, err
		}

		details, err := s.tmdb.FetchFilmDetails(id)
		if err != nil {
			return nil, errors.New("failed to fetch film details")
		}

		filmsdeets := Films{
			Title:          details.Title,
			TMDB_ID:        details.FilmID,
			Year:           details.Year,
			Poster:         details.PosterPath,
			BackdropPoster: details.BackdropPath,
			Description:    details.Description,
			TagLine:        details.TagLine,
			Time:           details.Time,
		}
		tx, err := s.repo.BeginTx(context.Background())
		if err != nil {
			return nil, err
		}

		films, err := s.repo.CreateFilm(context.Background(), tx, &filmsdeets)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		// TODO: loop through each cast and save to database for persons and films cast
		for _, value := range details.Credits.Cast {
			persons := Persons{
				TMDB_ID: value.CastID,
				Name:    value.Name,
			}

			person, err := s.repo.CreatePerson(context.Background(), tx, &persons)
			if err != nil {
				tx.Rollback()
				return nil, err
			}

			filmsCast := FilmsCast{
				FilmID:        films.FilmID,
				CastID:        person.ID,
				CharacterName: value.CharacterName,
			}

			err = s.repo.CreateFilmCast(context.Background(), tx, &filmsCast)
			if err != nil {
				tx.Rollback()
				return nil, err
			}

		}

		for _, value := range details.Credits.Genres {
			genres := Genres{
				Name: value.Name,
			}

			g, err := s.repo.CreateGenre(context.Background(), tx, &genres)
			if err != nil {
				tx.Rollback()
				return nil, err
			}

			filmsGenre := FilmsGenre{
				FilmID:  films.FilmID,
				GenreID: g.Id,
			}
			err = s.repo.CreateFilmGenre(context.Background(), tx, &filmsGenre)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}
		tx.Commit()
		return []Films{*films}, nil

	}

	return fetchedFilm, nil
}

func (s *Service) SubmitReview(review Review) error {

	if review.AudioURL == "" {
		return errors.New("Invalid file for audio")
	}

	if review.Stars < 0.5 || review.Stars > 5.0 || math.Mod(review.Stars, 0.5) != 0 {
		return errors.New("Wrong review format")
	}

	if review.FilmID == 0 {
		return errors.New("Invalid filmid")
	}

	_, err := s.repo.AddReview(context.Background(), &review)
	if err != nil {
		return err
	}

	return nil
}
