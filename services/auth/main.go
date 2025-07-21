package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/resume-optimizer/auth-service/internal/handlers"
	"github.com/resume-optimizer/auth-service/internal/middleware"
	"github.com/resume-optimizer/auth-service/internal/config"
	"github.com/resume-optimizer/auth-service/internal/database"
)

func main() {
	cfg := config.Load()
	database.InitDatabase(cfg)
	handlers.InitHandlers() // Initialize handlers after database is ready
	
	r := gin.Default()
	
	r.Use(middleware.CORS())
	
	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy", "service": "auth"})
	})
	
	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.GET("/google", handlers.GoogleAuth)
			auth.GET("/google/callback", handlers.GoogleCallback)
			auth.POST("/logout", handlers.Logout)
			auth.GET("/profile", middleware.RequireAuth(), handlers.GetProfile)
		}
		
		user := v1.Group("/user")
		user.Use(middleware.RequireAuth())
		{
			user.GET("/api-keys", handlers.GetUserAPIKeys)
			user.POST("/api-keys", handlers.CreateUserAPIKey)
			user.DELETE("/api-keys/:id", handlers.DeleteUserAPIKey)
		}
	}
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	log.Printf("Auth service starting on port %s", port)
	log.Fatal(r.Run(":" + port))
}