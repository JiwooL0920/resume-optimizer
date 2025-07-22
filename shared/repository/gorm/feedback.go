package gorm

import (
	"context"
	"fmt"

	"github.com/resume-optimizer/shared/errors"
	"github.com/resume-optimizer/shared/models"
	"github.com/resume-optimizer/shared/repository"
	"gorm.io/gorm"
)

type feedbackRepository struct {
	db *gorm.DB
}

func NewFeedbackRepository(db *gorm.DB) repository.FeedbackRepository {
	return &feedbackRepository{db: db}
}

func (r *feedbackRepository) Create(ctx context.Context, feedback *models.Feedback) error {
	if err := r.db.WithContext(ctx).Create(feedback).Error; err != nil {
		return errors.NewDatabaseError(fmt.Errorf("failed to create feedback: %w", err))
	}
	return nil
}

func (r *feedbackRepository) GetByID(ctx context.Context, id string) (*models.Feedback, error) {
	var feedback models.Feedback
	if err := r.db.WithContext(ctx).First(&feedback, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewAppError(errors.ErrCodeNotFound, "Feedback not found", err)
		}
		return nil, errors.NewDatabaseError(fmt.Errorf("failed to get feedback: %w", err))
	}
	return &feedback, nil
}

func (r *feedbackRepository) GetBySessionID(ctx context.Context, sessionID string) ([]*models.Feedback, error) {
	var feedback []*models.Feedback
	if err := r.db.WithContext(ctx).Where("session_id = ?", sessionID).Order("created_at ASC").Find(&feedback).Error; err != nil {
		return nil, errors.NewDatabaseError(fmt.Errorf("failed to get feedback by session ID: %w", err))
	}
	return feedback, nil
}

func (r *feedbackRepository) Update(ctx context.Context, feedback *models.Feedback) error {
	if err := r.db.WithContext(ctx).Save(feedback).Error; err != nil {
		return errors.NewDatabaseError(fmt.Errorf("failed to update feedback: %w", err))
	}
	return nil
}

func (r *feedbackRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&models.Feedback{}, "id = ?", id)
	if result.Error != nil {
		return errors.NewDatabaseError(fmt.Errorf("failed to delete feedback: %w", result.Error))
	}
	if result.RowsAffected == 0 {
		return errors.NewAppError(errors.ErrCodeNotFound, "Feedback not found", nil)
	}
	return nil
}

func (r *feedbackRepository) MarkAsProcessed(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Model(&models.Feedback{}).Where("id = ?", id).Update("is_processed", true)
	if result.Error != nil {
		return errors.NewDatabaseError(fmt.Errorf("failed to mark feedback as processed: %w", result.Error))
	}
	if result.RowsAffected == 0 {
		return errors.NewAppError(errors.ErrCodeNotFound, "Feedback not found", nil)
	}
	return nil
}

func (r *feedbackRepository) GetUnprocessed(ctx context.Context) ([]*models.Feedback, error) {
	var feedback []*models.Feedback
	if err := r.db.WithContext(ctx).Where("is_processed = ?", false).Order("created_at ASC").Find(&feedback).Error; err != nil {
		return nil, errors.NewDatabaseError(fmt.Errorf("failed to get unprocessed feedback: %w", err))
	}
	return feedback, nil
}