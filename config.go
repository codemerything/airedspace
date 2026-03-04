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
		DatabaseURL:  getEnv("DATABASE_URL", "root:fuckingmakeamess@tcp(localhost:3306)/airedspace"),
		Env:          getEnv("ENV", "development"),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return &cfg, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != " " {
		return value
	}
	return fallback
}
