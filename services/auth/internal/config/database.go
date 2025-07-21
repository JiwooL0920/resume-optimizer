package config

import (
	"fmt"
	"os"
)

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