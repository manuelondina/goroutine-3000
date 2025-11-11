package ratelimit

import (
	"sync"
	"time"
)

// Limiter manages rate limits for multiple keys (e.g., IP addresses, API keys)
// Uses token bucket algorithm for each key
type Limiter struct {
	buckets    map[string]*TokenBucket
	mu         sync.RWMutex
	capacity   int64
	refillRate int64
	interval   time.Duration
	cleanupInterval time.Duration
	lastCleanup time.Time
}

// NewLimiter creates a new multi-key rate limiter
func NewLimiter(capacity, refillRate int64, interval time.Duration) *Limiter {
	return &Limiter{
		buckets:         make(map[string]*TokenBucket),
		capacity:        capacity,
		refillRate:      refillRate,
		interval:        interval,
		cleanupInterval: 5 * time.Minute,
		lastCleanup:     time.Now(),
	}
}

// Allow checks if a request for the given key is allowed
func (l *Limiter) Allow(key string) bool {
	bucket := l.getBucket(key)
	return bucket.Allow()
}

// AllowN checks if n requests for the given key are allowed
func (l *Limiter) AllowN(key string, n int64) bool {
	bucket := l.getBucket(key)
	return bucket.AllowN(n)
}

// getBucket returns or creates a token bucket for the given key
func (l *Limiter) getBucket(key string) *TokenBucket {
	// Fast path: read lock for existing bucket
	l.mu.RLock()
	bucket, exists := l.buckets[key]
	l.mu.RUnlock()

	if exists {
		return bucket
	}

	// Slow path: write lock to create new bucket
	l.mu.Lock()
	defer l.mu.Unlock()

	// Double-check after acquiring write lock
	bucket, exists = l.buckets[key]
	if exists {
		return bucket
	}

	// Create new bucket
	bucket = NewTokenBucket(l.capacity, l.refillRate, l.interval)
	l.buckets[key] = bucket

	// Opportunistically cleanup old buckets
	l.cleanupIfNeeded()

	return bucket
}

// cleanupIfNeeded removes buckets that are at full capacity (inactive)
// Must be called with write lock held
func (l *Limiter) cleanupIfNeeded() {
	now := time.Now()
	if now.Sub(l.lastCleanup) < l.cleanupInterval {
		return
	}

	for key, bucket := range l.buckets {
		// Remove buckets that are full (haven't been used recently)
		if bucket.Available() == bucket.Capacity() {
			delete(l.buckets, key)
		}
	}

	l.lastCleanup = now
}

// Stats returns statistics about the limiter
func (l *Limiter) Stats() map[string]interface{} {
	l.mu.RLock()
	defer l.mu.RUnlock()

	return map[string]interface{}{
		"total_keys":   len(l.buckets),
		"capacity":     l.capacity,
		"refill_rate":  l.refillRate,
		"interval_ms":  l.interval.Milliseconds(),
	}
}

// Reset clears all rate limit buckets for the given key
func (l *Limiter) Reset(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	delete(l.buckets, key)
}

// ResetAll clears all rate limit buckets
func (l *Limiter) ResetAll() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.buckets = make(map[string]*TokenBucket)
}
