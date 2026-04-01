package main

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
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

type NavLink struct {
	Name string
	URL  string
}

type TemplateData struct {
	Nav      []NavLink
	Username any
	Films    []Films
	Title    string
}

func (h *Handler) Welcome(w http.ResponseWriter, r *http.Request) {

	data := TemplateData{
		Username: r.Context().Value("username"),
	}

	if data.Username != nil {
		data.Nav = []NavLink{
			{"Search", "/search"},
			{"Profile", "/profile"},
			{"Logout", "/logout"},
		}
	} else {
		data.Nav = []NavLink{
			{"Sign in", "/signin"},
			{"Create account", "/signup"},
		}
	}

	files := []string{
		"./templates/index.tmpl.html",
		"./templates/base.tmpl.html",
		"./templates/nav.partial.tmpl.html",
		"./templates/footer.partial.tmpl.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	ts.ExecuteTemplate(w, "base", data)
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	data := TemplateData{
		Username: r.Context().Value("username"),
	}

	if data.Username != nil {
		data.Nav = []NavLink{
			{"Search", "/search"},
			{"Profile", "/profile"},
			{"Logout", "/logout"},
		}
	} else {
		data.Nav = []NavLink{
			{"Sign in", "/signin"},
			{"Create account", "/signup"},
		}
	}

	files := []string{
		"./templates/index.tmpl.html",
		"./templates/base.tmpl.html",
		"./templates/nav.partial.tmpl.html",
		"./templates/footer.partial.tmpl.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal server error", 500)
	}
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

	if r.Method == "GET" {
		if r.URL.Path != "/signin" {
			http.NotFound(w, r)
			return
		}

		data := TemplateData{
			Username: r.Context().Value("username"),
		}

		if data.Username != nil {
			data.Nav = []NavLink{
				{"Search", "/search"},
				{"Profile", "/profile"},
				{"Logout", "/logout"},
			}
		} else {
			data.Nav = []NavLink{
				{"Sign in", "/signin"},
				{"Create account", "/signup"},
			}
		}

		files := []string{
			"./templates/signin.tmpl.html",
			"./templates/base.tmpl.html",
			"./templates/nav.partial.tmpl.html",
			"./templates/footer.partial.tmpl.html",
		}

		ts, err := template.ParseFiles(files...)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

		ts.Execute(w, nil)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Internal server error", 500)
		}

	}

	if r.Method == "POST" {
		//POST request
		var input LoginInput
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		input.Username = r.FormValue("username")
		input.Password = r.FormValue("password")

		token, err := h.service.Login(input)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("<div id='error-message' style='color:red'>Invalid credentials. Try again</div>"))
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    token,
			MaxAge:   86400,
			HttpOnly: true,
			SameSite: http.SameSiteDefaultMode,
		})

		w.Header().Set("HX-Redirect", "/welcome")
		w.WriteHeader(http.StatusOK)

	}

}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {

	cookie := &http.Cookie{
		Name:   "token",
		Path:   "/",
		MaxAge: -1,
	}

	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

/*FILM RELATED HANDLERS */
func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {

	search := r.URL.Query().Get("search")

	if search == "" {
		http.Error(w, "Query field cannot be empty", http.StatusBadRequest)
		return
	}

	var searchInput = Movie{
		Title: search,
	}

	films, err := h.service.Search(searchInput)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// b, err := json.Marshal(map[string]any{"films": films})
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)
	// w.Write(b)
	data := TemplateData{
		Films: films,
		Title: search,
	}

	files := []string{
		"./templates/results.tmpl.html",
		"./templates/base.tmpl.html",
		"./templates/nav.partial.tmpl.html",
		"./templates/footer.partial.tmpl.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "body", data)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal server error", 500)
	}
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
