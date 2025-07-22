package repository

import (
	"context"

	"github.com/resume-optimizer/shared/models"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByGoogleID(ctx context.Context, googleID string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, offset, limit int) ([]*models.User, error)
}

// ResumeRepository defines the interface for resume data operations
type ResumeRepository interface {
	Create(ctx context.Context, resume *models.Resume) error
	GetByID(ctx context.Context, id string) (*models.Resume, error)
	GetByUserID(ctx context.Context, userID string) ([]*models.Resume, error)
	Update(ctx context.Context, resume *models.Resume) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, offset, limit int) ([]*models.Resume, error)
	GetActiveByUserID(ctx context.Context, userID string) ([]*models.Resume, error)
}

// OptimizationSessionRepository defines the interface for optimization session operations
type OptimizationSessionRepository interface {
	Create(ctx context.Context, session *models.OptimizationSession) error
	GetByID(ctx context.Context, id string) (*models.OptimizationSession, error)
	GetByUserID(ctx context.Context, userID string) ([]*models.OptimizationSession, error)
	GetByResumeID(ctx context.Context, resumeID string) ([]*models.OptimizationSession, error)
	Update(ctx context.Context, session *models.OptimizationSession) error
	UpdateStatus(ctx context.Context, id, status string) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, offset, limit int) ([]*models.OptimizationSession, error)
	GetByStatus(ctx context.Context, status string) ([]*models.OptimizationSession, error)
}

// FeedbackRepository defines the interface for feedback operations
type FeedbackRepository interface {
	Create(ctx context.Context, feedback *models.Feedback) error
	GetByID(ctx context.Context, id string) (*models.Feedback, error)
	GetBySessionID(ctx context.Context, sessionID string) ([]*models.Feedback, error)
	Update(ctx context.Context, feedback *models.Feedback) error
	Delete(ctx context.Context, id string) error
	MarkAsProcessed(ctx context.Context, id string) error
	GetUnprocessed(ctx context.Context) ([]*models.Feedback, error)
}

// UserAPIKeyRepository defines the interface for API key operations
type UserAPIKeyRepository interface {
	Create(ctx context.Context, apiKey *models.UserAPIKey) error
	GetByID(ctx context.Context, id string) (*models.UserAPIKey, error)
	GetByUserID(ctx context.Context, userID string) ([]*models.UserAPIKey, error)
	GetByUserIDAndProvider(ctx context.Context, userID, provider string) (*models.UserAPIKey, error)
	Update(ctx context.Context, apiKey *models.UserAPIKey) error
	Delete(ctx context.Context, id string) error
	DeleteByUserIDAndProvider(ctx context.Context, userID, provider string) error
	List(ctx context.Context, offset, limit int) ([]*models.UserAPIKey, error)
}

// RepositoryManager provides access to all repositories
type RepositoryManager interface {
	User() UserRepository
	Resume() ResumeRepository
	OptimizationSession() OptimizationSessionRepository
	Feedback() FeedbackRepository
	UserAPIKey() UserAPIKeyRepository
	Transaction(fn func(RepositoryManager) error) error
}