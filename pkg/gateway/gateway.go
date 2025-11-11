package gateway

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"

	"github.com/manuelondina/goroutine-3000/pkg/middleware"
	"github.com/manuelondina/goroutine-3000/pkg/ratelimit"
)

// Backend represents a backend service
type Backend struct {
	URL          *url.URL
	Alive        bool
	mu           sync.RWMutex
	ReverseProxy *httputil.ReverseProxy
}

// SetAlive sets the alive status of the backend
func (b *Backend) SetAlive(alive bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.Alive = alive
}

// IsAlive returns the alive status of the backend
func (b *Backend) IsAlive() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.Alive
}

// Route represents a route configuration
type Route struct {
	Path     string
	Backends []*Backend
	current  int
	mu       sync.Mutex
}

// NextBackend returns the next available backend using round-robin
func (r *Route) NextBackend() *Backend {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Find next alive backend
	for i := 0; i < len(r.Backends); i++ {
		r.current = (r.current + 1) % len(r.Backends)
		backend := r.Backends[r.current]
		if backend.IsAlive() {
			return backend
		}
	}

	return nil
}

// Gateway is the main API gateway
type Gateway struct {
	routes       map[string]*Route
	limiter      *ratelimit.Limiter
	mu           sync.RWMutex
	healthCheck  time.Duration
	ctx          context.Context
	cancel       context.CancelFunc
}

// Config configures the gateway
type Config struct {
	// RateLimit settings (requests per interval)
	RateLimitCapacity  int64
	RateLimitRefill    int64
	RateLimitInterval  time.Duration

	// HealthCheck interval
	HealthCheckInterval time.Duration
}

// NewGateway creates a new API gateway
func NewGateway(config Config) *Gateway {
	ctx, cancel := context.WithCancel(context.Background())

	if config.HealthCheckInterval == 0 {
		config.HealthCheckInterval = 10 * time.Second
	}

	if config.RateLimitCapacity == 0 {
		config.RateLimitCapacity = 100
	}

	if config.RateLimitRefill == 0 {
		config.RateLimitRefill = 100
	}

	if config.RateLimitInterval == 0 {
		config.RateLimitInterval = time.Minute
	}

	return &Gateway{
		routes:      make(map[string]*Route),
		limiter:     ratelimit.NewLimiter(config.RateLimitCapacity, config.RateLimitRefill, config.RateLimitInterval),
		healthCheck: config.HealthCheckInterval,
		ctx:         ctx,
		cancel:      cancel,
	}
}

// AddRoute adds a new route to the gateway
func (g *Gateway) AddRoute(path string, backendURLs []string) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	route := &Route{
		Path:     path,
		Backends: make([]*Backend, 0, len(backendURLs)),
	}

	for _, backendURL := range backendURLs {
		u, err := url.Parse(backendURL)
		if err != nil {
			return fmt.Errorf("invalid backend URL %s: %w", backendURL, err)
		}

		proxy := httputil.NewSingleHostReverseProxy(u)
		
		// Customize error handler
		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			log.Printf("Gateway error: %v", err)
			w.WriteHeader(http.StatusBadGateway)
			fmt.Fprintf(w, `{"error":"bad gateway","message":"Backend service unavailable"}`)
		}

		backend := &Backend{
			URL:          u,
			Alive:        true,
			ReverseProxy: proxy,
		}

		route.Backends = append(route.Backends, backend)
	}

	g.routes[path] = route
	return nil
}

// Handler returns the HTTP handler for the gateway
func (g *Gateway) Handler() http.Handler {
	rateLimitMiddleware := middleware.RateLimit(middleware.RateLimitConfig{
		Limiter:             g.limiter,
		KeyExtractor:        middleware.IPKeyExtractor,
		OnRateLimitExceeded: middleware.DefaultRateLimitHandler,
	})

	return rateLimitMiddleware(http.HandlerFunc(g.handleRequest))
}

// handleRequest handles incoming requests
func (g *Gateway) handleRequest(w http.ResponseWriter, r *http.Request) {
	g.mu.RLock()
	route, exists := g.routes[r.URL.Path]
	g.mu.RUnlock()

	if !exists {
		http.NotFound(w, r)
		return
	}

	backend := route.NextBackend()
	if backend == nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, `{"error":"service unavailable","message":"No healthy backends available"}`)
		return
	}

	// Proxy the request
	backend.ReverseProxy.ServeHTTP(w, r)
}

// StartHealthCheck starts health checking for all backends
func (g *Gateway) StartHealthCheck() {
	ticker := time.NewTicker(g.healthCheck)
	
	go func() {
		for {
			select {
			case <-ticker.C:
				g.healthCheckAll()
			case <-g.ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()
}

// healthCheckAll checks health of all backends
func (g *Gateway) healthCheckAll() {
	g.mu.RLock()
	routes := make([]*Route, 0, len(g.routes))
	for _, route := range g.routes {
		routes = append(routes, route)
	}
	g.mu.RUnlock()

	var wg sync.WaitGroup
	for _, route := range routes {
		for _, backend := range route.Backends {
			wg.Add(1)
			go func(b *Backend) {
				defer wg.Done()
				alive := g.isBackendAlive(b.URL)
				b.SetAlive(alive)
				if !alive {
					log.Printf("Backend %s is down", b.URL)
				}
			}(backend)
		}
	}
	wg.Wait()
}

// isBackendAlive checks if a backend is alive
func (g *Gateway) isBackendAlive(u *url.URL) bool {
	timeout := 2 * time.Second
	ctx, cancel := context.WithTimeout(g.ctx, timeout)
	defer cancel()

	// Try /health endpoint first, fall back to root
	healthURL := *u
	healthURL.Path = "/health"
	
	req, err := http.NewRequestWithContext(ctx, "GET", healthURL.String(), nil)
	if err != nil {
		return false
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode >= 200 && resp.StatusCode < 500
}

// Stop stops the gateway
func (g *Gateway) Stop() {
	g.cancel()
}

// Stats returns gateway statistics
func (g *Gateway) Stats() map[string]interface{} {
	g.mu.RLock()
	defer g.mu.RUnlock()

	routeStats := make(map[string]interface{})
	for path, route := range g.routes {
		aliveCount := 0
		for _, backend := range route.Backends {
			if backend.IsAlive() {
				aliveCount++
			}
		}
		routeStats[path] = map[string]interface{}{
			"total_backends": len(route.Backends),
			"alive_backends": aliveCount,
		}
	}

	return map[string]interface{}{
		"routes":     routeStats,
		"rate_limit": g.limiter.Stats(),
	}
}
