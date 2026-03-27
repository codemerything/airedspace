package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
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

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		MaxAge:   86400,
		HttpOnly: true,
		SameSite: http.SameSiteDefaultMode,
	})

	b, err := json.Marshal(map[string]string{"message": "user signed in successfully"})
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

func (h *Handler) SubmitReview(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "No file uploaded", http.StatusBadRequest)
		return
	}
	defer file.Close()

	dirPath := "./uploads/audio"
	err = os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		http.Error(w, "failed to create dir", http.StatusInternalServerError)
		return
	}

	filename := filepath.Join(dirPath, randomHex(8)+filepath.Ext(header.Filename))
	dst, err := os.Create(filename)
	if err != nil {
		http.Error(w, "Failed to create file", http.StatusInternalServerError)
		return
	}

	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "failed to save file", http.StatusInternalServerError)
		return
	}

	userID := r.Context().Value("user_id").(int)
	filmIDstr := r.FormValue("film_id")
	filmStarstr := r.FormValue("stars")
	filmTMDBstr := r.FormValue("tmdb_id")
	filmIRstr := r.FormValue("is_rewatch")
	filmReviewLikes := r.FormValue("review_likes")

	filmID, err := strconv.ParseInt(filmIDstr, 10, 64)
	if err != nil {
		http.Error(w, "invalid film_id", http.StatusBadRequest)
		return
	}

	filmStar, err := strconv.ParseFloat(filmStarstr, 64)
	if err != nil {
		http.Error(w, "invalid film_stars", http.StatusBadRequest)
		return
	}

	filmIR, err := strconv.ParseInt(filmIRstr, 10, 8)
	if err != nil {
		http.Error(w, "invalid value", http.StatusBadRequest)
	}

	reviewLikes, err := strconv.ParseInt(filmReviewLikes, 10, 64)
	if err != nil {
		http.Error(w, "invalid film reviews", http.StatusBadRequest)
		return
	}

	review := Review{
		FilmID:      filmID,
		UserID:      int64(userID),
		AudioURL:    filename,
		Stars:       filmStar,
		TMDB_ID:     filmTMDBstr,
		IsRewatch:   filmIR,
		ReviewLikes: reviewLikes,
	}

	err = h.service.SubmitReview(review)
	if err != nil {
		http.Error(w, "failed to save review", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
