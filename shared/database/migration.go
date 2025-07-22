package database

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/resume-optimizer/shared/config"
	"github.com/resume-optimizer/shared/errors"
	"github.com/rs/zerolog/log"
)

// MigrationService handles database migrations
type MigrationService struct {
	migrate *migrate.Migrate
	config  *config.DatabaseConfig
}

// NewMigrationService creates a new migration service
func NewMigrationService(cfg *config.DatabaseConfig, migrationsPath string) (*MigrationService, error) {
	// Connect to database
	db, err := sql.Open("postgres", cfg.GetDSN())
	if err != nil {
		return nil, errors.NewDatabaseError(fmt.Errorf("failed to connect to database: %w", err))
	}

	// Create postgres driver instance
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, errors.NewDatabaseError(fmt.Errorf("failed to create postgres driver: %w", err))
	}

	// Create migrate instance
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"postgres",
		driver,
	)
	if err != nil {
		return nil, errors.NewDatabaseError(fmt.Errorf("failed to create migrate instance: %w", err))
	}

	return &MigrationService{
		migrate: m,
		config:  cfg,
	}, nil
}

// Up runs all pending migrations
func (m *MigrationService) Up() error {
	log.Info().Msg("Running database migrations...")
	
	err := m.migrate.Up()
	if err != nil {
		if err == migrate.ErrNoChange {
			log.Info().Msg("No pending migrations")
			return nil
		}
		return errors.NewDatabaseError(fmt.Errorf("migration up failed: %w", err))
	}

	log.Info().Msg("Database migrations completed successfully")
	return nil
}

// Down rolls back one migration
func (m *MigrationService) Down() error {
	log.Info().Msg("Rolling back one migration...")
	
	err := m.migrate.Steps(-1)
	if err != nil {
		if err == migrate.ErrNoChange {
			log.Info().Msg("No migrations to rollback")
			return nil
		}
		return errors.NewDatabaseError(fmt.Errorf("migration down failed: %w", err))
	}

	log.Info().Msg("Migration rollback completed")
	return nil
}

// Version returns the current migration version
func (m *MigrationService) Version() (uint, bool, error) {
	version, dirty, err := m.migrate.Version()
	if err != nil {
		if err == migrate.ErrNilVersion {
			return 0, false, nil
		}
		return 0, false, errors.NewDatabaseError(fmt.Errorf("failed to get migration version: %w", err))
	}

	return version, dirty, nil
}

// Force sets the migration version without running migrations
func (m *MigrationService) Force(version int) error {
	log.Warn().Int("version", version).Msg("Forcing migration version")
	
	err := m.migrate.Force(version)
	if err != nil {
		return errors.NewDatabaseError(fmt.Errorf("failed to force migration version: %w", err))
	}

	log.Info().Int("version", version).Msg("Migration version forced")
	return nil
}

// Drop drops the entire database schema
func (m *MigrationService) Drop() error {
	log.Warn().Msg("Dropping entire database schema...")
	
	err := m.migrate.Drop()
	if err != nil {
		return errors.NewDatabaseError(fmt.Errorf("failed to drop database: %w", err))
	}

	log.Info().Msg("Database schema dropped")
	return nil
}

// Close closes the migration service
func (m *MigrationService) Close() error {
	sourceErr, dbErr := m.migrate.Close()
	if sourceErr != nil {
		log.Error().Err(sourceErr).Msg("Error closing migration source")
	}
	if dbErr != nil {
		log.Error().Err(dbErr).Msg("Error closing migration database")
		return dbErr
	}
	return sourceErr
}

// Status returns a detailed status of migrations
func (m *MigrationService) Status() (*MigrationStatus, error) {
	version, dirty, err := m.Version()
	if err != nil {
		return nil, err
	}

	status := &MigrationStatus{
		CurrentVersion: version,
		IsDirty:        dirty,
	}

	log.Info().
		Uint("version", version).
		Bool("dirty", dirty).
		Msg("Migration status")

	return status, nil
}

// MigrationStatus represents the current migration status
type MigrationStatus struct {
	CurrentVersion uint `json:"current_version"`
	IsDirty        bool `json:"is_dirty"`
}