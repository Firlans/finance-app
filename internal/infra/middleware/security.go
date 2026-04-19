package middleware

import (
	"strconv"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/google/uuid"
)

func SecurityHeaders() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("X-Frame-Options", "DENY")
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-XSS-Protection", "1; mode=block")
		c.Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// Only set HSTS in production with HTTPS
		if c.Protocol() == "https" {
			c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}

		return c.Next()
	}
}

// RequestID adds unique request ID for tracing
func RequestID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestID := c.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set("X-Request-ID", requestID)
		c.Locals("request_id", requestID)
		return c.Next()
	}
}

// CORS returns configured CORS middleware
func CORS(allowedOrigins []string) fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     joinOrigins(allowedOrigins),
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Request-ID",
		ExposeHeaders:    "X-Request-ID",
		AllowCredentials: true,
		MaxAge:           86400, // 24 hours
	})
}

func joinOrigins(origins []string) string {
	if len(origins) == 0 {
		return "*" // Dev only
	}
	result := ""
	for i, origin := range origins {
		if i > 0 {
			result += ","
		}
		result += origin
	}
	return result
}

func buildLimitReachedResponse(c *fiber.Ctx, message string) fiber.Map {
	retryAfter := c.GetRespHeader(fiber.HeaderRetryAfter)
	seconds, err := strconv.Atoi(retryAfter)
	response := fiber.Map{
		"error": message,
	}
	if err == nil && seconds > 0 {
		response["retry_after_seconds"] = seconds
		response["retry_after_minutes"] = (seconds + 59) / 60
	}
	return response
}

// RateLimiter prevents brute force attacks
func RateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        100,
		Expiration: 15 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(buildLimitReachedResponse(c, "Too many requests. Please try again later."))
		},
	})
}

// AuthRateLimiter stricter limits for auth endpoints
func AuthRateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        5,
		Expiration: 15 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(buildLimitReachedResponse(c, "Too many authentication attempts. Account temporarily locked."))
		},
	})
}

// JWT Blacklist - In-memory implementation (use Redis in production)
type TokenBlacklist struct {
	tokens map[string]time.Time
	mu     sync.RWMutex
}

var blacklist = &TokenBlacklist{
	tokens: make(map[string]time.Time),
}

// AddToBlacklist adds token to blacklist
func (tb *TokenBlacklist) Add(token string, expiry time.Time) {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	tb.tokens[token] = expiry

	// Cleanup expired tokens
	go tb.cleanup()
}

// IsBlacklisted checks if token is blacklisted
func (tb *TokenBlacklist) IsBlacklisted(token string) bool {
	tb.mu.RLock()
	defer tb.mu.RUnlock()

	expiry, exists := tb.tokens[token]
	if !exists {
		return false
	}

	// Check if expired
	if time.Now().After(expiry) {
		return false
	}

	return true
}

// cleanup removes expired tokens
func (tb *TokenBlacklist) cleanup() {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	for token, expiry := range tb.tokens {
		if now.After(expiry) {
			delete(tb.tokens, token)
		}
	}
}

// GetBlacklist returns singleton instance
func GetBlacklist() *TokenBlacklist {
	return blacklist
}
