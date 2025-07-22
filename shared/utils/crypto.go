package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/resume-optimizer/shared/errors"
)

// EncryptionService handles encryption and decryption operations
type EncryptionService struct {
	key []byte
}

// NewEncryptionService creates a new encryption service
func NewEncryptionService(key string) (*EncryptionService, error) {
	if len(key) != 32 {
		return nil, errors.NewAppError(
			errors.ErrCodeValidation,
			"Encryption key must be exactly 32 characters long",
			nil,
		)
	}

	return &EncryptionService{
		key: []byte(key),
	}, nil
}

// Encrypt encrypts plaintext using AES-GCM
func (e *EncryptionService) Encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", errors.NewAppError(
			errors.ErrCodeInternal,
			"Failed to create cipher",
			err,
		)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", errors.NewAppError(
			errors.ErrCodeInternal,
			"Failed to create GCM",
			err,
		)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", errors.NewAppError(
			errors.ErrCodeInternal,
			"Failed to generate nonce",
			err,
		)
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return hex.EncodeToString(ciphertext), nil
}

// Decrypt decrypts ciphertext using AES-GCM
func (e *EncryptionService) Decrypt(ciphertext string) (string, error) {
	data, err := hex.DecodeString(ciphertext)
	if err != nil {
		return "", errors.NewAppError(
			errors.ErrCodeValidation,
			"Invalid ciphertext format",
			err,
		)
	}

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", errors.NewAppError(
			errors.ErrCodeInternal,
			"Failed to create cipher",
			err,
		)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", errors.NewAppError(
			errors.ErrCodeInternal,
			"Failed to create GCM",
			err,
		)
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.NewAppError(
			errors.ErrCodeValidation,
			"Ciphertext too short",
			nil,
		)
	}

	nonce, ciphertext_bytes := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext_bytes, nil)
	if err != nil {
		return "", errors.NewAppError(
			errors.ErrCodeInternal,
			"Failed to decrypt",
			err,
		)
	}

	return string(plaintext), nil
}

// MaskAPIKey masks an API key for safe display
func MaskAPIKey(apiKey string) string {
	if len(apiKey) <= 8 {
		return "****"
	}
	
	prefix := apiKey[:4]
	suffix := apiKey[len(apiKey)-4:]
	return fmt.Sprintf("%s****%s", prefix, suffix)
}

// ValidateAPIKey validates the format of an API key based on provider
func ValidateAPIKey(provider, apiKey string) error {
	if apiKey == "" {
		return errors.NewValidationError("API key cannot be empty")
	}

	switch provider {
	case "openai":
		if !isValidOpenAIKey(apiKey) {
			return errors.NewValidationError("Invalid OpenAI API key format")
		}
	case "anthropic":
		if !isValidAnthropicKey(apiKey) {
			return errors.NewValidationError("Invalid Anthropic API key format")
		}
	case "google":
		if len(apiKey) < 10 {
			return errors.NewValidationError("Invalid Google API key format")
		}
	default:
		if len(apiKey) < 10 {
			return errors.NewValidationError("API key too short")
		}
	}

	return nil
}

// isValidOpenAIKey checks if the API key matches OpenAI format
func isValidOpenAIKey(key string) bool {
	if len(key) < 50 {
		return false
	}
	return (len(key) >= 3 && key[:3] == "sk-") ||
		(len(key) >= 7 && key[:7] == "sk-proj")
}

// isValidAnthropicKey checks if the API key matches Anthropic format
func isValidAnthropicKey(key string) bool {
	return len(key) >= 50 && len(key) >= 7 && key[:7] == "sk-ant-"
}

// GenerateSecureKey generates a cryptographically secure random key
func GenerateSecureKey(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", errors.NewAppError(
			errors.ErrCodeInternal,
			"Failed to generate secure key",
			err,
		)
	}
	return hex.EncodeToString(bytes)[:length], nil
}