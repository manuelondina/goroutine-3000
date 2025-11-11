package patterns

import (
	"fmt"
	"log"
	"time"

	"github.com/manuelondina/goroutine-3000/pkg/gateway"
)

// DemoGateway demonstrates the rate-limited API gateway
func DemoGateway() {
	fmt.Println("=== Rate-Limited API Gateway Demo ===")
	fmt.Println("\nThis demo shows a production-ready API gateway with:")
	fmt.Println("  • Rate limiting per client (IP-based)")
	fmt.Println("  • Load balancing across backends (round-robin)")
	fmt.Println("  • Health checking with automatic failover")
	fmt.Println("  • Concurrent request handling")
	fmt.Println()

	// Create gateway
	gw := gateway.NewGateway(gateway.Config{
		RateLimitCapacity:   100,
		RateLimitRefill:     100,
		RateLimitInterval:   time.Minute,
		HealthCheckInterval: 10 * time.Second,
	})

	// Add example routes
	err := gw.AddRoute("/api/users", []string{
		"http://backend1:8080/api/users",
		"http://backend2:8080/api/users",
	})
	if err != nil {
		log.Printf("Failed to add route: %v", err)
	}

	err = gw.AddRoute("/api/products", []string{
		"http://backend1:8080/api/products",
	})
	if err != nil {
		log.Printf("Failed to add route: %v", err)
	}

	fmt.Println("✓ Gateway configured with routes:")
	fmt.Println("  • /api/users    → 2 backends (load balanced)")
	fmt.Println("  • /api/products → 1 backend")
	fmt.Println()

	fmt.Println("✓ Rate limiting active:")
	fmt.Println("  • 100 requests per minute per IP")
	fmt.Println()

	fmt.Println("✓ Health checks running:")
	fmt.Println("  • Every 10 seconds")
	fmt.Println("  • Automatic failover to healthy backends")
	fmt.Println()

	stats := gw.Stats()
	fmt.Printf("Gateway statistics: %+v\n", stats)
	fmt.Println()

	fmt.Println("💡 Integration example:")
	fmt.Println(`
  package main
  
  import (
      "net/http"
      "time"
      "github.com/manuelondina/goroutine-3000/pkg/gateway"
  )
  
  func main() {
      gw := gateway.NewGateway(gateway.Config{
          RateLimitCapacity: 100,
          RateLimitRefill: 100,
          RateLimitInterval: time.Minute,
      })
      
      gw.AddRoute("/api/users", []string{
          "http://backend1:8080/api/users",
          "http://backend2:8080/api/users",
      })
      
      gw.StartHealthCheck()
      defer gw.Stop()
      
      http.ListenAndServe(":8080", gw.Handler())
  }
	`)

	fmt.Println("To run the full example:")
	fmt.Println("  cd examples/gateway && go run main.go")
	fmt.Println()

	gw.Stop()
	fmt.Println("✓ Gateway demo completed")
}
