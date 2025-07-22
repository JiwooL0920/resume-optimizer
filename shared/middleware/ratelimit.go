package middleware

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/resume-optimizer/shared/errors"
	"github.com/rs/zerolog/log"
)

// RateLimiter implements a simple in-memory rate limiter
type RateLimiter struct {
	visitors map[string]*Visitor
	mutex    sync.RWMutex
	rate     int           // requests per second
	burst    int           // maximum burst size
	cleanup  time.Duration // cleanup interval
}

// Visitor represents a client's visit information
type Visitor struct {
	limiter  *TokenBucket
	lastSeen time.Time
}

// TokenBucket implements a token bucket rate limiter
type TokenBucket struct {
	tokens    int
	capacity  int
	refillRate int
	lastRefill time.Time
	mutex     sync.Mutex
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(rate, burst int) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*Visitor),
		rate:     rate,
		burst:    burst,
		cleanup:  5 * time.Minute,
	}

	// Start cleanup goroutine
	go rl.cleanupVisitors()

	return rl
}

// NewTokenBucket creates a new token bucket
func NewTokenBucket(capacity, refillRate int) *TokenBucket {
	return &TokenBucket{
		tokens:     capacity,
		capacity:   capacity,
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

// Allow checks if a request should be allowed
func (tb *TokenBucket) Allow() bool {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastRefill)
	
	// Add tokens based on time elapsed
	tokensToAdd := int(elapsed.Seconds()) * tb.refillRate
	if tokensToAdd > 0 {
		tb.tokens = min(tb.capacity, tb.tokens+tokensToAdd)
		tb.lastRefill = now
	}

	// Check if we have tokens available
	if tb.tokens > 0 {
		tb.tokens--
		return true
	}

	return false
}

// getVisitor gets or creates a visitor for the given IP
func (rl *RateLimiter) getVisitor(ip string) *Visitor {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	visitor, exists := rl.visitors[ip]
	if !exists {
		visitor = &Visitor{
			limiter:  NewTokenBucket(rl.burst, rl.rate),
			lastSeen: time.Now(),
		}
		rl.visitors[ip] = visitor
	} else {
		visitor.lastSeen = time.Now()
	}

	return visitor
}

// cleanupVisitors removes old visitors
func (rl *RateLimiter) cleanupVisitors() {
	ticker := time.NewTicker(rl.cleanup)
	defer ticker.Stop()

	for range ticker.C {
		rl.mutex.Lock()
		for ip, visitor := range rl.visitors {
			if time.Since(visitor.lastSeen) > rl.cleanup {
				delete(rl.visitors, ip)
			}
		}
		rl.mutex.Unlock()
	}
}

// RateLimitMiddleware creates rate limiting middleware
func RateLimitMiddleware(rate, burst int) gin.HandlerFunc {
	if rate <= 0 || burst <= 0 {
		// If rate limiting is disabled, return a no-op middleware
		return func(c *gin.Context) {
			c.Next()
		}
	}

	limiter := NewRateLimiter(rate, burst)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		visitor := limiter.getVisitor(ip)

		if !visitor.limiter.Allow() {
			log.Warn().
				Str("ip", ip).
				Str("path", c.Request.URL.Path).
				Str("method", c.Request.Method).
				Msg("Rate limit exceeded")

			appErr := errors.NewAppError(
				errors.ErrCodeRateLimit,
				"Rate limit exceeded. Please try again later.",
				nil,
			)
			
			c.Header("Retry-After", "60") // Suggest retry after 60 seconds
			c.JSON(appErr.HTTPStatus, appErr)
			c.Abort()
			return
		}

		c.Next()
	}
}

// PerUserRateLimitMiddleware creates per-user rate limiting middleware
func PerUserRateLimitMiddleware(rate, burst int) gin.HandlerFunc {
	if rate <= 0 || burst <= 0 {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	limiter := NewRateLimiter(rate, burst)

	return func(c *gin.Context) {
		// Try to get user ID, fallback to IP
		userID, exists := c.Get("userID")
		key := c.ClientIP() // Default to IP
		
		if exists && userID != nil {
			key = "user:" + userID.(string)
		}

		visitor := limiter.getVisitor(key)

		if !visitor.limiter.Allow() {
			log.Warn().
				Str("key", key).
				Str("path", c.Request.URL.Path).
				Str("method", c.Request.Method).
				Msg("User rate limit exceeded")

			appErr := errors.NewAppError(
				errors.ErrCodeRateLimit,
				"Rate limit exceeded. Please try again later.",
				nil,
			)
			
			c.Header("Retry-After", "60")
			c.JSON(appErr.HTTPStatus, appErr)
			c.Abort()
			return
		}

		c.Next()
	}
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}