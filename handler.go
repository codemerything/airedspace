package main

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	service *Service
}

type SignUpRequest struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	LastName string `json:"lastname"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Welcome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world"))
}

func SignIn(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	var user SignUpRequest
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	input := CreateUserInput{
		Username: user.Username,
		Password: user.Password,
		Email:    user.Email,
		Name:     user.Name,
		LastName: user.LastName,
	}

	err = h.service.CreateUser(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "signup successful"})

}

func Movies(w http.ResponseWriter, r *http.Request) {}
