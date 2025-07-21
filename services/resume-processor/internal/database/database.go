package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/resume-optimizer/resume-processor/internal/models"
)

var DB *gorm.DB

func InitDatabase(databaseURL string) {
	var err error
	DB, err = gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate the schema
	err = DB.AutoMigrate(
		&models.User{},
		&models.Resume{},
		&models.OptimizationSession{},
		&models.Feedback{},
		&models.UserAPIKey{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database connection established and migrated")
}

func GetDB() *gorm.DB {
	return DB
}
