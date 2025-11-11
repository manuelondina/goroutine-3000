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
	// In production, these would be your actual backend services
	err := gw.AddRoute("/api/hello", []string{
		"http://localhost:8081/api/hello",
		"http://localhost:8082/api/hello",
	})
	if err != nil {
		log.Fatalf("Failed to add route: %v", err)
	}

	err = gw.AddRoute("/api/slow", []string{
		"http://localhost:8081/api/slow",
	})
	if err != nil {
		log.Fatalf("Failed to add route: %v", err)
	}

	err = gw.AddRoute("/api/data", []string{
		"http://localhost:8081/api/data",
		"http://localhost:8082/api/data",
	})
	if err != nil {
		log.Fatalf("Failed to add route: %v", err)
	}

	// Start health checking
	gw.StartHealthCheck()
	defer gw.Stop()

	// Create mux and add gateway handler
	mux := http.NewServeMux()
	mux.Handle("/", gw.Handler())

	// Add stats endpoint
	mux.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(gw.Stats())
	})

	log.Println("Gateway starting on :8080")
	log.Println("Rate limit: 10 requests per 10 seconds per IP")
	log.Println("\nExample usage:")
	log.Println("  curl http://localhost:8080/api/hello")
	log.Println("  curl http://localhost:8080/api/data")
	log.Println("  curl http://localhost:8080/stats")
	log.Println("\nTo test rate limiting:")
	log.Println("  for i in {1..15}; do curl http://localhost:8080/api/hello; done")
	
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
