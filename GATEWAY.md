# Rate-Limited API Gateway

A production-ready, goroutine-powered API gateway with rate limiting, load balancing, and health checking.

## 🚀 Features

### Core Capabilities

✅ **Rate Limiting**
- Token bucket algorithm for smooth rate limiting
- Per-client limits (IP, API key, custom)
- Configurable burst capacity and refill rates
- Automatic cleanup of inactive limiters

✅ **Load Balancing**
- Round-robin distribution across backends
- Automatic failover to healthy backends
- Support for multiple backend services per route

✅ **Health Checking**
- Concurrent health checks using goroutines
- Configurable check intervals
- Automatic backend marking (alive/dead)

✅ **Easy Integration**
- Drop-in HTTP middleware
- Works with standard `net/http`
- Minimal configuration required

✅ **Production Ready**
- Thread-safe (tested with race detector)
- 90.8% test coverage
- Comprehensive error handling
- Statistics and monitoring endpoints

## 📦 What's Included

### Packages

```
pkg/
├── ratelimit/          # Core rate limiting logic
│   ├── token_bucket.go    # Token bucket implementation
│   ├── limiter.go         # Multi-key limiter
│   └── *_test.go          # Comprehensive tests
├── middleware/         # HTTP middleware integration
│   └── ratelimit.go       # Rate limit middleware
└── gateway/            # Full gateway implementation
    └── gateway.go         # Reverse proxy with rate limiting
```

### Examples

```
examples/
├── middleware/         # Middleware integration example
├── gateway/           # Full gateway example
└── backend/           # Sample backend service
```

### Documentation

```
docs/
└── GATEWAY_QUICKSTART.md  # Complete integration guide
```

## 🎯 Use Cases

### 1. **As Middleware** (Recommended for most cases)

Add rate limiting to your existing HTTP service:

```go
limiter := ratelimit.NewLimiter(100, 100, time.Minute)
config := middleware.RateLimitConfig{
    Limiter: limiter,
    KeyExtractor: middleware.IPKeyExtractor,
}

mux.Handle("/api", middleware.RateLimit(config)(yourHandler))
```

**Best for:**
- Adding rate limiting to existing services
- Per-endpoint rate limiting
- Custom rate limit strategies

### 2. **As Gateway** (Full-featured proxy)

Deploy as a standalone gateway in front of your services:

```go
gw := gateway.NewGateway(gateway.Config{
    RateLimitCapacity: 100,
    RateLimitInterval: time.Minute,
})

gw.AddRoute("/api/users", []string{
    "http://backend1:8080/api/users",
    "http://backend2:8080/api/users",
})

http.ListenAndServe(":8080", gw.Handler())
```

**Best for:**
- Microservices architecture
- Load balancing across multiple backends
- Centralized rate limiting
- Health checking and failover

## 🧪 Testing

All tests pass with race detection:

```bash
$ go test ./pkg/ratelimit/... -race -cover
ok      github.com/manuelondina/goroutine-3000/pkg/ratelimit    1.423s  coverage: 90.8%
```

Run the examples:

```bash
# Demo the gateway
go run cmd/gateway-demo/main.go

# Run middleware example
go run examples/middleware/main.go

# Run full gateway (with backends)
PORT=8081 go run examples/backend/main.go &
PORT=8082 go run examples/backend/main.go &
go run examples/gateway/main.go
```

## 🎨 Architecture Highlights

### Goroutine Usage

The gateway leverages goroutines for maximum concurrency:

1. **Health Checks**: Each backend is checked concurrently
2. **Request Handling**: Each request processed in its own goroutine (via HTTP server)
3. **Rate Limiting**: Lock-free reads with optimized locking for writes

### Design Patterns

- **Token Bucket**: Smooth rate limiting with burst support
- **Round-Robin**: Fair load distribution
- **Worker Pool**: Health checking uses WaitGroup for coordination
- **Middleware Pattern**: Composable HTTP handlers

### Concurrency Safety

- All data structures are thread-safe
- Tested with `-race` detector
- Proper use of mutexes (RWMutex for read-heavy operations)
- No data races in 12 concurrent tests

## 📊 Performance

### Rate Limiter Benchmarks

```
BenchmarkTokenBucketAllow-8              5000000    250 ns/op
BenchmarkTokenBucketAllowConcurrent-8    3000000    400 ns/op
BenchmarkLimiterAllow-8                  4000000    280 ns/op
BenchmarkLimiterConcurrent-8             2500000    450 ns/op
```

### Characteristics

- **Lightweight**: Token bucket operations in ~250ns
- **Scalable**: Handles thousands of concurrent requests
- **Memory Efficient**: Automatic cleanup of inactive limiters
- **Low Latency**: Minimal overhead on request processing

## 📖 Documentation

- **[GATEWAY_QUICKSTART.md](docs/GATEWAY_QUICKSTART.md)** - Complete integration guide
- **[README.md](README.md)** - Main project documentation
- **Examples** - See `examples/` directory for working code

## 🛠️ Quick Examples

### Different Rate Limit Strategies

```go
// By IP address
middleware.IPKeyExtractor

// By API key
middleware.APIKeyExtractor

// By IP + path (per-endpoint per-user)
middleware.PathBasedKeyExtractor

// Custom extractor
func(r *http.Request) string {
    return getUserID(r)
}
```

### Different Rate Limits

```go
// 100 requests per minute
ratelimit.NewLimiter(100, 100, time.Minute)

// 10 requests per second
ratelimit.NewLimiter(10, 10, time.Second)

// Burst of 50, refill 10/sec
ratelimit.NewLimiter(50, 10, time.Second)
```

## 🎓 Learning

This gateway demonstrates several goroutine patterns from the main project:

- **Worker Pool**: Health checking uses coordinated goroutines
- **Context**: Proper timeout and cancellation handling
- **Fan-Out**: Concurrent health checks across all backends
- **Error Handling**: Graceful error handling in concurrent operations

## 🚦 Integration Examples

### With Existing Router

```go
// Works with any http.Handler
handler := yourExistingRouter()
rateLimited := middleware.RateLimit(config)(handler)
http.ListenAndServe(":8080", rateLimited)
```

### Multiple Rate Limiters

```go
// Different limits for different endpoints
publicLimiter := ratelimit.NewLimiter(10, 10, time.Minute)
apiLimiter := ratelimit.NewLimiter(100, 100, time.Minute)

mux.Handle("/public", middleware.RateLimit(
    middleware.RateLimitConfig{Limiter: publicLimiter})(publicHandler))
mux.Handle("/api", middleware.RateLimit(
    middleware.RateLimitConfig{Limiter: apiLimiter})(apiHandler))
```

## 💡 Why This Implementation?

1. **Simple Integration**: Just a few lines to add rate limiting
2. **Flexible**: Multiple configuration options
3. **Production Ready**: Tested, benchmarked, race-free
4. **Go Idiomatic**: Uses standard library patterns
5. **Educational**: Shows real-world goroutine usage

## 🎯 Next Steps

1. **Read the quickstart**: [docs/GATEWAY_QUICKSTART.md](docs/GATEWAY_QUICKSTART.md)
2. **Run the examples**: `go run examples/middleware/main.go`
3. **Integrate into your service**: Copy examples and adapt
4. **Customize**: Extend with your own key extractors

---

Built with ❤️ using goroutines and the power of Go's concurrency primitives.
