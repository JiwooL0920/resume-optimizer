package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/resume-optimizer/shared/errors"
	"github.com/resume-optimizer/shared/utils"
	"github.com/rs/zerolog/log"
)

// AuthMiddleware creates authentication middleware
func AuthMiddleware(jwtService *utils.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Warn().
				Str("path", c.FullPath()).
				Str("method", c.Request.Method).
				Msg("Missing authorization header")
			
			appErr := errors.ErrUnauthorized.WithDetails("Authorization header required")
			c.JSON(appErr.HTTPStatus, appErr)
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			log.Warn().
				Str("auth_header", authHeader).
				Msg("Invalid authorization header format")
			
			appErr := errors.ErrUnauthorized.WithDetails("Invalid authorization header format")
			c.JSON(appErr.HTTPStatus, appErr)
			c.Abort()
			return
		}

		token := tokenParts[1]
		claims, err := jwtService.ValidateToken(token)
		if err != nil {
			log.Warn().
				Err(err).
				Str("token", maskToken(token)).
				Msg("Token validation failed")
			
			appErr := errors.GetAppError(err)
			c.JSON(appErr.HTTPStatus, appErr)
			c.Abort()
			return
		}

		// Store user information in context
		c.Set("userID", claims.UserID)
		c.Set("userEmail", claims.Email)
		c.Set("claims", claims)

		log.Debug().
			Str("user_id", claims.UserID).
			Str("email", claims.Email).
			Str("path", c.FullPath()).
			Msg("Request authenticated")

		c.Next()
	}
}

// OptionalAuthMiddleware creates optional authentication middleware
// Sets user context if token is valid, but doesn't block requests without tokens
func OptionalAuthMiddleware(jwtService *utils.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.Next()
			return
		}

		token := tokenParts[1]
		claims, err := jwtService.ValidateToken(token)
		if err != nil {
			log.Debug().
				Err(err).
				Str("token", maskToken(token)).
				Msg("Optional auth token validation failed")
			c.Next()
			return
		}

		// Store user information in context
		c.Set("userID", claims.UserID)
		c.Set("userEmail", claims.Email)
		c.Set("claims", claims)

		c.Next()
	}
}

// RequireRole creates role-based authorization middleware
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := c.Get("claims")
		if !exists {
			appErr := errors.ErrUnauthorized.WithDetails("User not authenticated")
			c.JSON(appErr.HTTPStatus, appErr)
			c.Abort()
			return
		}

		userClaims, ok := claims.(*utils.Claims)
		if !ok {
			appErr := errors.ErrInternal.WithDetails("Invalid user claims")
			c.JSON(appErr.HTTPStatus, appErr)
			c.Abort()
			return
		}

		// TODO: Add role checking when user roles are implemented
		// For now, all authenticated users are allowed
		_ = userClaims
		_ = allowedRoles

		c.Next()
	}
}

// GetUserID extracts user ID from context
func GetUserID(c *gin.Context) (string, error) {
	userID, exists := c.Get("userID")
	if !exists {
		return "", errors.ErrUnauthorized
	}

	userIDStr, ok := userID.(string)
	if !ok {
		return "", errors.ErrInternal.WithDetails("Invalid user ID in context")
	}

	return userIDStr, nil
}

// GetUserEmail extracts user email from context
func GetUserEmail(c *gin.Context) (string, error) {
	email, exists := c.Get("userEmail")
	if !exists {
		return "", errors.ErrUnauthorized
	}

	emailStr, ok := email.(string)
	if !ok {
		return "", errors.ErrInternal.WithDetails("Invalid user email in context")
	}

	return emailStr, nil
}

// GetClaims extracts full claims from context
func GetClaims(c *gin.Context) (*utils.Claims, error) {
	claims, exists := c.Get("claims")
	if !exists {
		return nil, errors.ErrUnauthorized
	}

	userClaims, ok := claims.(*utils.Claims)
	if !ok {
		return nil, errors.ErrInternal.WithDetails("Invalid user claims")
	}

	return userClaims, nil
}

// maskToken masks a JWT token for logging
func maskToken(token string) string {
	if len(token) <= 8 {
		return "****"
	}
	return token[:4] + "..." + token[len(token)-4:]
}