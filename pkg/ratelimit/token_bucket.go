package ratelimit

import (
	"sync"
	"time"
)

// TokenBucket implements a token bucket rate limiter
// Safe for concurrent use by multiple goroutines
type TokenBucket struct {
	capacity  int64         // Maximum tokens
	tokens    int64         // Current tokens
	refillRate int64        // Tokens added per refill interval
	interval  time.Duration // Refill interval
	lastRefill time.Time
	mu        sync.Mutex
}

// NewTokenBucket creates a new token bucket rate limiter
// capacity: maximum number of tokens
// refillRate: number of tokens to add per interval
// interval: how often to refill tokens
func NewTokenBucket(capacity, refillRate int64, interval time.Duration) *TokenBucket {
	return &TokenBucket{
		capacity:   capacity,
		tokens:     capacity,
		refillRate: refillRate,
		interval:   interval,
		lastRefill: time.Now(),
	}
}

// Allow checks if a request can proceed and consumes a token if available
func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.refill()

	if tb.tokens > 0 {
		tb.tokens--
		return true
	}

	return false
}

// AllowN checks if n tokens are available and consumes them if so
func (tb *TokenBucket) AllowN(n int64) bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.refill()

	if tb.tokens >= n {
		tb.tokens -= n
		return true
	}

	return false
}

// refill adds tokens based on elapsed time since last refill
// Must be called with lock held
func (tb *TokenBucket) refill() {
	now := time.Now()
	elapsed := now.Sub(tb.lastRefill)

	if elapsed < tb.interval {
		return
	}

	// Calculate how many refill periods have passed
	periods := elapsed / tb.interval
	tokensToAdd := int64(periods) * tb.refillRate

	tb.tokens += tokensToAdd
	if tb.tokens > tb.capacity {
		tb.tokens = tb.capacity
	}

	tb.lastRefill = tb.lastRefill.Add(time.Duration(periods) * tb.interval)
}

// Available returns the current number of available tokens
func (tb *TokenBucket) Available() int64 {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.refill()
	return tb.tokens
}

// Capacity returns the maximum capacity of the bucket
func (tb *TokenBucket) Capacity() int64 {
	return tb.capacity
}
