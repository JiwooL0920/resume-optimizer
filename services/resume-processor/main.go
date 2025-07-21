package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/resume-optimizer/resume-processor/internal/config"
	"github.com/resume-optimizer/resume-processor/internal/handlers"
	"github.com/resume-optimizer/resume-processor/internal/middleware"
	"github.com/resume-optimizer/resume-processor/internal/database"
)

func main() {
	cfg := config.Load()
	database.InitDatabase(cfg.DatabaseURL)

	r := gin.Default()
	
	r.Use(middleware.CORS())
	
	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy", "service": "resume-processor"})
	})
	
	v1 := r.Group("/api/v1")
	{
		resumes := v1.Group("/resumes")
		resumes.Use(middleware.RequireAuth())
		{
			resumes.POST("/upload", handlers.UploadResume)
			resumes.GET("/:id", handlers.GetResume)
			resumes.GET("/", handlers.ListResumes)
			resumes.DELETE("/:id", handlers.DeleteResume)
		}
		
		optimize := v1.Group("/optimize")
		optimize.Use(middleware.RequireAuth())
		{
			optimize.POST("/", handlers.OptimizeResume)
			optimize.POST("/feedback", handlers.ApplyFeedback)
		}
	}
	
	log.Printf("Resume processor service starting on port %s", cfg.Port)
	log.Fatal(r.Run(":" + cfg.Port))
}