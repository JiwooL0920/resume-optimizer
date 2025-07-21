package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/resume-optimizer/auth-service/internal/config"
	"github.com/resume-optimizer/auth-service/internal/models"
)

var DB *gorm.DB

func InitDatabase(cfg *config.Config) {
	var err error
	DB, err = gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate the schema
	err = DB.AutoMigrate(&models.User{}, &models.UserAPIKey{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database connection established and migrated")
}

func GetDB() *gorm.DB {
	return DB
}
