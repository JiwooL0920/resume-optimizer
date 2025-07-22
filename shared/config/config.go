package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Database string `mapstructure:"database"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	SSLMode  string `mapstructure:"ssl_mode"`
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	JWTSecret        string `mapstructure:"jwt_secret"`
	JWTExpiration    int    `mapstructure:"jwt_expiration"`
	GoogleClientID   string `mapstructure:"google_client_id"`
	GoogleClientSecret string `mapstructure:"google_client_secret"`
	GoogleRedirectURL string `mapstructure:"google_redirect_url"`
}

// SecurityConfig holds security-related configuration
type SecurityConfig struct {
	EncryptionKey    string   `mapstructure:"encryption_key"`
	AllowedOrigins   []string `mapstructure:"allowed_origins"`
	RateLimitEnabled bool     `mapstructure:"rate_limit_enabled"`
	RateLimitRPS     int      `mapstructure:"rate_limit_rps"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port    int    `mapstructure:"port"`
	Host    string `mapstructure:"host"`
	Mode    string `mapstructure:"mode"` // debug, release, test
	Timeout int    `mapstructure:"timeout"`
}

// LogConfig holds logging configuration
type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"` // json, console
}

// MigrationConfig holds migration configuration
type MigrationConfig struct {
	Path string `mapstructure:"path"`
}

// ClientConfig holds client application configuration
type ClientConfig struct {
	BaseURL string `mapstructure:"base_url"`
}

// Config represents the application configuration
type Config struct {
	Database  DatabaseConfig  `mapstructure:"database"`
	Auth      AuthConfig      `mapstructure:"auth"`
	Security  SecurityConfig  `mapstructure:"security"`
	Server    ServerConfig    `mapstructure:"server"`
	Log       LogConfig       `mapstructure:"log"`
	Migration MigrationConfig `mapstructure:"migration"`
	Client    ClientConfig    `mapstructure:"client"`
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	config := &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			Database: getEnv("DB_NAME", "resume_optimizer"),
			Username: getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		Auth: AuthConfig{
			JWTSecret:         getEnv("JWT_SECRET", ""),
			JWTExpiration:     getEnvAsInt("JWT_EXPIRATION", 24*60*60), // 24 hours in seconds
			GoogleClientID:    getEnv("GOOGLE_CLIENT_ID", ""),
			GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
			GoogleRedirectURL: getEnv("GOOGLE_REDIRECT_URL", ""),
		},
		Security: SecurityConfig{
			EncryptionKey:    getEnv("ENCRYPTION_KEY", ""),
			AllowedOrigins:   getEnvAsSlice("ALLOWED_ORIGINS", []string{"http://localhost:3000"}),
			RateLimitEnabled: getEnvAsBool("RATE_LIMIT_ENABLED", true),
			RateLimitRPS:     getEnvAsInt("RATE_LIMIT_RPS", 10),
		},
		Server: ServerConfig{
			Port:    getEnvAsInt("PORT", 8080),
			Host:    getEnv("HOST", "0.0.0.0"),
			Mode:    getEnv("GIN_MODE", "debug"),
			Timeout: getEnvAsInt("SERVER_TIMEOUT", 30),
		},
		Log: LogConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
		Migration: MigrationConfig{
			Path: getEnv("MIGRATIONS_PATH", "./shared/database/migrations"),
		},
		Client: ClientConfig{
			BaseURL: getEnv("CLIENT_BASE_URL", "http://localhost:3000"),
		},
	}

	// Validate required configuration
	if err := validateConfig(config); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return config, nil
}

// validateConfig validates that required configuration is present
func validateConfig(config *Config) error {
	if config.Auth.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}

	if len(config.Auth.JWTSecret) < 32 {
		return fmt.Errorf("JWT_SECRET must be at least 32 characters long")
	}

	if config.Security.EncryptionKey == "" {
		return fmt.Errorf("ENCRYPTION_KEY is required")
	}

	if len(config.Security.EncryptionKey) != 32 {
		return fmt.Errorf("ENCRYPTION_KEY must be exactly 32 characters long")
	}

	if config.Database.Password == "" {
		return fmt.Errorf("DB_PASSWORD is required")
	}

	return nil
}

// GetDSN returns the database connection string
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.Username, c.Password, c.Database, c.SSLMode)
}

// Helper functions for environment variable parsing
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvAsSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, ",")
	}
	return defaultValue
}