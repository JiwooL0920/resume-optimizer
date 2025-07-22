package models

import "time"

// Resume represents a user's uploaded resume
type Resume struct {
	ID              string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID          string    `json:"user_id" gorm:"not null;type:uuid;index"`
	Title           string    `json:"title" gorm:"not null"`
	OriginalContent string    `json:"original_content" gorm:"not null;type:text"` // File path
	ExtractedText   string    `json:"extracted_text" gorm:"type:text"`           // Extracted text content
	FileType        string    `json:"file_type" gorm:"not null;default:'.pdf'"`
	FileSize        *int      `json:"file_size"`
	IsActive        bool      `json:"is_active" gorm:"default:true"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	
	// Associations
	User                 User                  `json:"user,omitempty" gorm:"foreignKey:UserID"`
	OptimizationSessions []OptimizationSession `json:"optimization_sessions,omitempty" gorm:"foreignKey:ResumeID"`
}

// OptimizationSession represents an AI optimization session
type OptimizationSession struct {
	ID                 string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID             string    `json:"user_id" gorm:"not null;type:uuid;index"`
	ResumeID           string    `json:"resume_id" gorm:"not null;type:uuid;index"`
	JobDescriptionURL  *string   `json:"job_description_url"`
	JobDescriptionText *string   `json:"job_description_text" gorm:"type:text"`
	AIModel            string    `json:"ai_model" gorm:"not null"`
	KeepOnePage        bool      `json:"keep_one_page" gorm:"default:false"`
	OptimizedContent   *string   `json:"optimized_content" gorm:"type:text"`
	Status             string    `json:"status" gorm:"default:'pending'"` // pending, processing, completed, failed
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	
	// Associations
	User     User       `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Resume   Resume     `json:"resume,omitempty" gorm:"foreignKey:ResumeID"`
	Feedback []Feedback `json:"feedback,omitempty" gorm:"foreignKey:SessionID"`
}

// Feedback represents user feedback on optimization results
type Feedback struct {
	ID               string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	SessionID        string    `json:"session_id" gorm:"not null;type:uuid;index"`
	SectionHighlight string    `json:"section_highlight" gorm:"not null;type:text"`
	UserComment      string    `json:"user_comment" gorm:"not null;type:text"`
	IsProcessed      bool      `json:"is_processed" gorm:"default:false"`
	CreatedAt        time.Time `json:"created_at"`
	
	// Associations
	Session OptimizationSession `json:"session,omitempty" gorm:"foreignKey:SessionID"`
}