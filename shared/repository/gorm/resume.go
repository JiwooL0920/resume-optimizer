package gorm

import (
	"context"
	"fmt"

	"github.com/resume-optimizer/shared/errors"
	"github.com/resume-optimizer/shared/models"
	"github.com/resume-optimizer/shared/repository"
	"gorm.io/gorm"
)

type resumeRepository struct {
	db *gorm.DB
}

// NewResumeRepository creates a new resume repository
func NewResumeRepository(db *gorm.DB) repository.ResumeRepository {
	return &resumeRepository{db: db}
}

// Create creates a new resume
func (r *resumeRepository) Create(ctx context.Context, resume *models.Resume) error {
	if err := r.db.WithContext(ctx).Create(resume).Error; err != nil {
		if isForeignKeyConstraintError(err) {
			return errors.NewAppError(
				errors.ErrCodeNotFound,
				"User not found",
				err,
			)
		}
		return errors.NewDatabaseError(fmt.Errorf("failed to create resume: %w", err))
	}
	return nil
}

// GetByID retrieves a resume by ID
func (r *resumeRepository) GetByID(ctx context.Context, id string) (*models.Resume, error) {
	var resume models.Resume
	if err := r.db.WithContext(ctx).First(&resume, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewAppError(
				errors.ErrCodeNotFound,
				"Resume not found",
				err,
			)
		}
		return nil, errors.NewDatabaseError(fmt.Errorf("failed to get resume by ID: %w", err))
	}
	return &resume, nil
}

// GetByUserID retrieves all resumes for a user
func (r *resumeRepository) GetByUserID(ctx context.Context, userID string) ([]*models.Resume, error) {
	var resumes []*models.Resume
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&resumes).Error; err != nil {
		return nil, errors.NewDatabaseError(fmt.Errorf("failed to get resumes by user ID: %w", err))
	}
	return resumes, nil
}

// GetActiveByUserID retrieves active resumes for a user
func (r *resumeRepository) GetActiveByUserID(ctx context.Context, userID string) ([]*models.Resume, error) {
	var resumes []*models.Resume
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_active = ?", userID, true).
		Order("created_at DESC").
		Find(&resumes).Error; err != nil {
		return nil, errors.NewDatabaseError(fmt.Errorf("failed to get active resumes by user ID: %w", err))
	}
	return resumes, nil
}

// Update updates a resume
func (r *resumeRepository) Update(ctx context.Context, resume *models.Resume) error {
	if err := r.db.WithContext(ctx).Save(resume).Error; err != nil {
		return errors.NewDatabaseError(fmt.Errorf("failed to update resume: %w", err))
	}
	return nil
}

// Delete deletes a resume
func (r *resumeRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&models.Resume{}, "id = ?", id)
	if result.Error != nil {
		return errors.NewDatabaseError(fmt.Errorf("failed to delete resume: %w", result.Error))
	}
	if result.RowsAffected == 0 {
		return errors.NewAppError(
			errors.ErrCodeNotFound,
			"Resume not found",
			nil,
		)
	}
	return nil
}

// List retrieves resumes with pagination
func (r *resumeRepository) List(ctx context.Context, offset, limit int) ([]*models.Resume, error) {
	var resumes []*models.Resume
	query := r.db.WithContext(ctx)
	query = applyPagination(query, offset, limit)
	query = applySorting(query, "created_at", "desc")
	
	if err := query.Find(&resumes).Error; err != nil {
		return nil, errors.NewDatabaseError(fmt.Errorf("failed to list resumes: %w", err))
	}
	return resumes, nil
}