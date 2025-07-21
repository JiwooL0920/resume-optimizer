package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/resume-optimizer/resume-processor/internal/handlers"
	"github.com/resume-optimizer/resume-processor/internal/middleware"
)

func main() {
	r := gin.Default()
	
	r.Use(middleware.CORS())
	
	v1 := r.Group("/api/v1")
	{
		resumes := v1.Group("/resumes")
		{
			resumes.POST("/upload", handlers.UploadResume)
			resumes.GET("/:id", handlers.GetResume)
			resumes.GET("/", handlers.ListResumes)
			resumes.DELETE("/:id", handlers.DeleteResume)
		}
		
		optimize := v1.Group("/optimize")
		{
			optimize.POST("/", handlers.OptimizeResume)
			optimize.POST("/feedback", handlers.ApplyFeedback)
		}
	}
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	
	log.Printf("Resume processor service starting on port %s", port)
	log.Fatal(r.Run(":" + port))
}