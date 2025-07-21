package models

import (
	"time"
	"gorm.io/gorm"
)

type User struct {
	ID         string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Email      string    `json:"email" gorm:"uniqueIndex;not null"`
	GoogleID   *string   `json:"google_id" gorm:"uniqueIndex"`
	Name       string    `json:"name" gorm:"not null"`
	PictureURL *string   `json:"picture_url"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	
	Resumes             []Resume             `json:"resumes,omitempty" gorm:"foreignKey:UserID"`
	OptimizationSessions []OptimizationSession `json:"optimization_sessions,omitempty" gorm:"foreignKey:UserID"`
	APIKeys             []UserAPIKey         `json:"api_keys,omitempty" gorm:"foreignKey:UserID"`
}