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
