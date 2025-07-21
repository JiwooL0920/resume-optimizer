package config

import (
	"fmt"
	"os"
)

// Config holds application configuration
type Config struct {
	DatabaseURL string
	Port        string
}

// DatabaseConfig holds database connection settings
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

// BuildDatabaseURL constructs the PostgreSQL connection string
func (dc *DatabaseConfig) BuildDatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dc.User, dc.Password, dc.Host, dc.Port, dc.Name, dc.SSLMode)
}

// LoadDatabaseConfig loads database configuration from environment variables
func LoadDatabaseConfig() *DatabaseConfig {
	getEnvOrDefault := func(key, defaultValue string) string {
		if value := os.Getenv(key); value != "" {
			return value
		}
		return defaultValue
	}
	
	return &DatabaseConfig{
		Host:     getEnvOrDefault("DB_HOST", "localhost"),
		Port:     getEnvOrDefault("DB_PORT", "5432"),
		User:     getEnvOrDefault("DB_USER", "user"),
		Password: getEnvOrDefault("DB_PASSWORD", "password"),
		Name:     getEnvOrDefault("DB_NAME", "resume_optimizer"),
		SSLMode:  getEnvOrDefault("DB_SSL_MODE", "disable"),
	}
}

// Load loads configuration from environment variables
func Load() *Config {
	fmt.Println("Loading resume-processor configuration...")
	
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
		DatabaseURL: databaseURL,
		Port:        getEnv("PORT", "8081"),
	}
	
	fmt.Printf("Resume-processor configuration loaded successfully\n")
	return config
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}