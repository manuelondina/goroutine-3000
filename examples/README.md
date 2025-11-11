# Gateway Examples

This directory contains working examples of the rate-limited API gateway.

## Examples

### 1. Middleware Integration (`middleware/`)

Shows how to integrate rate limiting into your existing HTTP service as middleware.

**Run:**
```bash
go run middleware/main.go
```

**Test:**
```bash
# Should succeed (under limit)
curl http://localhost:9090/api/public

# Send 10 requests to hit the limit
for i in {1..10}; do 
    curl http://localhost:9090/api/public
    echo
done
```

**Features Demonstrated:**
- IP-based rate limiting
- API key-based rate limiting
- Per-path rate limiting
- Custom rate limit handlers

### 2. Full Gateway (`gateway/`)

Complete reverse proxy gateway with load balancing and health checking.

**Run:**
```bash
# Terminal 1: Start first backend
PORT=8081 go run backend/main.go

# Terminal 2: Start second backend
PORT=8082 go run backend/main.go

# Terminal 3: Start gateway
go run gateway/main.go
```

**Test:**
```bash
# Make requests - observe round-robin load balancing
for i in {1..6}; do 
    curl http://localhost:8080/api/hello
    echo
done

# Check gateway stats
curl http://localhost:8080/stats | jq

# Test rate limiting
for i in {1..15}; do 
    curl -w "\nStatus: %{http_code}\n" http://localhost:8080/api/hello
done
```

**Features Demonstrated:**
- Reverse proxy with load balancing
- Rate limiting per client IP
- Health checking with automatic failover
- Statistics endpoint

### 3. Backend Service (`backend/`)

Sample backend service for testing the gateway.

**Endpoints:**
- `/health` - Health check endpoint
- `/api/hello` - Simple API endpoint
- `/api/slow` - Endpoint with 1s delay
- `/api/data` - Returns sample data

**Run:**
```bash
# Default port 8081
go run backend/main.go

# Custom port
PORT=8082 go run backend/main.go
```

## Quick Start

The fastest way to see everything in action:

```bash
# In goroutine-3000 root directory
cd examples

# Terminal 1
PORT=8081 go run backend/main.go

# Terminal 2
PORT=8082 go run backend/main.go

# Terminal 3
go run gateway/main.go

# Terminal 4
curl http://localhost:8080/api/hello
```

## Integration Patterns

### Pattern 1: Simple Middleware

```go
import "github.com/manuelondina/goroutine-3000/pkg/middleware"
import "github.com/manuelondina/goroutine-3000/pkg/ratelimit"

limiter := ratelimit.NewLimiter(100, 100, time.Minute)
config := middleware.RateLimitConfig{
    Limiter: limiter,
    KeyExtractor: middleware.IPKeyExtractor,
}

http.Handle("/api", middleware.RateLimit(config)(yourHandler))
```

### Pattern 2: Full Gateway

```go
import "github.com/manuelondina/goroutine-3000/pkg/gateway"

gw := gateway.NewGateway(gateway.Config{
    RateLimitCapacity: 100,
    RateLimitInterval: time.Minute,
})

gw.AddRoute("/api", []string{
    "http://backend1:8080/api",
    "http://backend2:8080/api",
})

gw.StartHealthCheck()
http.ListenAndServe(":8080", gw.Handler())
```

## Testing Rate Limiting

### Test Script

```bash
#!/bin/bash
# test-rate-limit.sh

echo "Sending 15 requests (limit: 10 per 10 seconds)"
for i in {1..15}; do
    response=$(curl -s -w "%{http_code}" http://localhost:9090/api/public)
    echo "Request $i: $response"
    sleep 0.1
done
```

### Expected Output

```
Request 1: 200
Request 2: 200
...
Request 10: 200
Request 11: 429 Too Many Requests
Request 12: 429 Too Many Requests
...
```

## Load Testing

### Using Apache Bench

```bash
# Install
brew install apache-bench  # macOS

# Test
ab -n 1000 -c 10 http://localhost:8080/api/hello
```

### Using hey

```bash
# Install
go install github.com/rakyll/hey@latest

# Test
hey -n 1000 -c 10 http://localhost:8080/api/hello
```

## Monitoring

### Check Gateway Stats

```bash
# Pretty print with jq
curl -s http://localhost:8080/stats | jq

# Watch stats in real-time
watch -n 1 'curl -s http://localhost:8080/stats | jq'
```

### Example Stats Output

```json
{
  "routes": {
    "/api/hello": {
      "total_backends": 2,
      "alive_backends": 2
    }
  },
  "rate_limit": {
    "total_keys": 3,
    "capacity": 100,
    "refill_rate": 100,
    "interval_ms": 60000
  }
}
```

## Troubleshooting

### Gateway can't connect to backends

1. Make sure backend services are running
2. Check the backend URLs in gateway config
3. Look for error messages in gateway logs

### Rate limiting not working

1. Verify the rate limit configuration
2. Check if requests are coming from the same IP
3. Look at the rate limit headers in responses

### Health checks failing

1. Ensure backends have a working endpoint
2. Check health check interval settings
3. Review backend logs for errors

## More Information

- [Quick Start Guide](../docs/GATEWAY_QUICKSTART.md)
- [Main README](../README.md)
- [Gateway Documentation](../GATEWAY.md)
