package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/resume-optimizer/auth-service/internal/handlers"
	"github.com/resume-optimizer/auth-service/internal/middleware"
	"github.com/resume-optimizer/auth-service/internal/config"
)

func main() {
	_ = config.Load()
	
	r := gin.Default()
	
	r.Use(middleware.CORS())
	
	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.GET("/google", handlers.GoogleAuth)
			auth.GET("/google/callback", handlers.GoogleCallback)
			auth.POST("/logout", handlers.Logout)
			auth.GET("/profile", middleware.RequireAuth(), handlers.GetProfile)
		}
	}
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	log.Printf("Auth service starting on port %s", port)
	log.Fatal(r.Run(":" + port))
}