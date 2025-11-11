package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/manuelondina/goroutine-3000/pkg/middleware"
	"github.com/manuelondina/goroutine-3000/pkg/ratelimit"
)

// This example shows how to integrate rate limiting middleware into your existing service

func main() {
	// Create a rate limiter: 5 requests per 10 seconds
	limiter := ratelimit.NewLimiter(5, 5, 10*time.Second)

	// Configure rate limiting middleware
	rateLimitConfig := middleware.RateLimitConfig{
		Limiter:      limiter,
		KeyExtractor: middleware.IPKeyExtractor, // Rate limit by IP
		OnRateLimitExceeded: func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error":   "rate_limit_exceeded",
				"message": "You've made too many requests. Please slow down!",
			})
		},
	}

	// Create your HTTP handlers
	mux := http.NewServeMux()

	// Regular endpoint with rate limiting
	mux.Handle("/api/public", middleware.RateLimit(rateLimitConfig)(
		http.HandlerFunc(publicHandler),
	))

	// API key-based rate limiting example
	apiKeyLimiter := ratelimit.NewLimiter(100, 100, time.Minute)
	apiKeyConfig := middleware.RateLimitConfig{
		Limiter:      apiKeyLimiter,
		KeyExtractor: middleware.APIKeyExtractor, // Rate limit by API key
	}

	mux.Handle("/api/protected", middleware.RateLimit(apiKeyConfig)(
		http.HandlerFunc(protectedHandler),
	))

	// Per-path rate limiting example
	pathLimiter := ratelimit.NewLimiter(3, 3, 10*time.Second)
	pathConfig := middleware.RateLimitConfig{
		Limiter:      pathLimiter,
		KeyExtractor: middleware.PathBasedKeyExtractor, // Rate limit by IP+Path
	}

	mux.Handle("/api/expensive", middleware.RateLimit(pathConfig)(
		http.HandlerFunc(expensiveHandler),
	))

	// Stats endpoint (no rate limiting)
	mux.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"public_limiter":    limiter.Stats(),
			"api_key_limiter":   apiKeyLimiter.Stats(),
			"path_limiter":      pathLimiter.Stats(),
		})
	})

	log.Println("Server starting on :9090")
	log.Println("\nEndpoints:")
	log.Println("  /api/public     - 5 requests per 10 seconds (by IP)")
	log.Println("  /api/protected  - 100 requests per minute (by API key)")
	log.Println("  /api/expensive  - 3 requests per 10 seconds (by IP+Path)")
	log.Println("  /stats          - View rate limit statistics")
	log.Println("\nExample usage:")
	log.Println("  curl http://localhost:9090/api/public")
	log.Println("  curl -H 'Authorization: my-api-key' http://localhost:9090/api/protected")
	log.Println("  for i in {1..10}; do curl http://localhost:9090/api/public; echo; done")

	if err := http.ListenAndServe(":9090", mux); err != nil {
		log.Fatal(err)
	}
}

func publicHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Hello from public API!",
		"time":    time.Now().Format(time.RFC3339),
	})
}

func protectedHandler(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Header.Get("Authorization")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Hello from protected API!",
		"api_key": fmt.Sprintf("%.10s...", apiKey),
		"time":    time.Now().Format(time.RFC3339),
	})
}

func expensiveHandler(w http.ResponseWriter, r *http.Request) {
	// Simulate expensive operation
	time.Sleep(500 * time.Millisecond)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "This is an expensive operation",
		"time":    time.Now().Format(time.RFC3339),
	})
}
