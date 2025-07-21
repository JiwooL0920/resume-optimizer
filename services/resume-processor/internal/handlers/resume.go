package handlers

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/resume-optimizer/resume-processor/internal/database"
	"github.com/resume-optimizer/resume-processor/internal/models"
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
	resume := models.Resume{
		ID:              fileID,
		OriginalContent: destPath,
		FileType:        filepath.Ext(file.Filename),
		FileSize:        &fileSize,
	}

	if err := database.GetDB().Create(&resume).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": fileID})
}

// GetResume retrieves a resume by ID
func GetResume(c *gin.Context) {
	id := c.Param("id")
	var resume models.Resume

	if err := database.GetDB().First(&resume, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Resume not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"resume": resume})
}

// ListResumes lists all resumes
func ListResumes(c *gin.Context) {
	var resumes []models.Resume
	if err := database.GetDB().Find(&resumes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"resumes": resumes})
}

// DeleteResume deletes a resume by ID
func DeleteResume(c *gin.Context) {
	id := c.Param("id")
	if err := database.GetDB().Delete(&models.Resume{}, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Resume deleted", "id": id})
}

// OptimizeResume optimizes a resume
func OptimizeResume(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Resume optimization in development"})
}

// ApplyFeedback applies feedback to a resume
func ApplyFeedback(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Apply feedback in development"})
}
