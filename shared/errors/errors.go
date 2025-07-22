package errors

import (
	"fmt"
	"net/http"
)

// ErrorCode represents application error codes
type ErrorCode string

const (
	// Authentication & Authorization errors
	ErrCodeUnauthorized     ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden        ErrorCode = "FORBIDDEN"
	ErrCodeInvalidToken     ErrorCode = "INVALID_TOKEN"
	ErrCodeTokenExpired     ErrorCode = "TOKEN_EXPIRED"
	
	// Validation errors
	ErrCodeValidation       ErrorCode = "VALIDATION_ERROR"
	ErrCodeInvalidInput     ErrorCode = "INVALID_INPUT"
	ErrCodeMissingField     ErrorCode = "MISSING_FIELD"
	
	// Database errors
	ErrCodeDatabase         ErrorCode = "DATABASE_ERROR"
	ErrCodeNotFound         ErrorCode = "NOT_FOUND"
	ErrCodeDuplicate        ErrorCode = "DUPLICATE_ENTRY"
	
	// File errors
	ErrCodeFileUpload       ErrorCode = "FILE_UPLOAD_ERROR"
	ErrCodeFileType         ErrorCode = "INVALID_FILE_TYPE"
	ErrCodeFileSize         ErrorCode = "FILE_TOO_LARGE"
	
	// External service errors
	ErrCodeExternalService  ErrorCode = "EXTERNAL_SERVICE_ERROR"
	ErrCodeAIService        ErrorCode = "AI_SERVICE_ERROR"
	
	// Internal errors
	ErrCodeInternal         ErrorCode = "INTERNAL_ERROR"
	ErrCodeRateLimit        ErrorCode = "RATE_LIMIT_EXCEEDED"
)

// AppError represents an application error with context
type AppError struct {
	Code       ErrorCode `json:"code"`
	Message    string    `json:"message"`
	Details    string    `json:"details,omitempty"`
	HTTPStatus int       `json:"-"`
	Err        error     `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Err
}

// NewAppError creates a new application error
func NewAppError(code ErrorCode, message string, err error) *AppError {
	appErr := &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
	
	// Set default HTTP status based on error code
	appErr.HTTPStatus = getHTTPStatus(code)
	
	return appErr
}

// NewAppErrorWithDetails creates a new application error with details
func NewAppErrorWithDetails(code ErrorCode, message, details string, err error) *AppError {
	appErr := NewAppError(code, message, err)
	appErr.Details = details
	return appErr
}

// WithHTTPStatus sets a custom HTTP status code
func (e *AppError) WithHTTPStatus(status int) *AppError {
	e.HTTPStatus = status
	return e
}

// WithDetails adds details to the error
func (e *AppError) WithDetails(details string) *AppError {
	e.Details = details
	return e
}

// getHTTPStatus returns the appropriate HTTP status code for an error code
func getHTTPStatus(code ErrorCode) int {
	switch code {
	case ErrCodeUnauthorized, ErrCodeInvalidToken, ErrCodeTokenExpired:
		return http.StatusUnauthorized
	case ErrCodeForbidden:
		return http.StatusForbidden
	case ErrCodeValidation, ErrCodeInvalidInput, ErrCodeMissingField, ErrCodeFileType, ErrCodeFileSize:
		return http.StatusBadRequest
	case ErrCodeNotFound:
		return http.StatusNotFound
	case ErrCodeDuplicate:
		return http.StatusConflict
	case ErrCodeFileUpload:
		return http.StatusBadRequest
	case ErrCodeRateLimit:
		return http.StatusTooManyRequests
	case ErrCodeExternalService, ErrCodeAIService:
		return http.StatusBadGateway
	case ErrCodeDatabase, ErrCodeInternal:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// Predefined common errors
var (
	ErrUnauthorized = NewAppError(ErrCodeUnauthorized, "Unauthorized access", nil)
	ErrForbidden    = NewAppError(ErrCodeForbidden, "Forbidden", nil)
	ErrNotFound     = NewAppError(ErrCodeNotFound, "Resource not found", nil)
	ErrInternal     = NewAppError(ErrCodeInternal, "Internal server error", nil)
	ErrInvalidInput = NewAppError(ErrCodeInvalidInput, "Invalid input provided", nil)
)

// Validation error helpers
func NewValidationError(message string) *AppError {
	return NewAppError(ErrCodeValidation, message, nil)
}

func NewDatabaseError(err error) *AppError {
	return NewAppError(ErrCodeDatabase, "Database operation failed", err)
}

func NewFileUploadError(message string, err error) *AppError {
	return NewAppError(ErrCodeFileUpload, message, err)
}

func NewExternalServiceError(service string, err error) *AppError {
	return NewAppError(ErrCodeExternalService, fmt.Sprintf("External service error: %s", service), err)
}

// IsAppError checks if an error is an AppError
func IsAppError(err error) bool {
	_, ok := err.(*AppError)
	return ok
}

// GetAppError extracts AppError from error, or creates a generic one
func GetAppError(err error) *AppError {
	if appErr, ok := err.(*AppError); ok {
		return appErr
	}
	return NewAppError(ErrCodeInternal, "Internal server error", err)
}