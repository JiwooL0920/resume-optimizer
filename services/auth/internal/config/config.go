package config

import (
	"fmt"
	"os"
)

type Config struct {
	DatabaseURL      string
	GoogleClientID   string
	GoogleClientSecret string
	JWTSecret       string
	RedirectURL     string
}

func Load() *Config {
	fmt.Println("Loading configuration...")
	
	// Try to get DATABASE_URL directly, otherwise build from components
	databaseURL := getEnv("DATABASE_URL", "")
	if databaseURL == "" {
		fmt.Println("DATABASE_URL not found, building from components...")
		dbConfig := LoadDatabaseConfig()
		databaseURL = dbConfig.BuildDatabaseURL()
		fmt.Printf("Built database URL: %s\n", databaseURL)
	} else {
		fmt.Printf("Using provided DATABASE_URL: %s\n", databaseURL)
	}
	
	config := &Config{
		DatabaseURL:        databaseURL,
		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		JWTSecret:         getEnv("JWT_SECRET", "your-secret-key"),
		RedirectURL:       getEnv("REDIRECT_URL", "http://localhost:8080/api/v1/auth/google/callback"),
	}
	
	fmt.Printf("Configuration loaded: GoogleClientID exists: %t\n", config.GoogleClientID != "")
	return config
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}