# Rate-Limited API Gateway - Quick Start Guide

This guide will help you integrate the rate-limited API gateway into your Go services.

## Table of Contents

1. [Middleware Integration (Easiest)](#middleware-integration)
2. [Gateway Mode (Full Featured)](#gateway-mode)
3. [Configuration Options](#configuration-options)
4. [Testing Your Integration](#testing)
5. [Production Best Practices](#production-best-practices)

---

## Middleware Integration

The easiest way to add rate limiting to your existing HTTP service.

### Basic Example

```go
package main

import (
    "encoding/json"
    "net/http"
    "time"
    
    "github.com/manuelondina/goroutine-3000/pkg/middleware"
    "github.com/manuelondina/goroutine-3000/pkg/ratelimit"
)

func main() {
    // Create rate limiter: 100 requests per minute per IP
    limiter := ratelimit.NewLimiter(100, 100, time.Minute)
    
    // Configure middleware
    config := middleware.RateLimitConfig{
        Limiter:      limiter,
        KeyExtractor: middleware.IPKeyExtractor, // Rate limit by IP
    }
    
    // Create your handlers
    mux := http.NewServeMux()
    
    // Apply rate limiting to specific endpoints
    mux.Handle("/api/public", middleware.RateLimit(config)(
        http.HandlerFunc(yourHandler),
    ))
    
    http.ListenAndServe(":8080", mux)
}

func yourHandler(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(map[string]string{
        "message": "Hello!",
    })
}
```

### Different Rate Limit Strategies

#### 1. IP-Based Rate Limiting

```go
config := middleware.RateLimitConfig{
    Limiter:      ratelimit.NewLimiter(100, 100, time.Minute),
    KeyExtractor: middleware.IPKeyExtractor, // By client IP
}
```

#### 2. API Key-Based Rate Limiting

```go
config := middleware.RateLimitConfig{
    Limiter:      ratelimit.NewLimiter(1000, 1000, time.Minute),
    KeyExtractor: middleware.APIKeyExtractor, // By Authorization header
}
```

#### 3. Per-Path Rate Limiting

```go
config := middleware.RateLimitConfig{
    Limiter:      ratelimit.NewLimiter(50, 50, time.Minute),
    KeyExtractor: middleware.PathBasedKeyExtractor, // By IP + Path
}
```

#### 4. Custom Key Extractor

```go
config := middleware.RateLimitConfig{
    Limiter: ratelimit.NewLimiter(100, 100, time.Minute),
    KeyExtractor: func(r *http.Request) string {
        // Custom logic - e.g., rate limit by user ID from JWT
        userID := getUserIDFromJWT(r)
        return userID
    },
}
```

### Custom Rate Limit Response

```go
config := middleware.RateLimitConfig{
    Limiter:      limiter,
    KeyExtractor: middleware.IPKeyExtractor,
    OnRateLimitExceeded: func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusTooManyRequests)
        json.NewEncoder(w).Encode(map[string]interface{}{
            "error": "slow_down",
            "message": "Too many requests. Try again later.",
            "retry_after": 60,
        })
    },
}
```

### Skip Rate Limiting for Certain Requests

```go
config := middleware.RateLimitConfig{
    Limiter:      limiter,
    KeyExtractor: middleware.IPKeyExtractor,
    SkipFunc: func(r *http.Request) bool {
        // Skip rate limiting for health checks
        return r.URL.Path == "/health"
    },
}
```

---

## Gateway Mode

Full-featured reverse proxy with rate limiting, load balancing, and health checks.

### Basic Gateway

```go
package main

import (
    "log"
    "net/http"
    "time"
    
    "github.com/manuelondina/goroutine-3000/pkg/gateway"
)

func main() {
    // Create gateway
    gw := gateway.NewGateway(gateway.Config{
        RateLimitCapacity:   100,              // 100 requests
        RateLimitRefill:     100,              // refill 100 tokens
        RateLimitInterval:   time.Minute,      // every minute
        HealthCheckInterval: 10 * time.Second, // check backends every 10s
    })
    
    // Add routes with multiple backends (automatic load balancing)
    err := gw.AddRoute("/api/users", []string{
        "http://backend1:8080/api/users",
        "http://backend2:8080/api/users",
        "http://backend3:8080/api/users",
    })
    if err != nil {
        log.Fatal(err)
    }
    
    err = gw.AddRoute("/api/products", []string{
        "http://backend1:8080/api/products",
    })
    if err != nil {
        log.Fatal(err)
    }
    
    // Start health checking
    gw.StartHealthCheck()
    defer gw.Stop()
    
    // Start serving
    log.Println("Gateway listening on :8080")
    http.ListenAndServe(":8080", gw.Handler())
}
```

### Gateway with Monitoring

```go
mux := http.NewServeMux()

// Gateway handler
mux.Handle("/", gw.Handler())

// Statistics endpoint
mux.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(gw.Stats())
})

// Health endpoint
mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("OK"))
})

http.ListenAndServe(":8080", mux)
```

---

## Configuration Options

### Rate Limiter Configuration

```go
limiter := ratelimit.NewLimiter(
    capacity,    // Maximum number of requests allowed in burst
    refillRate,  // Number of requests to refill per interval
    interval,    // Time interval for refill
)
```

**Examples:**

```go
// 100 requests per minute
ratelimit.NewLimiter(100, 100, time.Minute)

// 10 requests per second
ratelimit.NewLimiter(10, 10, time.Second)

// 1000 requests per hour
ratelimit.NewLimiter(1000, 1000, time.Hour)

// Burst of 50, refill 10 per second
ratelimit.NewLimiter(50, 10, time.Second)
```

### Gateway Configuration

```go
gateway.Config{
    RateLimitCapacity:   100,              // Burst capacity
    RateLimitRefill:     100,              // Refill rate
    RateLimitInterval:   time.Minute,      // Refill interval
    HealthCheckInterval: 10 * time.Second, // Health check frequency
}
```

---

## Testing

### Test Rate Limiting

```bash
# Send 15 requests (assuming 10 req/10sec limit)
for i in {1..15}; do
    curl -w "\nStatus: %{http_code}\n" http://localhost:8080/api/hello
    echo "---"
done
```

Expected output:
- First 10 requests: `200 OK`
- Remaining 5 requests: `429 Too Many Requests`

### Test Load Balancing

Run multiple backend servers:

```bash
# Terminal 1
PORT=8081 go run examples/backend/main.go

# Terminal 2
PORT=8082 go run examples/backend/main.go

# Terminal 3
go run examples/gateway/main.go
```

Make requests and observe responses from different ports:

```bash
for i in {1..10}; do
    curl http://localhost:8080/api/hello
done
```

### Benchmark Performance

```bash
# Install Apache Bench (if needed)
brew install apache-bench  # macOS

# Benchmark: 1000 requests, 10 concurrent
ab -n 1000 -c 10 http://localhost:8080/api/hello
```

---

## Production Best Practices

### 1. **Different Limits for Different Endpoints**

```go
// Expensive operations - lower limit
expensiveConfig := middleware.RateLimitConfig{
    Limiter: ratelimit.NewLimiter(10, 10, time.Minute),
    KeyExtractor: middleware.IPKeyExtractor,
}
mux.Handle("/api/expensive", middleware.RateLimit(expensiveConfig)(expensiveHandler))

// Regular operations - higher limit
regularConfig := middleware.RateLimitConfig{
    Limiter: ratelimit.NewLimiter(100, 100, time.Minute),
    KeyExtractor: middleware.IPKeyExtractor,
}
mux.Handle("/api/regular", middleware.RateLimit(regularConfig)(regularHandler))
```

### 2. **Tiered Rate Limits**

```go
// Free tier: 100 req/minute
freeLimiter := ratelimit.NewLimiter(100, 100, time.Minute)

// Premium tier: 1000 req/minute
premiumLimiter := ratelimit.NewLimiter(1000, 1000, time.Minute)

// Choose limiter based on user tier
config := middleware.RateLimitConfig{
    Limiter: limiter, // Choose based on user
    KeyExtractor: func(r *http.Request) string {
        user := getUserFromAuth(r)
        if user.IsPremium {
            return "premium-" + user.ID
        }
        return "free-" + user.ID
    },
}
```

### 3. **Graceful Shutdown**

```go
gw := gateway.NewGateway(config)
gw.StartHealthCheck()

// Catch interrupt signal
sigChan := make(chan os.Signal, 1)
signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

go func() {
    <-sigChan
    log.Println("Shutting down gateway...")
    gw.Stop()
    os.Exit(0)
}()

http.ListenAndServe(":8080", gw.Handler())
```

### 4. **Monitoring and Alerting**

```go
// Expose metrics endpoint
mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
    stats := gw.Stats()
    
    // Check for unhealthy backends
    for path, routeStats := range stats["routes"].(map[string]interface{}) {
        rs := routeStats.(map[string]interface{})
        if rs["alive_backends"].(int) == 0 {
            log.Printf("ALERT: No healthy backends for %s", path)
        }
    }
    
    json.NewEncoder(w).Encode(stats)
})
```

### 5. **Environment-Based Configuration**

```go
func getConfig() gateway.Config {
    capacity := getEnvInt("RATE_LIMIT_CAPACITY", 100)
    interval := getEnvDuration("RATE_LIMIT_INTERVAL", time.Minute)
    
    return gateway.Config{
        RateLimitCapacity:   int64(capacity),
        RateLimitRefill:     int64(capacity),
        RateLimitInterval:   interval,
        HealthCheckInterval: 10 * time.Second,
    }
}
```

---

## Complete Example

See the `examples/` directory for complete, runnable examples:

- **`examples/middleware/`** - Middleware integration example
- **`examples/gateway/`** - Full gateway example
- **`examples/backend/`** - Sample backend service

```bash
# Run the middleware example
go run examples/middleware/main.go

# Run the full gateway (needs backends)
PORT=8081 go run examples/backend/main.go &
PORT=8082 go run examples/backend/main.go &
go run examples/gateway/main.go
```

---

## Need Help?

- Check the main [README.md](../README.md)
- Review the [examples/](../examples/) directory
- Open an issue on GitHub

Happy rate-limiting! 🚀
