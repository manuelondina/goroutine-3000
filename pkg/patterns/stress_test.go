package patterns

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestStressTest(t *testing.T) {
	t.Run("spawns thousands of goroutines", func(t *testing.T) {
		const numGoroutines = 1000
		var counter atomic.Int64
		var wg sync.WaitGroup

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				counter.Add(1)
			}()
		}

		wg.Wait()

		if counter.Load() != numGoroutines {
			t.Errorf("Expected counter to be %d, got %d", numGoroutines, counter.Load())
		}
	})

	t.Run("goroutines execute concurrently", func(t *testing.T) {
		const numGoroutines = 100
		const sleepDuration = 10 * time.Millisecond
		var wg sync.WaitGroup

		start := time.Now()

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				time.Sleep(sleepDuration)
			}()
		}

		wg.Wait()
		duration := time.Since(start)

		// If sequential, would take numGoroutines * sleepDuration
		// With concurrency, should be much less
		maxExpected := sleepDuration * 20 // Allow for some overhead
		if duration > maxExpected {
			t.Errorf("Execution took too long: %v (expected < %v)", duration, maxExpected)
		}
	})

	t.Run("atomic operations are thread-safe", func(t *testing.T) {
		const numGoroutines = 1000
		const increments = 100
		var counter atomic.Int64
		var wg sync.WaitGroup

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < increments; j++ {
					counter.Add(1)
				}
			}()
		}

		wg.Wait()
		expected := int64(numGoroutines * increments)

		if counter.Load() != expected {
			t.Errorf("Expected counter to be %d, got %d", expected, counter.Load())
		}
	})
}

func BenchmarkGoroutineCreation(b *testing.B) {
	var wg sync.WaitGroup
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
		}()
	}

	wg.Wait()
}

func BenchmarkGoroutineWithWork(b *testing.B) {
	var wg sync.WaitGroup
	var counter atomic.Int64
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				counter.Add(1)
			}
		}()
	}

	wg.Wait()
}

func BenchmarkChannelCommunication(b *testing.B) {
	ch := make(chan int, 100)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < b.N; i++ {
			ch <- i
		}
		close(ch)
	}()

	b.ResetTimer()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for range ch {
		}
	}()

	wg.Wait()
}

func TestGoroutineCleanup(t *testing.T) {
	t.Run("goroutines are cleaned up after completion", func(t *testing.T) {
		before := runtime.NumGoroutine()

		var wg sync.WaitGroup
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				time.Sleep(10 * time.Millisecond)
			}()
		}

		wg.Wait()
		time.Sleep(100 * time.Millisecond) // Allow time for cleanup

		after := runtime.NumGoroutine()

		// Should be back to approximately the same number
		// Allow for some variance due to runtime goroutines
		if after > before+5 {
			t.Errorf("Goroutines not cleaned up properly. Before: %d, After: %d", before, after)
		}
	})
}
