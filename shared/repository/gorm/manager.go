package gorm

import (
	"github.com/resume-optimizer/shared/repository"
	"gorm.io/gorm"
)

// repositoryManager implements the RepositoryManager interface
type repositoryManager struct {
	db                      *gorm.DB
	userRepo                repository.UserRepository
	resumeRepo              repository.ResumeRepository
	optimizationSessionRepo repository.OptimizationSessionRepository
	feedbackRepo            repository.FeedbackRepository
	userAPIKeyRepo          repository.UserAPIKeyRepository
}

// NewRepositoryManager creates a new repository manager
func NewRepositoryManager(db *gorm.DB) repository.RepositoryManager {
	return &repositoryManager{
		db:                      db,
		userRepo:                NewUserRepository(db),
		resumeRepo:              NewResumeRepository(db),
		optimizationSessionRepo: NewOptimizationSessionRepository(db),
		feedbackRepo:            NewFeedbackRepository(db),
		userAPIKeyRepo:          NewUserAPIKeyRepository(db),
	}
}

// User returns the user repository
func (rm *repositoryManager) User() repository.UserRepository {
	return rm.userRepo
}

// Resume returns the resume repository
func (rm *repositoryManager) Resume() repository.ResumeRepository {
	return rm.resumeRepo
}

// OptimizationSession returns the optimization session repository
func (rm *repositoryManager) OptimizationSession() repository.OptimizationSessionRepository {
	return rm.optimizationSessionRepo
}

// Feedback returns the feedback repository
func (rm *repositoryManager) Feedback() repository.FeedbackRepository {
	return rm.feedbackRepo
}

// UserAPIKey returns the user API key repository
func (rm *repositoryManager) UserAPIKey() repository.UserAPIKeyRepository {
	return rm.userAPIKeyRepo
}

// Transaction executes a function within a database transaction
func (rm *repositoryManager) Transaction(fn func(repository.RepositoryManager) error) error {
	return rm.db.Transaction(func(tx *gorm.DB) error {
		txManager := &repositoryManager{
			db:                      tx,
			userRepo:                NewUserRepository(tx),
			resumeRepo:              NewResumeRepository(tx),
			optimizationSessionRepo: NewOptimizationSessionRepository(tx),
			feedbackRepo:            NewFeedbackRepository(tx),
			userAPIKeyRepo:          NewUserAPIKeyRepository(tx),
		}
		return fn(txManager)
	})
}