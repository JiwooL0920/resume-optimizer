package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/resume-optimizer/shared/config"
	"github.com/resume-optimizer/shared/errors"
	"github.com/resume-optimizer/shared/middleware"
	"github.com/resume-optimizer/shared/models"
	"github.com/resume-optimizer/shared/repository"
	"github.com/resume-optimizer/shared/utils"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// AuthHandlers handles authentication-related requests
type AuthHandlers struct {
	config      *config.Config
	jwtService  *utils.JWTService
	repoManager repository.RepositoryManager
	oauthConfig *oauth2.Config
}

// NewAuthHandlers creates a new auth handlers instance
func NewAuthHandlers(cfg *config.Config, jwtService *utils.JWTService, repoManager repository.RepositoryManager) *AuthHandlers {
	oauthConfig := &oauth2.Config{
		ClientID:     cfg.Auth.GoogleClientID,
		ClientSecret: cfg.Auth.GoogleClientSecret,
		RedirectURL:  cfg.Auth.GoogleRedirectURL,
		Scopes:       []string{"openid", "profile", "email"},
		Endpoint:     google.Endpoint,
	}

	return &AuthHandlers{
		config:      cfg,
		jwtService:  jwtService,
		repoManager: repoManager,
		oauthConfig: oauthConfig,
	}
}

// GoogleAuth initiates Google OAuth flow
func (h *AuthHandlers) GoogleAuth(c *gin.Context) {
	state := uuid.New().String()
	// In production, store state in session or cache for validation
	url := h.oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)

	log.Info().
		Str("state", state).
		Msg("Initiating Google OAuth flow")

	c.JSON(http.StatusOK, gin.H{
		"auth_url": url,
		"state":    state,
	})
}

// GoogleCallback handles Google OAuth callback
func (h *AuthHandlers) GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	if code == "" {
		appErr := errors.NewValidationError("Authorization code is required")
		c.JSON(appErr.HTTPStatus, appErr)
		return
	}

	// TODO: Validate state parameter in production
	// For now, just log it
	log.Debug().Str("state", state).Msg("OAuth callback received")

	// Exchange code for token
	token, err := h.oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Error().Err(err).Msg("Failed to exchange OAuth code for token")
		appErr := errors.NewExternalServiceError("Google OAuth", err)
		c.JSON(appErr.HTTPStatus, appErr)
		return
	}

	// Get user info from Google
	client := h.oauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user info from Google")
		appErr := errors.NewExternalServiceError("Google API", err)
		c.JSON(appErr.HTTPStatus, appErr)
		return
	}
	defer resp.Body.Close()

	var googleUser struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		log.Error().Err(err).Msg("Failed to decode Google user info")
		appErr := errors.NewExternalServiceError("Google API", err)
		c.JSON(appErr.HTTPStatus, appErr)
		return
	}

	ctx := context.Background()

	// Check if user exists
	user, err := h.repoManager.User().GetByGoogleID(ctx, googleUser.ID)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok && appErr.Code == errors.ErrCodeNotFound {
			// Check if user exists by email
			existingUser, emailErr := h.repoManager.User().GetByEmail(ctx, googleUser.Email)
			if emailErr != nil {
				if appErr2, ok := emailErr.(*errors.AppError); ok && appErr2.Code == errors.ErrCodeNotFound {
					// Create new user
					user = &models.User{
						Email:      googleUser.Email,
						GoogleID:   &googleUser.ID,
						Name:       googleUser.Name,
						PictureURL: &googleUser.Picture,
					}

					if err := h.repoManager.User().Create(ctx, user); err != nil {
						log.Error().Err(err).Msg("Failed to create user")
						appErr := errors.GetAppError(err)
						c.JSON(appErr.HTTPStatus, appErr)
						return
					}

					log.Info().
						Str("user_id", user.ID).
						Str("email", user.Email).
						Msg("Created new user")
				} else {
					log.Error().Err(emailErr).Msg("Database error checking user by email")
					appErr := errors.GetAppError(emailErr)
					c.JSON(appErr.HTTPStatus, appErr)
					return
				}
			} else {
				// User exists with email but no Google ID, link accounts
				existingUser.GoogleID = &googleUser.ID
				existingUser.PictureURL = &googleUser.Picture
				if err := h.repoManager.User().Update(ctx, existingUser); err != nil {
					log.Error().Err(err).Msg("Failed to link Google account")
					appErr := errors.GetAppError(err)
					c.JSON(appErr.HTTPStatus, appErr)
					return
				}
				user = existingUser

				log.Info().
					Str("user_id", user.ID).
					Str("email", user.Email).
					Msg("Linked existing user with Google account")
			}
		} else {
			log.Error().Err(err).Msg("Database error getting user by Google ID")
			appErr := errors.GetAppError(err)
			c.JSON(appErr.HTTPStatus, appErr)
			return
		}
	}

	// Generate JWT token
	jwtToken, err := h.jwtService.GenerateToken(user.ID, user.Email)
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate JWT token")
		appErr := errors.GetAppError(err)
		c.JSON(appErr.HTTPStatus, appErr)
		return
	}

	log.Info().
		Str("user_id", user.ID).
		Str("email", user.Email).
		Msg("User authenticated successfully")

	// Redirect to frontend callback with token
	frontendCallbackURL := fmt.Sprintf("%s/auth/callback?token=%s", h.config.Client.BaseURL, jwtToken)
	c.Redirect(http.StatusFound, frontendCallbackURL)
}

// GetProfile returns the current user's profile
func (h *AuthHandlers) GetProfile(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		appErr := errors.GetAppError(err)
		c.JSON(appErr.HTTPStatus, appErr)
		return
	}

	user, err := h.repoManager.User().GetByID(context.Background(), userID)
	if err != nil {
		appErr := errors.GetAppError(err)
		c.JSON(appErr.HTTPStatus, appErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// Logout handles user logout
func (h *AuthHandlers) Logout(c *gin.Context) {
	// In a stateless JWT system, logout is typically handled client-side
	// by removing the token from storage
	// For additional security, you could maintain a blacklist of tokens

	log.Info().Msg("User logged out")
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}