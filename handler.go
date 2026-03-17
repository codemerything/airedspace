package main

import (
	"encoding/json"
	"io"
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

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	var user SignUpRequest
	body, _ := io.ReadAll(r.Body)
	err := json.Unmarshal(body, &user)

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

	b, err := json.Marshal(map[string]string{"message": "signup sucessful"})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	var input LoginInput
	body, _ := io.ReadAll(r.Body)
	err := json.Unmarshal(body, &input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := h.service.Login(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	b, err := json.Marshal(map[string]string{"token": token})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {

	title := r.URL.Query().Get("title")

	if title == "" {
		http.Error(w, "Query field cannot be empty", http.StatusBadRequest)
		return
	}

	var searchInput = Movie{
		Title: title,
	}

	films, err := h.service.Search(searchInput)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	b, err := json.Marshal(map[string]any{"films": films})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
