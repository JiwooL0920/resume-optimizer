package gorm

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/resume-optimizer/shared/errors"
	"github.com/resume-optimizer/shared/models"
	"github.com/resume-optimizer/shared/repository"
	"gorm.io/gorm"
)

type userAPIKeyRepository struct {
	db *gorm.DB
}

func NewUserAPIKeyRepository(db *gorm.DB) repository.UserAPIKeyRepository {
	return &userAPIKeyRepository{db: db}
}

func (r *userAPIKeyRepository) Create(ctx context.Context, apiKey *models.UserAPIKey) error {
	// Generate UUID if not already set
	if apiKey.ID == "" {
		apiKey.ID = uuid.New().String()
	}
	
	if err := r.db.WithContext(ctx).Create(apiKey).Error; err != nil {
		if isUniqueConstraintError(err) {
			return errors.NewAppError(errors.ErrCodeDuplicate, "API key for this provider already exists", err)
		}
		return errors.NewDatabaseError(fmt.Errorf("failed to create API key: %w", err))
	}
	return nil
}

func (r *userAPIKeyRepository) GetByID(ctx context.Context, id string) (*models.UserAPIKey, error) {
	var apiKey models.UserAPIKey
	if err := r.db.WithContext(ctx).First(&apiKey, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewAppError(errors.ErrCodeNotFound, "API key not found", err)
		}
		return nil, errors.NewDatabaseError(fmt.Errorf("failed to get API key: %w", err))
	}
	return &apiKey, nil
}

func (r *userAPIKeyRepository) GetByUserID(ctx context.Context, userID string) ([]*models.UserAPIKey, error) {
	var apiKeys []*models.UserAPIKey
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&apiKeys).Error; err != nil {
		return nil, errors.NewDatabaseError(fmt.Errorf("failed to get API keys: %w", err))
	}
	return apiKeys, nil
}

func (r *userAPIKeyRepository) GetByUserIDAndProvider(ctx context.Context, userID, provider string) (*models.UserAPIKey, error) {
	var apiKey models.UserAPIKey
	if err := r.db.WithContext(ctx).First(&apiKey, "user_id = ? AND provider = ?", userID, provider).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewAppError(errors.ErrCodeNotFound, "API key not found", err)
		}
		return nil, errors.NewDatabaseError(fmt.Errorf("failed to get API key: %w", err))
	}
	return &apiKey, nil
}

func (r *userAPIKeyRepository) Update(ctx context.Context, apiKey *models.UserAPIKey) error {
	if err := r.db.WithContext(ctx).Save(apiKey).Error; err != nil {
		return errors.NewDatabaseError(fmt.Errorf("failed to update API key: %w", err))
	}
	return nil
}

func (r *userAPIKeyRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&models.UserAPIKey{}, "id = ?", id)
	if result.Error != nil {
		return errors.NewDatabaseError(fmt.Errorf("failed to delete API key: %w", result.Error))
	}
	if result.RowsAffected == 0 {
		return errors.NewAppError(errors.ErrCodeNotFound, "API key not found", nil)
	}
	return nil
}

func (r *userAPIKeyRepository) DeleteByUserIDAndProvider(ctx context.Context, userID, provider string) error {
	result := r.db.WithContext(ctx).Delete(&models.UserAPIKey{}, "user_id = ? AND provider = ?", userID, provider)
	if result.Error != nil {
		return errors.NewDatabaseError(fmt.Errorf("failed to delete API key: %w", result.Error))
	}
	return nil
}

func (r *userAPIKeyRepository) List(ctx context.Context, offset, limit int) ([]*models.UserAPIKey, error) {
	var apiKeys []*models.UserAPIKey
	query := applyPagination(r.db.WithContext(ctx), offset, limit)
	if err := query.Find(&apiKeys).Error; err != nil {
		return nil, errors.NewDatabaseError(fmt.Errorf("failed to list API keys: %w", err))
	}
	return apiKeys, nil
}