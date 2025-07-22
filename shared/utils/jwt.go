package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/resume-optimizer/shared/errors"
)

// JWTService handles JWT operations
type JWTService struct {
	secret     []byte
	expiration time.Duration
}

// Claims represents JWT claims
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// NewJWTService creates a new JWT service
func NewJWTService(secret string, expirationSeconds int) *JWTService {
	return &JWTService{
		secret:     []byte(secret),
		expiration: time.Duration(expirationSeconds) * time.Second,
	}
}

// GenerateToken generates a JWT token for a user
func (j *JWTService) GenerateToken(userID, email string) (string, error) {
	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "resume-optimizer",
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(j.secret)
	if err != nil {
		return "", errors.NewAppError(
			errors.ErrCodeInternal,
			"Failed to generate token",
			err,
		)
	}

	return tokenString, nil
}

// ValidateToken validates and parses a JWT token
func (j *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secret, nil
	})

	if err != nil {
		// Check for specific error types in JWT v5
		if err.Error() == "token is expired" {
			return nil, errors.NewAppError(
				errors.ErrCodeTokenExpired,
				"Token has expired",
				err,
			)
		}
		return nil, errors.NewAppError(
			errors.ErrCodeInvalidToken,
			"Failed to parse token",
			err,
		)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.NewAppError(
		errors.ErrCodeInvalidToken,
		"Invalid token claims",
		nil,
	)
}

// RefreshToken generates a new token with extended expiration
func (j *JWTService) RefreshToken(tokenString string) (string, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		// Allow refresh even if token is expired (but not for other validation errors)
		if appErr, ok := err.(*errors.AppError); ok && appErr.Code == errors.ErrCodeTokenExpired {
			// Try to parse the expired token to get user info
			token, parseErr := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
				return j.secret, nil
			})
			
			if parseErr != nil {
				return "", err // Return original expiration error
			}
			
			if expiredClaims, ok := token.Claims.(*Claims); ok {
				claims = expiredClaims
			} else {
				return "", err
			}
		} else {
			return "", err
		}
	}

	// Generate new token
	return j.GenerateToken(claims.UserID, claims.Email)
}

// GetTokenExpiration returns the token expiration duration
func (j *JWTService) GetTokenExpiration() time.Duration {
	return j.expiration
}