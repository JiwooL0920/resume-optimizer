package config

import "os"

type Config struct {
	DatabaseURL      string
	GoogleClientID   string
	GoogleClientSecret string
	JWTSecret       string
	RedirectURL     string
}

func Load() *Config {
	return &Config{
		DatabaseURL:        getEnv("DATABASE_URL", "postgres://user:password@localhost/resume_optimizer?sslmode=disable"),
		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		JWTSecret:         getEnv("JWT_SECRET", "your-secret-key"),
		RedirectURL:       getEnv("REDIRECT_URL", "http://localhost:3000/auth/callback"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}