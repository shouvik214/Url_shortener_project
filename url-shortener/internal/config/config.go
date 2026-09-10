package config

import (
	"errors"
	"fmt"
	"os"
)

type Config struct {
	Port        string
	DatabaseURL string
	RedisAddr   string
	BaseURL     string
	FrontendURL string
}

func Load() (*Config, error) {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, errors.New("DATABASE_URL is not set")
	}

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:" + port
	}

	frontendURL := os.Getenv("FRONTEND_URL")

	if frontendURL == "" {
		return nil, fmt.Errorf("FRONTEND_URL is not set")
	}

	return &Config{
		Port:        port,
		DatabaseURL: databaseURL,
		RedisAddr:   redisAddr,
		BaseURL:     baseURL,
		FrontendURL: frontendURL,
	}, nil

}
