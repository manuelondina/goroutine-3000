package ratelimit

import (
	"sync"
	"testing"
	"time"
)

func TestNewTokenBucket(t *testing.T) {
	tb := NewTokenBucket(10, 5, time.Second)

	if tb.Capacity() != 10 {
		t.Errorf("Expected capacity 10, got %d", tb.Capacity())
	}

	if tb.Available() != 10 {
		t.Errorf("Expected available tokens 10, got %d", tb.Available())
	}
}

func TestTokenBucketAllow(t *testing.T) {
	tb := NewTokenBucket(5, 5, time.Second)

	// Should allow 5 requests
	for i := 0; i < 5; i++ {
		if !tb.Allow() {
			t.Errorf("Request %d should be allowed", i+1)
		}
	}

	// 6th request should be denied
	if tb.Allow() {
		t.Error("Request 6 should be denied")
	}

	// Available should be 0
	if tb.Available() != 0 {
		t.Errorf("Expected 0 available tokens, got %d", tb.Available())
	}
}

func TestTokenBucketAllowN(t *testing.T) {
	tb := NewTokenBucket(10, 10, time.Second)

	// Should allow consuming 5 tokens
	if !tb.AllowN(5) {
		t.Error("AllowN(5) should succeed")
	}

	// Should have 5 tokens left
	if tb.Available() != 5 {
		t.Errorf("Expected 5 available tokens, got %d", tb.Available())
	}

	// Should not allow consuming 6 tokens
	if tb.AllowN(6) {
		t.Error("AllowN(6) should fail")
	}

	// Should allow consuming 5 tokens
	if !tb.AllowN(5) {
		t.Error("AllowN(5) should succeed")
	}
}

func TestTokenBucketRefill(t *testing.T) {
	tb := NewTokenBucket(10, 10, 100*time.Millisecond)

	// Consume all tokens
	for i := 0; i < 10; i++ {
		tb.Allow()
	}

	if tb.Available() != 0 {
		t.Errorf("Expected 0 available tokens, got %d", tb.Available())
	}

	// Wait for refill
	time.Sleep(150 * time.Millisecond)

	// Should have refilled
	available := tb.Available()
	if available != 10 {
		t.Errorf("Expected 10 available tokens after refill, got %d", available)
	}
}

func TestTokenBucketConcurrency(t *testing.T) {
	tb := NewTokenBucket(100, 100, time.Second)
	var wg sync.WaitGroup
	allowed := 0
	denied := 0
	var mu sync.Mutex

	// Try 200 concurrent requests
	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if tb.Allow() {
				mu.Lock()
				allowed++
				mu.Unlock()
			} else {
				mu.Lock()
				denied++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	if allowed != 100 {
		t.Errorf("Expected 100 allowed requests, got %d", allowed)
	}

	if denied != 100 {
		t.Errorf("Expected 100 denied requests, got %d", denied)
	}
}

func BenchmarkTokenBucketAllow(b *testing.B) {
	tb := NewTokenBucket(int64(b.N), int64(b.N), time.Minute)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tb.Allow()
	}
}

func BenchmarkTokenBucketAllowConcurrent(b *testing.B) {
	tb := NewTokenBucket(int64(b.N), int64(b.N), time.Minute)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			tb.Allow()
		}
	})
}
