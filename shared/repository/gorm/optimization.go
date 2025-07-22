package gorm

import (
	"context"
	"fmt"

	"github.com/resume-optimizer/shared/errors"
	"github.com/resume-optimizer/shared/models"
	"github.com/resume-optimizer/shared/repository"
	"gorm.io/gorm"
)

type optimizationSessionRepository struct {
	db *gorm.DB
}

func NewOptimizationSessionRepository(db *gorm.DB) repository.OptimizationSessionRepository {
	return &optimizationSessionRepository{db: db}
}

func (r *optimizationSessionRepository) Create(ctx context.Context, session *models.OptimizationSession) error {
	if err := r.db.WithContext(ctx).Create(session).Error; err != nil {
		return errors.NewDatabaseError(fmt.Errorf("failed to create optimization session: %w", err))
	}
	return nil
}

func (r *optimizationSessionRepository) GetByID(ctx context.Context, id string) (*models.OptimizationSession, error) {
	var session models.OptimizationSession
	if err := r.db.WithContext(ctx).First(&session, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewAppError(errors.ErrCodeNotFound, "Optimization session not found", err)
		}
		return nil, errors.NewDatabaseError(fmt.Errorf("failed to get optimization session: %w", err))
	}
	return &session, nil
}

func (r *optimizationSessionRepository) GetByUserID(ctx context.Context, userID string) ([]*models.OptimizationSession, error) {
	var sessions []*models.OptimizationSession
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Find(&sessions).Error; err != nil {
		return nil, errors.NewDatabaseError(fmt.Errorf("failed to get optimization sessions: %w", err))
	}
	return sessions, nil
}

func (r *optimizationSessionRepository) GetByResumeID(ctx context.Context, resumeID string) ([]*models.OptimizationSession, error) {
	var sessions []*models.OptimizationSession
	if err := r.db.WithContext(ctx).Where("resume_id = ?", resumeID).Order("created_at DESC").Find(&sessions).Error; err != nil {
		return nil, errors.NewDatabaseError(fmt.Errorf("failed to get optimization sessions: %w", err))
	}
	return sessions, nil
}

func (r *optimizationSessionRepository) Update(ctx context.Context, session *models.OptimizationSession) error {
	if err := r.db.WithContext(ctx).Save(session).Error; err != nil {
		return errors.NewDatabaseError(fmt.Errorf("failed to update optimization session: %w", err))
	}
	return nil
}

func (r *optimizationSessionRepository) UpdateStatus(ctx context.Context, id, status string) error {
	result := r.db.WithContext(ctx).Model(&models.OptimizationSession{}).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		return errors.NewDatabaseError(fmt.Errorf("failed to update status: %w", result.Error))
	}
	if result.RowsAffected == 0 {
		return errors.NewAppError(errors.ErrCodeNotFound, "Optimization session not found", nil)
	}
	return nil
}

func (r *optimizationSessionRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&models.OptimizationSession{}, "id = ?", id)
	if result.Error != nil {
		return errors.NewDatabaseError(fmt.Errorf("failed to delete optimization session: %w", result.Error))
	}
	if result.RowsAffected == 0 {
		return errors.NewAppError(errors.ErrCodeNotFound, "Optimization session not found", nil)
	}
	return nil
}

func (r *optimizationSessionRepository) List(ctx context.Context, offset, limit int) ([]*models.OptimizationSession, error) {
	var sessions []*models.OptimizationSession
	query := applyPagination(r.db.WithContext(ctx), offset, limit)
	query = applySorting(query, "created_at", "desc")
	if err := query.Find(&sessions).Error; err != nil {
		return nil, errors.NewDatabaseError(fmt.Errorf("failed to list optimization sessions: %w", err))
	}
	return sessions, nil
}

func (r *optimizationSessionRepository) GetByStatus(ctx context.Context, status string) ([]*models.OptimizationSession, error) {
	var sessions []*models.OptimizationSession
	if err := r.db.WithContext(ctx).Where("status = ?", status).Order("created_at ASC").Find(&sessions).Error; err != nil {
		return nil, errors.NewDatabaseError(fmt.Errorf("failed to get sessions by status: %w", err))
	}
	return sessions, nil
}