package models

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID         string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Email      string    `json:"email" gorm:"uniqueIndex;not null"`
	GoogleID   *string   `json:"google_id" gorm:"uniqueIndex"`
	Name       string    `json:"name" gorm:"not null"`
	PictureURL *string   `json:"picture_url"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	
	// Associations (use omitempty to avoid circular loading)
	Resumes              []Resume              `json:"resumes,omitempty" gorm:"foreignKey:UserID"`
	OptimizationSessions []OptimizationSession `json:"optimization_sessions,omitempty" gorm:"foreignKey:UserID"`
	APIKeys              []UserAPIKey          `json:"api_keys,omitempty" gorm:"foreignKey:UserID"`
}

// UserAPIKey represents an encrypted API key for external services
type UserAPIKey struct {
	ID           string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID       string    `json:"user_id" gorm:"not null;type:uuid;index"`
	Provider     string    `json:"provider" gorm:"not null"` // openai, anthropic, google, etc.
	EncryptedKey string    `json:"-" gorm:"not null;type:text"`
	MaskedKey    string    `json:"masked_key" gorm:"-"` // Generated at runtime, not stored
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	
	// Associations
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName returns the table name for UserAPIKey
func (UserAPIKey) TableName() string {
	return "user_api_keys"
}