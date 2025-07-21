package services

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"github.com/resume-optimizer/auth-service/internal/database"
	"github.com/resume-optimizer/auth-service/internal/models"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService() *UserService {
	return &UserService{
		db: database.GetDB(),
	}
}

func (s *UserService) CreateUser(email, name string, googleID *string, pictureURL *string) (*models.User, error) {
	user := &models.User{
		ID:         uuid.New().String(),
		Email:      email,
		Name:       name,
		GoogleID:   googleID,
		PictureURL: pictureURL,
	}

	if err := s.db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (s *UserService) GetUserByGoogleID(googleID string) (*models.User, error) {
	var user models.User
	if err := s.db.Where("google_id = ?", googleID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (s *UserService) GetUserByID(userID string) (*models.User, error) {
	var user models.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (s *UserService) UpdateUser(userID string, updates map[string]interface{}) error {
	return s.db.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error
}

func (s *UserService) GetUserAPIKeys(userID string) ([]*models.UserAPIKey, error) {
	var apiKeys []*models.UserAPIKey
	if err := s.db.Where("user_id = ?", userID).Find(&apiKeys).Error; err != nil {
		return nil, err
	}
	return apiKeys, nil
}

func (s *UserService) CreateUserAPIKey(apiKey *models.UserAPIKey) error {
	return s.db.Create(apiKey).Error
}

func (s *UserService) DeleteUserAPIKey(userID, keyID string) error {
	return s.db.Where("user_id = ? AND id = ?", userID, keyID).Delete(&models.UserAPIKey{}).Error
}

func (s *UserService) GetUserAPIKey(userID, keyID string) (*models.UserAPIKey, error) {
	var apiKey models.UserAPIKey
	if err := s.db.Where("user_id = ? AND id = ?", userID, keyID).First(&apiKey).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &apiKey, nil
}
