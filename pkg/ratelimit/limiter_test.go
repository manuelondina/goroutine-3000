package ratelimit

import (
	"sync"
	"testing"
	"time"
)

func TestNewLimiter(t *testing.T) {
	limiter := NewLimiter(10, 10, time.Minute)

	if limiter == nil {
		t.Fatal("Expected limiter to be created")
	}

	stats := limiter.Stats()
	if stats["capacity"] != int64(10) {
		t.Errorf("Expected capacity 10, got %v", stats["capacity"])
	}
}

func TestLimiterAllow(t *testing.T) {
	limiter := NewLimiter(5, 5, time.Minute)

	// Allow 5 requests for key1
	for i := 0; i < 5; i++ {
		if !limiter.Allow("key1") {
			t.Errorf("Request %d for key1 should be allowed", i+1)
		}
	}

	// 6th request should be denied
	if limiter.Allow("key1") {
		t.Error("Request 6 for key1 should be denied")
	}

	// Different key should have separate limit
	if !limiter.Allow("key2") {
		t.Error("Request 1 for key2 should be allowed")
	}
}

func TestLimiterMultipleKeys(t *testing.T) {
	limiter := NewLimiter(3, 3, time.Minute)

	keys := []string{"user1", "user2", "user3"}

	for _, key := range keys {
		// Each key should get 3 requests
		for i := 0; i < 3; i++ {
			if !limiter.Allow(key) {
				t.Errorf("Request %d for %s should be allowed", i+1, key)
			}
		}

		// 4th request should be denied
		if limiter.Allow(key) {
			t.Errorf("Request 4 for %s should be denied", key)
		}
	}

	stats := limiter.Stats()
	if stats["total_keys"] != 3 {
		t.Errorf("Expected 3 keys, got %v", stats["total_keys"])
	}
}

func TestLimiterReset(t *testing.T) {
	limiter := NewLimiter(2, 2, time.Minute)

	// Consume all tokens for key1
	limiter.Allow("key1")
	limiter.Allow("key1")

	if limiter.Allow("key1") {
		t.Error("Request should be denied before reset")
	}

	// Reset key1
	limiter.Reset("key1")

	// Should allow again
	if !limiter.Allow("key1") {
		t.Error("Request should be allowed after reset")
	}
}

func TestLimiterResetAll(t *testing.T) {
	limiter := NewLimiter(1, 1, time.Minute)

	// Consume tokens for multiple keys
	limiter.Allow("key1")
	limiter.Allow("key2")
	limiter.Allow("key3")

	stats := limiter.Stats()
	if stats["total_keys"] != 3 {
		t.Errorf("Expected 3 keys before reset, got %v", stats["total_keys"])
	}

	// Reset all
	limiter.ResetAll()

	stats = limiter.Stats()
	if stats["total_keys"] != 0 {
		t.Errorf("Expected 0 keys after reset, got %v", stats["total_keys"])
	}
}

func TestLimiterConcurrency(t *testing.T) {
	limiter := NewLimiter(100, 100, time.Minute)
	var wg sync.WaitGroup

	// Multiple goroutines accessing different keys
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			key := "user" + string(rune(id))
			for j := 0; j < 50; j++ {
				limiter.Allow(key)
			}
		}(i)
	}

	wg.Wait()

	stats := limiter.Stats()
	if stats["total_keys"] != 10 {
		t.Errorf("Expected 10 keys, got %v", stats["total_keys"])
	}
}

func TestLimiterAllowN(t *testing.T) {
	limiter := NewLimiter(10, 10, time.Minute)

	if !limiter.AllowN("key1", 5) {
		t.Error("AllowN(5) should succeed")
	}

	if !limiter.AllowN("key1", 5) {
		t.Error("AllowN(5) should succeed")
	}

	if limiter.AllowN("key1", 1) {
		t.Error("AllowN(1) should fail - no tokens left")
	}
}

func BenchmarkLimiterAllow(b *testing.B) {
	limiter := NewLimiter(int64(b.N), int64(b.N), time.Minute)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		limiter.Allow("benchmark-key")
	}
}

func BenchmarkLimiterAllowMultipleKeys(b *testing.B) {
	limiter := NewLimiter(int64(b.N), int64(b.N), time.Minute)
	keys := []string{"key1", "key2", "key3", "key4", "key5"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		limiter.Allow(keys[i%len(keys)])
	}
}

func BenchmarkLimiterConcurrent(b *testing.B) {
	limiter := NewLimiter(int64(b.N), int64(b.N), time.Minute)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			limiter.Allow("key" + string(rune(i%10)))
			i++
		}
	})
}
