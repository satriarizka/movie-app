package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// Simple in-memory rate limiter per IP
type IPLimiter struct {
	ips map[string]*rate.Limiter
	mu  *sync.RWMutex
	r   rate.Limit
	b   int
}

func NewIPLimiter(r rate.Limit, b int) *IPLimiter {
	return &IPLimiter{
		ips: make(map[string]*rate.Limiter),
		mu:  &sync.RWMutex{},
		r:   r,
		b:   b,
	}
}

func (i *IPLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter, exists := i.ips[ip]
	if !exists {
		limiter = rate.NewLimiter(i.r, i.b)
		i.ips[ip] = limiter
	}

	return limiter
}

func RateLimiterMiddleware() gin.HandlerFunc {
	// Allow 5 requests per second with burst of 10
	limiter := NewIPLimiter(5, 10)
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !limiter.GetLimiter(ip).Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"status":  false,
				"message": "Too many requests",
			})
			return
		}
		c.Next()
	}
}
