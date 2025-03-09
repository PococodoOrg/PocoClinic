package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/dksch/pococlinic/internal/pkg/errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// SecurityHeaders adds security-related headers to the response
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Header("Content-Security-Policy", "default-src 'self'")
		c.Next()
	}
}

// IPRateLimiter stores rate limiters for IP addresses
type IPRateLimiter struct {
	ips    map[string]*rate.Limiter
	mu     *sync.RWMutex
	rate   rate.Limit
	burst  int
	TTL    time.Duration
	lastIP map[string]time.Time
}

// NewIPRateLimiter creates a new rate limiter
func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	return &IPRateLimiter{
		ips:    make(map[string]*rate.Limiter),
		mu:     &sync.RWMutex{},
		rate:   r,
		burst:  b,
		TTL:    time.Hour,
		lastIP: make(map[string]time.Time),
	}
}

// GetLimiter returns the rate limiter for the provided IP address
func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter, exists := i.ips[ip]
	if !exists {
		limiter = rate.NewLimiter(i.rate, i.burst)
		i.ips[ip] = limiter
	}
	i.lastIP[ip] = time.Now()
	return limiter
}

// RateLimiterMiddleware creates a new rate limiter middleware
func RateLimiterMiddleware(limiter *IPRateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		// Try to take a single token
		if !limiter.GetLimiter(ip).AllowN(time.Now(), 1) {
			c.JSON(http.StatusTooManyRequests, errors.NewAPIError(errors.ErrRateLimit, "Rate limit exceeded"))
			c.Abort()
			return
		}
		c.Next()
	}
}

// Recovery returns a middleware that recovers from panics
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.JSON(http.StatusInternalServerError, errors.NewAPIError(errors.ErrInternalServer, "An internal error occurred"))
				c.Abort()
			}
		}()
		c.Next()
	}
}

// CleanupTask starts a goroutine that periodically removes old IP entries
func (i *IPRateLimiter) CleanupTask() {
	ticker := time.NewTicker(10 * time.Minute)
	go func() {
		for range ticker.C {
			i.mu.Lock()
			for ip, lastSeen := range i.lastIP {
				if time.Since(lastSeen) > i.TTL {
					delete(i.ips, ip)
					delete(i.lastIP, ip)
				}
			}
			i.mu.Unlock()
		}
	}()
}
