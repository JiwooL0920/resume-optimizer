package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/resume-optimizer/shared/errors"
	"github.com/resume-optimizer/shared/middleware"
	"github.com/resume-optimizer/shared/models"
	"github.com/resume-optimizer/shared/repository"
	"github.com/resume-optimizer/shared/utils"
	"github.com/rs/zerolog/log"
)

// APIKeyHandlers handles API key-related requests
type APIKeyHandlers struct {
	encryptionService *utils.EncryptionService
	repoManager       repository.RepositoryManager
}

// NewAPIKeyHandlers creates a new API key handlers instance
func NewAPIKeyHandlers(encryptionService *utils.EncryptionService, repoManager repository.RepositoryManager) *APIKeyHandlers {
	return &APIKeyHandlers{
		encryptionService: encryptionService,
		repoManager:       repoManager,
	}
}

// CreateUserAPIKeyRequest represents the request body for creating an API key
type CreateUserAPIKeyRequest struct {
	Provider string `json:"provider" binding:"required"`
	APIKey   string `json:"api_key" binding:"required"`
}

// GetUserAPIKeys retrieves all API keys for the authenticated user
func (h *APIKeyHandlers) GetUserAPIKeys(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		appErr := errors.GetAppError(err)
		c.JSON(appErr.HTTPStatus, appErr)
		return
	}

	apiKeys, err := h.repoManager.UserAPIKey().GetByUserID(context.Background(), userID)
	if err != nil {
		appErr := errors.GetAppError(err)
		c.JSON(appErr.HTTPStatus, appErr)
		return
	}

	// Add masked keys for display
	for _, key := range apiKeys {
		key.MaskedKey = utils.MaskAPIKey(key.EncryptedKey)
	}

	c.JSON(http.StatusOK, gin.H{"api_keys": apiKeys})
}

// CreateUserAPIKey creates a new API key for the authenticated user
func (h *APIKeyHandlers) CreateUserAPIKey(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		appErr := errors.GetAppError(err)
		c.JSON(appErr.HTTPStatus, appErr)
		return
	}

	var req CreateUserAPIKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := errors.NewValidationError("Invalid request body")
		c.JSON(appErr.HTTPStatus, appErr)
		return
	}

	// Validate API key format
	if err := utils.ValidateAPIKey(req.Provider, req.APIKey); err != nil {
		appErr := errors.GetAppError(err)
		c.JSON(appErr.HTTPStatus, appErr)
		return
	}

	// Check if API key for this provider already exists
	ctx := context.Background()
	existingKey, err := h.repoManager.UserAPIKey().GetByUserIDAndProvider(ctx, userID, req.Provider)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); !ok || appErr.Code != errors.ErrCodeNotFound {
			appErr := errors.GetAppError(err)
			c.JSON(appErr.HTTPStatus, appErr)
			return
		}
	}

	if existingKey != nil {
		// Update existing key
		encryptedKey, err := h.encryptionService.Encrypt(req.APIKey)
		if err != nil {
			log.Error().Err(err).Msg("Failed to encrypt API key")
			appErr := errors.GetAppError(err)
			c.JSON(appErr.HTTPStatus, appErr)
			return
		}

		existingKey.EncryptedKey = encryptedKey
		if err := h.repoManager.UserAPIKey().Update(ctx, existingKey); err != nil {
			// Log the detailed error before wrapping
			log.Error().Err(err).
				Str("user_id", userID).
				Str("provider", req.Provider).
				Str("key_id", existingKey.ID).
				Msg("Failed to update API key in database")
			appErr := errors.GetAppError(err)
			c.JSON(appErr.HTTPStatus, appErr)
			return
		}

		existingKey.MaskedKey = utils.MaskAPIKey(req.APIKey)
		existingKey.EncryptedKey = "" // Don't return encrypted key

		log.Info().
			Str("user_id", userID).
			Str("provider", req.Provider).
			Msg("Updated user API key")

		c.JSON(http.StatusOK, gin.H{
			"message": "API key updated successfully",
			"api_key": existingKey,
		})
		return
	}

	// Create new API key
	encryptedKey, err := h.encryptionService.Encrypt(req.APIKey)
	if err != nil {
		log.Error().Err(err).Msg("Failed to encrypt API key")
		appErr := errors.GetAppError(err)
		c.JSON(appErr.HTTPStatus, appErr)
		return
	}

	apiKey := &models.UserAPIKey{
		UserID:       userID,
		Provider:     req.Provider,
		EncryptedKey: encryptedKey,
	}

	if err := h.repoManager.UserAPIKey().Create(ctx, apiKey); err != nil {
		// Log the detailed error before wrapping
		log.Error().Err(err).
			Str("user_id", userID).
			Str("provider", req.Provider).
			Msg("Failed to create API key in database")
		appErr := errors.GetAppError(err)
		c.JSON(appErr.HTTPStatus, appErr)
		return
	}

	// Prepare response
	apiKey.MaskedKey = utils.MaskAPIKey(req.APIKey)
	apiKey.EncryptedKey = "" // Don't return encrypted key

	log.Info().
		Str("user_id", userID).
		Str("provider", req.Provider).
		Str("key_id", apiKey.ID).
		Msg("Created user API key")

	c.JSON(http.StatusCreated, gin.H{
		"message": "API key created successfully",
		"api_key": apiKey,
	})
}

// DeleteUserAPIKey deletes an API key for the authenticated user
func (h *APIKeyHandlers) DeleteUserAPIKey(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		appErr := errors.GetAppError(err)
		c.JSON(appErr.HTTPStatus, appErr)
		return
	}

	keyID := c.Param("id")
	if keyID == "" {
		appErr := errors.NewValidationError("API key ID is required")
		c.JSON(appErr.HTTPStatus, appErr)
		return
	}

	ctx := context.Background()

	// Verify the API key belongs to the user
	apiKey, err := h.repoManager.UserAPIKey().GetByID(ctx, keyID)
	if err != nil {
		appErr := errors.GetAppError(err)
		c.JSON(appErr.HTTPStatus, appErr)
		return
	}

	if apiKey.UserID != userID {
		appErr := errors.ErrForbidden.WithDetails("API key does not belong to the authenticated user")
		c.JSON(appErr.HTTPStatus, appErr)
		return
	}

	// Delete the API key
	if err := h.repoManager.UserAPIKey().Delete(ctx, keyID); err != nil {
		appErr := errors.GetAppError(err)
		c.JSON(appErr.HTTPStatus, appErr)
		return
	}

	log.Info().
		Str("user_id", userID).
		Str("key_id", keyID).
		Str("provider", apiKey.Provider).
		Msg("Deleted user API key")

	c.JSON(http.StatusOK, gin.H{
		"message": "API key deleted successfully",
		"id":      keyID,
	})
}