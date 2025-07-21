package handlers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/resume-optimizer/resume-processor/internal/database"
	"github.com/resume-optimizer/resume-processor/internal/models"
	"github.com/resume-optimizer/resume-processor/internal/services"
	"gorm.io/gorm"
)

var (
	storagePath = "./uploaded_files"
)

func init() {
	if os.MkdirAll(storagePath, 0755) != nil {
		panic("Unable to create storage directory")
	}
}

// UploadResume uploads a new resume to the server
func UploadResume(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request: " + err.Error()})
		return
	}

	fileID := uuid.New().String()
	destPath := filepath.Join(storagePath, fileID)
	if err := c.SaveUploadedFile(file, destPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save file: " + err.Error()})
		return
	}

	fileSize := int(file.Size)
	userIDStr := userID.(string)
	resume := models.Resume{
		ID:              fileID,
		UserID:          &userIDStr,
		Title:           file.Filename,
		OriginalContent: destPath,
		FileType:        filepath.Ext(file.Filename),
		FileSize:        &fileSize,
	}

	if err := database.GetDB().Create(&resume).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, resume)
}

// GetResume retrieves a resume by ID for the authenticated user
func GetResume(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	id := c.Param("id")
	var resume models.Resume

	if err := database.GetDB().Where("id = ? AND user_id = ?", id, userID.(string)).First(&resume).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Resume not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"resume": resume})
}

// ListResumes lists user's resumes
func ListResumes(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var resumes []models.Resume
	if err := database.GetDB().Where("user_id = ?", userID.(string)).Find(&resumes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"resumes": resumes})
}

// DeleteResume deletes a resume by ID for the authenticated user
func DeleteResume(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	id := c.Param("id")
	if err := database.GetDB().Where("id = ? AND user_id = ?", id, userID.(string)).Delete(&models.Resume{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Resume deleted", "id": id})
}

// OptimizeResume optimizes a resume
func OptimizeResume(c *gin.Context) {
	var req struct {
		ResumeID           string `json:"resumeId" binding:"required"`
		JobDescriptionURL  string `json:"jobDescriptionUrl"`
		JobDescriptionText string `json:"jobDescriptionText"`
		AIModel           string `json:"aiModel" binding:"required"`
		KeepOnePage       bool   `json:"keepOnePage"`
		UserAPIKey        string `json:"userApiKey"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Validate that either URL or text is provided
	if req.JobDescriptionURL == "" && req.JobDescriptionText == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Either job description URL or text must be provided"})
		return
	}

	// Get the resume from database
	var resume models.Resume
	if err := database.GetDB().First(&resume, "id = ?", req.ResumeID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Resume not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
		return
	}

	// Read the resume file content
	resumeContent, err := readResumeContent(resume.OriginalContent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read resume content: " + err.Error()})
		return
	}

	// Get job description
	jobDescription := req.JobDescriptionText
	if req.JobDescriptionURL != "" {
		scraper := services.NewJobScraper()
		fetchedDesc, err := scraper.FetchJobDescription(req.JobDescriptionURL)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to fetch job description from URL: " + err.Error()})
			return
		}
		jobDescription = fetchedDesc
	}

	// For now, we'll use a placeholder API key from environment if user doesn't provide one
	apiKey := req.UserAPIKey
	if apiKey == "" {
		apiKey = os.Getenv("DEFAULT_AI_API_KEY")
		if apiKey == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "API key is required. Please provide your own API key or contact administrator."})
			return
		}
	}

	// Create optimization session in database
	sessionID := uuid.New().String()
	userID := ""
	if resume.UserID != nil {
		userID = *resume.UserID
	}
	
	session := models.OptimizationSession{
		ID:                 sessionID,
		UserID:            userID,
		ResumeID:          req.ResumeID,
		JobDescriptionURL:  &req.JobDescriptionURL,
		JobDescriptionText: &jobDescription,
		AIModel:           req.AIModel,
		KeepOnePage:       req.KeepOnePage,
		Status:            "processing",
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	if err := database.GetDB().Create(&session).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create optimization session: " + err.Error()})
		return
	}

	// Optimize the resume using AI
	optimizer := services.NewAIOptimizer()
	optimizationReq := services.OptimizationRequest{
		ResumeContent:  resumeContent,
		JobDescription: jobDescription,
		AIModel:       req.AIModel,
		KeepOnePage:   req.KeepOnePage,
		UserAPIKey:    apiKey,
	}

	result, err := optimizer.OptimizeResume(optimizationReq)
	if err != nil {
		// Update session status to failed
		database.GetDB().Model(&session).Updates(map[string]interface{}{
			"status":     "failed",
			"updated_at": time.Now(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to optimize resume: " + err.Error()})
		return
	}

	// Update session with results
	session.OptimizedContent = &result.OptimizedContent
	session.Status = "completed"
	session.UpdatedAt = time.Now()

	if err := database.GetDB().Save(&session).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save optimization results: " + err.Error()})
		return
	}

	// Return the completed session
	c.JSON(http.StatusOK, gin.H{
		"session": session,
		"summary": result.Summary,
		"changes": result.Changes,
	})
}

// readResumeContent reads the content from a resume file
func readResumeContent(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	// For now, return raw content. In production, you might want to:
	// - Parse PDF content using a library like github.com/ledongthuc/pdf
	// - Parse DOCX content using a library like github.com/nguyenthenguyen/docx
	// - Extract plain text for AI processing
	return string(content), nil
}

// ApplyFeedback applies feedback to a resume
func ApplyFeedback(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Apply feedback in development"})
}
