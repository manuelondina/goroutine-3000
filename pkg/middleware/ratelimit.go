package middleware

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/manuelondina/goroutine-3000/pkg/ratelimit"
)

// KeyExtractor is a function that extracts a rate limit key from a request
type KeyExtractor func(*http.Request) string

// RateLimitConfig configures the rate limiting middleware
type RateLimitConfig struct {
	// Limiter is the rate limiter to use
	Limiter *ratelimit.Limiter

	// KeyExtractor extracts the key for rate limiting (defaults to IP-based)
	KeyExtractor KeyExtractor

	// OnRateLimitExceeded is called when a request is rate limited
	// Defaults to returning 429 Too Many Requests
	OnRateLimitExceeded func(http.ResponseWriter, *http.Request)

	// SkipFunc determines if rate limiting should be skipped for a request
	SkipFunc func(*http.Request) bool
}

// RateLimit returns HTTP middleware that applies rate limiting
func RateLimit(config RateLimitConfig) func(http.Handler) http.Handler {
	// Set defaults
	if config.KeyExtractor == nil {
		config.KeyExtractor = IPKeyExtractor
	}

	if config.OnRateLimitExceeded == nil {
		config.OnRateLimitExceeded = DefaultRateLimitHandler
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip rate limiting if configured
			if config.SkipFunc != nil && config.SkipFunc(r) {
				next.ServeHTTP(w, r)
				return
			}

			// Extract key and check rate limit
			key := config.KeyExtractor(r)
			if !config.Limiter.Allow(key) {
				config.OnRateLimitExceeded(w, r)
				return
			}

			// Add rate limit headers
			AddRateLimitHeaders(w, config.Limiter, key)

			next.ServeHTTP(w, r)
		})
	}
}

// IPKeyExtractor extracts the client IP address as the rate limit key
func IPKeyExtractor(r *http.Request) string {
	// Try X-Forwarded-For header first (for proxied requests)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// Take the first IP in the list
		if host, _, err := net.SplitHostPort(xff); err == nil {
			return host
		}
		return xff
	}

	// Try X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fall back to RemoteAddr
	if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		return host
	}

	return r.RemoteAddr
}

// APIKeyExtractor extracts an API key from the Authorization header
func APIKeyExtractor(r *http.Request) string {
	apiKey := r.Header.Get("Authorization")
	if apiKey == "" {
		// Fall back to IP-based limiting if no API key
		return IPKeyExtractor(r)
	}
	return apiKey
}

// PathBasedKeyExtractor combines path and IP for per-endpoint rate limiting
func PathBasedKeyExtractor(r *http.Request) string {
	ip := IPKeyExtractor(r)
	return fmt.Sprintf("%s:%s", ip, r.URL.Path)
}

// DefaultRateLimitHandler returns a 429 response when rate limit is exceeded
func DefaultRateLimitHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusTooManyRequests)
	fmt.Fprintf(w, `{"error":"rate limit exceeded","message":"Too many requests. Please try again later."}`)
}

// AddRateLimitHeaders adds standard rate limit headers to the response
func AddRateLimitHeaders(w http.ResponseWriter, limiter *ratelimit.Limiter, key string) {
	stats := limiter.Stats()
	w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", stats["capacity"]))
	w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(time.Duration(stats["interval_ms"].(int64))*time.Millisecond).Unix()))
}

// NewDefaultConfig creates a rate limit config with sensible defaults
// 100 requests per minute per IP
func NewDefaultConfig() RateLimitConfig {
	return RateLimitConfig{
		Limiter:             ratelimit.NewLimiter(100, 100, time.Minute),
		KeyExtractor:        IPKeyExtractor,
		OnRateLimitExceeded: DefaultRateLimitHandler,
	}
}

// NewAPIKeyConfig creates a rate limit config for API key-based limiting
// 1000 requests per minute per API key
func NewAPIKeyConfig() RateLimitConfig {
	return RateLimitConfig{
		Limiter:             ratelimit.NewLimiter(1000, 1000, time.Minute),
		KeyExtractor:        APIKeyExtractor,
		OnRateLimitExceeded: DefaultRateLimitHandler,
	}
}
