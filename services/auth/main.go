package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/resume-optimizer/auth-service/internal/handlers"
	"github.com/resume-optimizer/shared/config"
	"github.com/resume-optimizer/shared/database"
	"github.com/resume-optimizer/shared/middleware"
	"github.com/resume-optimizer/shared/repository/gorm"
	"github.com/resume-optimizer/shared/utils"
	"github.com/rs/zerolog/log"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Setup logging
	middleware.SetupLogger(cfg.Log.Level, cfg.Log.Format)

	// Connect to database
	if err := database.Connect(&cfg.Database); err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer database.Disconnect()

	// Initialize services
	jwtService := utils.NewJWTService(cfg.Auth.JWTSecret, cfg.Auth.JWTExpiration)
	encryptionService, err := utils.NewEncryptionService(cfg.Security.EncryptionKey)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize encryption service")
	}

	// Initialize repositories
	repoManager := gorm.NewRepositoryManager(database.GetDB())

	// Initialize handlers
	authHandlers := handlers.NewAuthHandlers(cfg, jwtService, repoManager)
	apiKeyHandlers := handlers.NewAPIKeyHandlers(encryptionService, repoManager)

	// Setup Gin
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()

	// Middleware
	r.Use(middleware.RequestIDMiddleware())
	r.Use(middleware.LoggingMiddleware())
	r.Use(middleware.CORSWithConfig(cfg.Security.AllowedOrigins))
	r.Use(middleware.RateLimitMiddleware(cfg.Security.RateLimitRPS, cfg.Security.RateLimitRPS*2))

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		if err := database.HealthCheck(); err != nil {
			c.JSON(503, gin.H{
				"status":  "unhealthy",
				"service": "auth",
				"error":   "database connection failed",
			})
			return
		}
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "auth",
		})
	})

	// API routes
	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.GET("/google", authHandlers.GoogleAuth)
			auth.GET("/callback", authHandlers.GoogleCallback)
			auth.POST("/logout", authHandlers.Logout)
			auth.GET("/profile", middleware.AuthMiddleware(jwtService), authHandlers.GetProfile)
		}

		user := v1.Group("/user")
		user.Use(middleware.AuthMiddleware(jwtService))
		{
			user.GET("/api-keys", apiKeyHandlers.GetUserAPIKeys)
			user.POST("/api-keys", apiKeyHandlers.CreateUserAPIKey)
			user.DELETE("/api-keys/:id", apiKeyHandlers.DeleteUserAPIKey)
		}
	}

	log.Info().
		Str("host", cfg.Server.Host).
		Int("port", cfg.Server.Port).
		Str("mode", cfg.Server.Mode).
		Msg("Starting auth service")

	if err := r.Run(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}