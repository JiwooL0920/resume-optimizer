package middleware

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/resume-optimizer/auth-service/internal/utils"
)

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("Auth middleware called for %s %s", c.Request.Method, c.Request.URL.Path)
		
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Printf("No authorization header provided")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "No authorization token provided",
			})
			c.Abort()
			return
		}

		log.Printf("Authorization header: %s", authHeader)

		// Extract token from "Bearer token" format
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			log.Printf("Invalid authorization header format")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization header format",
			})
			c.Abort()
			return
		}

		tokenPreview := tokenString
		if len(tokenString) > 20 {
			tokenPreview = tokenString[:20] + "..."
		}
		log.Printf("Extracted token: %s", tokenPreview)

		jwtSecret := os.Getenv("JWT_SECRET")
		claims, err := utils.ValidateJWT(tokenString, jwtSecret)
		if err != nil {
			log.Printf("JWT validation failed: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token: " + err.Error(),
			})
			c.Abort()
			return
		}

		log.Printf("JWT validation successful for user: %s", claims.Email)

		// Set user information in context
		c.Set("userID", claims.UserID)
		c.Set("userEmail", claims.Email)
		c.Next()
	}
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}