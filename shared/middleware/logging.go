package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// LoggingMiddleware creates a structured logging middleware
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get user ID if available
		userID, _ := c.Get("userID")
		userIDStr := ""
		if userID != nil {
			userIDStr = userID.(string)
		}

		// Create log event
		logEvent := log.Info()
		
	// Add error information if request failed
	if len(c.Errors) > 0 {
		errorStrings := make([]string, len(c.Errors))
		for i, err := range c.Errors {
			errorStrings[i] = err.Error()
		}
		logEvent = log.Error().Strs("errors", errorStrings)
	} else if c.Writer.Status() >= 400 {
		logEvent = log.Warn()
	}

		// Build full path with query parameters
		fullPath := path
		if raw != "" {
			fullPath = path + "?" + raw
		}

		// Log the request
		logEvent.
			Str("method", c.Request.Method).
			Str("path", fullPath).
			Str("ip", c.ClientIP()).
			Str("user_agent", c.Request.UserAgent()).
			Str("user_id", userIDStr).
			Int("status", c.Writer.Status()).
			Int64("size", int64(c.Writer.Size())).
			Dur("latency", latency).
			Msg("HTTP Request")
	}
}

// SetupLogger configures the global logger
func SetupLogger(level string, format string) {
	// Set log level
	switch level {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// Set log format
	if format == "console" {
		log.Logger = log.Output(zerolog.ConsoleWriter{
			TimeFormat: time.RFC3339,
		})
	}
	// Default is JSON format, no additional setup needed
}

// RequestIDMiddleware adds a unique request ID to each request
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			// Generate a simple request ID
			requestID = generateRequestID()
		}
		
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		
		// Add request ID to logger context
		logger := log.With().Str("request_id", requestID).Logger()
		c.Set("logger", &logger)
		
		c.Next()
	}
}

// generateRequestID generates a simple request ID
func generateRequestID() string {
	// Use current timestamp + random component for simplicity
	return time.Now().Format("20060102-150405") + "-" + randomString(6)
}

// randomString generates a random string of given length
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}