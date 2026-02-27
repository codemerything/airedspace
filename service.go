package main

type CreateUserInput struct {
	Username string
	Email    string
	Password string
}

type Service struct {
	repo *Repository
}

func (s *Service) CreateUser(input CreateUserInput) (int64, error) {

}
