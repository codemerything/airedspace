package main

import (
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port         int
	DatabaseURL  string
	Env          string
	JWTSecret    string
	TMDBKey      string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func Load() (*Config, error) {
	portStr := getEnv("PORT", ":8080")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Invalid port: %s", portStr)
	}

	cfg := Config{
		Port:         port,
		DatabaseURL:  requireEnv("DATABASE_URL"),
		Env:          getEnv("ENV", "development"),
		JWTSecret:    requireEnv("JWT_SECRET"),
		TMDBKey:      requireEnv("TMDB_API_KEY"),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return &cfg, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func requireEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("required environment variable %s is not set", key)
	}
	return value
}
