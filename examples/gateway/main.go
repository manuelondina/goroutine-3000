package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/manuelondina/goroutine-3000/pkg/gateway"
)

func main() {
	// Create gateway with rate limiting: 10 requests per 10 seconds
	gw := gateway.NewGateway(gateway.Config{
		RateLimitCapacity:   10,
		RateLimitRefill:     10,
		RateLimitInterval:   10 * time.Second,
		HealthCheckInterval: 5 * time.Second,
	})

	// Add routes with multiple backends for load balancing
	// Backend URLs should be just the base URL, the path will be forwarded
	err := gw.AddRoute("/api/hello", []string{
		"http://localhost:8081",
		"http://localhost:8082",
	})
	if err != nil {
		log.Fatalf("Failed to add route: %v", err)
	}

	err = gw.AddRoute("/api/slow", []string{
		"http://localhost:8081",
	})
	if err != nil {
		log.Fatalf("Failed to add route: %v", err)
	}

	err = gw.AddRoute("/api/data", []string{
		"http://localhost:8081",
		"http://localhost:8082",
	})
	if err != nil {
		log.Fatalf("Failed to add route: %v", err)
	}
	// Start health checking
	gw.StartHealthCheck()
	defer gw.Stop()

	// Custom handler that routes to stats or gateway
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/stats" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(gw.Stats())
			return
		}
		gw.Handler().ServeHTTP(w, r)
	})

	log.Println("Gateway starting on :8080")
	log.Println("Rate limit: 10 requests per 10 seconds per IP")
	log.Println("\nExample usage:")
	log.Println("  curl http://localhost:8080/api/hello")
	log.Println("  curl http://localhost:8080/api/data")
	log.Println("  curl http://localhost:8080/stats")
	log.Println("\nTo test rate limiting:")
	log.Println("  for i in {1..15}; do curl http://localhost:8080/api/hello; done")
	
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
