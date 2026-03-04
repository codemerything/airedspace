package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type Handlers struct {
	h *Handler
}

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	cfg, err := Load()
	if err != nil {
		log.Fatalf("Failed to load config %v", err)
	}

	// database connection
	db, err := NewConnection(cfg.DatabaseURL)

	// setting up routes
	repo := &Repository{db: db.DB}
	service := &Service{repo: repo}
	h := &Handler{service: service}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", Welcome)
	mux.HandleFunc("POST /signin", SignIn)
	mux.HandleFunc("POST /signup", h.SignUp)
	mux.HandleFunc("GET /movies", Movies)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		Handler:      mux,
	}

	go func() {
		log.Printf("Starting server on port %d", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutting off server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
