package database

import (
	"fmt"
	"time"

	"github.com/resume-optimizer/shared/config"
	"github.com/resume-optimizer/shared/errors"
	"github.com/resume-optimizer/shared/models"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

// Connect establishes database connection
func Connect(cfg *config.DatabaseConfig) error {
	dsn := cfg.GetDSN()
	
	// Configure GORM logger
	gormLogger := logger.Default.LogMode(logger.Silent)
	if log.Debug().Enabled() {
		gormLogger = logger.Default.LogMode(logger.Info)
	}

	// Connect to database
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                 gormLogger,
		SkipDefaultTransaction: true, // Better performance
		PrepareStmt:           true,  // Cache prepared statements
	})
	if err != nil {
		return errors.NewDatabaseError(fmt.Errorf("failed to connect to database: %w", err))
	}

	// Get underlying sql.DB
	sqlDB, err := database.DB()
	if err != nil {
		return errors.NewDatabaseError(fmt.Errorf("failed to get sql.DB: %w", err))
	}

	// Configure connection pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	db = database

	log.Info().
		Str("host", cfg.Host).
		Int("port", cfg.Port).
		Str("database", cfg.Database).
		Msg("Connected to database")

	// Run auto migration to ensure tables exist
	if err := AutoMigrate(); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return db
}

// Disconnect closes the database connection
func Disconnect() error {
	if db == nil {
		return nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		return errors.NewDatabaseError(fmt.Errorf("failed to get sql.DB: %w", err))
	}

	if err := sqlDB.Close(); err != nil {
		return errors.NewDatabaseError(fmt.Errorf("failed to close database connection: %w", err))
	}

	log.Info().Msg("Database connection closed")
	return nil
}

// AutoMigrate runs GORM auto migration (for development only)
func AutoMigrate() error {
	if db == nil {
		return errors.NewDatabaseError(fmt.Errorf("database not connected"))
	}

	log.Info().Msg("Running GORM auto migration...")

	// Clean up any existing data with NULL user_id values before migration
	log.Info().Msg("Cleaning up data inconsistencies...")
	db.Exec("DELETE FROM resumes WHERE user_id IS NULL")
	db.Exec("DELETE FROM optimization_sessions WHERE user_id IS NULL")
	db.Exec("DELETE FROM feedback WHERE session_id NOT IN (SELECT id FROM optimization_sessions)")

	err := db.AutoMigrate(
		&models.User{},
		&models.UserAPIKey{},
		&models.Resume{},
		&models.OptimizationSession{},
		&models.Feedback{},
	)

	if err != nil {
		return errors.NewDatabaseError(fmt.Errorf("auto migration failed: %w", err))
	}

	log.Info().Msg("GORM auto migration completed")
	return nil
}

// HealthCheck performs a database health check
func HealthCheck() error {
	if db == nil {
		return errors.NewDatabaseError(fmt.Errorf("database not connected"))
	}

	sqlDB, err := db.DB()
	if err != nil {
		return errors.NewDatabaseError(fmt.Errorf("failed to get sql.DB: %w", err))
	}

	if err := sqlDB.Ping(); err != nil {
		return errors.NewDatabaseError(fmt.Errorf("database ping failed: %w", err))
	}

	return nil
}

// GetStats returns database connection statistics
func GetStats() (*DatabaseStats, error) {
	if db == nil {
		return nil, errors.NewDatabaseError(fmt.Errorf("database not connected"))
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, errors.NewDatabaseError(fmt.Errorf("failed to get sql.DB: %w", err))
	}

	stats := sqlDB.Stats()
	return &DatabaseStats{
		OpenConnections: stats.OpenConnections,
		InUse:          stats.InUse,
		Idle:           stats.Idle,
		WaitCount:      stats.WaitCount,
		WaitDuration:   stats.WaitDuration,
	}, nil
}

// DatabaseStats represents database connection statistics
type DatabaseStats struct {
	OpenConnections int           `json:"open_connections"`
	InUse          int           `json:"in_use"`
	Idle           int           `json:"idle"`
	WaitCount      int64         `json:"wait_count"`
	WaitDuration   time.Duration `json:"wait_duration"`
}

// Transaction executes a function within a database transaction
func Transaction(fn func(*gorm.DB) error) error {
	if db == nil {
		return errors.NewDatabaseError(fmt.Errorf("database not connected"))
	}

	return db.Transaction(fn)
}

// WithContext returns a database instance with context
func WithContext(db *gorm.DB) *gorm.DB {
	return db.WithContext(db.Statement.Context)
}