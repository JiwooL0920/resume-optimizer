package gorm

import (
	"strings"

	"gorm.io/gorm"
)

// isUniqueConstraintError checks if the error is a unique constraint violation
func isUniqueConstraintError(err error) bool {
	if err == nil {
		return false
	}
	
	errStr := strings.ToLower(err.Error())
	return strings.Contains(errStr, "unique") || 
		   strings.Contains(errStr, "duplicate") ||
		   strings.Contains(errStr, "already exists")
}

// isForeignKeyConstraintError checks if the error is a foreign key constraint violation
func isForeignKeyConstraintError(err error) bool {
	if err == nil {
		return false
	}
	
	errStr := strings.ToLower(err.Error())
	return strings.Contains(errStr, "foreign key") ||
		   strings.Contains(errStr, "violates foreign key constraint")
}

// preloadWithContext applies common preloading options
func preloadWithContext(db *gorm.DB, preloads ...string) *gorm.DB {
	for _, preload := range preloads {
		db = db.Preload(preload)
	}
	return db
}

// applySorting applies sorting to the query
func applySorting(db *gorm.DB, sortBy, sortOrder string) *gorm.DB {
	if sortBy == "" {
		return db.Order("created_at DESC")
	}
	
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}
	
	return db.Order(sortBy + " " + strings.ToUpper(sortOrder))
}

// applyPagination applies pagination to the query
func applyPagination(db *gorm.DB, offset, limit int) *gorm.DB {
	if limit <= 0 {
		limit = 50 // Default limit
	}
	if limit > 100 {
		limit = 100 // Max limit
	}
	if offset < 0 {
		offset = 0
	}
	
	return db.Offset(offset).Limit(limit)
}