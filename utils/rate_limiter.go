package utils

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	rate     rate.Limit
	burst    int
}

func NewRateLimiter(requestsPerSecond int, burst int) *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rate:     rate.Limit(requestsPerSecond),
		burst:    burst,
	}
}

func (rl *RateLimiter) GetLimiter(key string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.limiters[key]
	if !exists {
		limiter = rate.NewLimiter(rl.rate, rl.burst)
		rl.limiters[key] = limiter
	}

	return limiter
}

func (rl *RateLimiter) Allow(key string) bool {
	limiter := rl.GetLimiter(key)
	return limiter.Allow()
}

func RateLimitMiddleware(requestsPerSecond int, burst int) gin.HandlerFunc {
	limiter := NewRateLimiter(requestsPerSecond, burst)

	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		if !limiter.Allow(clientIP) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "Rate limit exceeded",
				"message": "Too many requests. Please try again later.",
				"code":    http.StatusTooManyRequests,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func (rl *RateLimiter) CleanupOldLimiters(maxAge time.Duration) {
	ticker := time.NewTicker(maxAge)
	go func() {
		for range ticker.C {
			rl.mu.Lock()
			for key := range rl.limiters {
				_ = key
			}
			rl.mu.Unlock()
		}
	}()
}
