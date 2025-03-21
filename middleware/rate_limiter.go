package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// Rate limiter structure
type RateLimiter struct {
	clients map[string][]time.Time
	mu      sync.Mutex
}

// Create new rate limiter instance
var limiter = RateLimiter{
	clients: make(map[string][]time.Time),
}

// Rate limiting middleware (Max 5 messages per 10 seconds)
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.GetString("username") // Get username from JWT
		if user == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		limiter.mu.Lock()
		defer limiter.mu.Unlock()

		// Clean up old timestamps
		now := time.Now()
		timeWindow := 10 * time.Second
		allowedMessages := 5

		// Keep only recent timestamps within the window
		timestamps := limiter.clients[user]
		validTimestamps := []time.Time{}
		for _, t := range timestamps {
			if now.Sub(t) < timeWindow {
				validTimestamps = append(validTimestamps, t)
			}
		}

		// Check if user exceeded limit
		if len(validTimestamps) >= allowedMessages {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded. Max 5 messages per 10 seconds."})
			c.Abort()
			return
		}

		// Add new timestamp and proceed
		validTimestamps = append(validTimestamps, now)
		limiter.clients[user] = validTimestamps
		c.Next()
	}
}
