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

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{db: db}
}

// Create creates a new user
func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	// Generate UUID if not already set
	if user.ID == "" {
		user.ID = uuid.New().String()
	}
	
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		if isUniqueConstraintError(err) {
			return errors.NewAppError(
				errors.ErrCodeDuplicate,
				"User with this email already exists",
				err,
			)
		}
		return errors.NewDatabaseError(fmt.Errorf("failed to create user: %w", err))
	}
	return nil
}

// GetByID retrieves a user by ID
func (r *userRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewAppError(
				errors.ErrCodeNotFound,
				"User not found",
				err,
			)
		}
		return nil, errors.NewDatabaseError(fmt.Errorf("failed to get user by ID: %w", err))
	}
	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).First(&user, "email = ?", email).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewAppError(
				errors.ErrCodeNotFound,
				"User not found",
				err,
			)
		}
		return nil, errors.NewDatabaseError(fmt.Errorf("failed to get user by email: %w", err))
	}
	return &user, nil
}

// GetByGoogleID retrieves a user by Google ID
func (r *userRepository) GetByGoogleID(ctx context.Context, googleID string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).First(&user, "google_id = ?", googleID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewAppError(
				errors.ErrCodeNotFound,
				"User not found",
				err,
			)
		}
		return nil, errors.NewDatabaseError(fmt.Errorf("failed to get user by Google ID: %w", err))
	}
	return &user, nil
}

// Update updates a user
func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		if isUniqueConstraintError(err) {
			return errors.NewAppError(
				errors.ErrCodeDuplicate,
				"User with this email already exists",
				err,
			)
		}
		return errors.NewDatabaseError(fmt.Errorf("failed to update user: %w", err))
	}
	return nil
}

// Delete deletes a user
func (r *userRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&models.User{}, "id = ?", id)
	if result.Error != nil {
		return errors.NewDatabaseError(fmt.Errorf("failed to delete user: %w", result.Error))
	}
	if result.RowsAffected == 0 {
		return errors.NewAppError(
			errors.ErrCodeNotFound,
			"User not found",
			nil,
		)
	}
	return nil
}

// List retrieves users with pagination
func (r *userRepository) List(ctx context.Context, offset, limit int) ([]*models.User, error) {
	var users []*models.User
	if err := r.db.WithContext(ctx).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&users).Error; err != nil {
		return nil, errors.NewDatabaseError(fmt.Errorf("failed to list users: %w", err))
	}
	return users, nil
}