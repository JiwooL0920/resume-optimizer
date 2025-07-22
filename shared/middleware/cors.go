package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposedHeaders   []string
	AllowCredentials bool
	MaxAge           int
}

// DefaultCORSConfig returns a default CORS configuration
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowedHeaders: []string{
			"Origin",
			"Content-Length",
			"Content-Type",
			"Authorization",
			"Accept",
			"X-Requested-With",
		},
		ExposedHeaders:   []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60, // 12 hours
	}
}

// CORSMiddleware creates CORS middleware with the given configuration
func CORSMiddleware(config CORSConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		
		// Check if origin is allowed
		if isOriginAllowed(origin, config.AllowedOrigins) {
			c.Header("Access-Control-Allow-Origin", origin)
		}
		
		c.Header("Access-Control-Allow-Methods", strings.Join(config.AllowedMethods, ", "))
		c.Header("Access-Control-Allow-Headers", strings.Join(config.AllowedHeaders, ", "))
		
		if len(config.ExposedHeaders) > 0 {
			c.Header("Access-Control-Expose-Headers", strings.Join(config.ExposedHeaders, ", "))
		}
		
		if config.AllowCredentials {
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		
		if config.MaxAge > 0 {
			c.Header("Access-Control-Max-Age", string(rune(config.MaxAge)))
		}

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// isOriginAllowed checks if the origin is in the allowed list
func isOriginAllowed(origin string, allowedOrigins []string) bool {
	if origin == "" {
		return false
	}
	
	for _, allowed := range allowedOrigins {
		if allowed == "*" {
			return true
		}
		if allowed == origin {
			return true
		}
		// Support wildcard subdomains (e.g., *.example.com)
		if strings.HasPrefix(allowed, "*.") {
			domain := allowed[2:]
			if strings.HasSuffix(origin, "."+domain) || origin == domain {
				return true
			}
		}
	}
	
	return false
}

// CORSWithConfig creates CORS middleware with custom configuration
func CORSWithConfig(allowedOrigins []string) gin.HandlerFunc {
	config := DefaultCORSConfig()
	config.AllowedOrigins = allowedOrigins
	return CORSMiddleware(config)
}