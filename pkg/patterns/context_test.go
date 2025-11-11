package patterns

import (
	"context"
	"testing"
	"time"
)

func TestContext(t *testing.T) {
	t.Run("context timeout cancels operation", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		done := make(chan bool)
		cancelled := false

		go func() {
			ticker := time.NewTicker(50 * time.Millisecond)
			defer ticker.Stop()
			defer func() { done <- true }()

			for {
				select {
				case <-ctx.Done():
					cancelled = true
					return
				case <-ticker.C:
					// Continue working
				}
			}
		}()

		<-done

		if !cancelled {
			t.Error("Expected context to cancel the operation")
		}
	})

	t.Run("manual cancellation stops operation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		done := make(chan bool)
		cancelled := false

		go func() {
			ticker := time.NewTicker(50 * time.Millisecond)
			defer ticker.Stop()
			defer func() { done <- true }()

			for {
				select {
				case <-ctx.Done():
					cancelled = true
					return
				case <-ticker.C:
					// Continue working
				}
			}
		}()

		time.Sleep(150 * time.Millisecond)
		cancel()
		<-done

		if !cancelled {
			t.Error("Expected manual cancellation to stop the operation")
		}
	})

	t.Run("context deadline is respected", func(t *testing.T) {
		deadline := time.Now().Add(100 * time.Millisecond)
		ctx, cancel := context.WithDeadline(context.Background(), deadline)
		defer cancel()

		done := make(chan bool)
		cancelled := false
		start := time.Now()

		go func() {
			ticker := time.NewTicker(50 * time.Millisecond)
			defer ticker.Stop()
			defer func() { done <- true }()

			for {
				select {
				case <-ctx.Done():
					cancelled = true
					return
				case <-ticker.C:
					// Continue working
				}
			}
		}()

		<-done
		duration := time.Since(start)

		if !cancelled {
			t.Error("Expected context deadline to cancel the operation")
		}

		// Should cancel around the deadline time
		if duration < 100*time.Millisecond || duration > 200*time.Millisecond {
			t.Errorf("Expected cancellation around 100ms, got %v", duration)
		}
	})

	t.Run("context error is correct for timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		time.Sleep(50 * time.Millisecond)

		if ctx.Err() != context.DeadlineExceeded {
			t.Errorf("Expected DeadlineExceeded error, got %v", ctx.Err())
		}
	})

	t.Run("context error is correct for cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		time.Sleep(10 * time.Millisecond)

		if ctx.Err() != context.Canceled {
			t.Errorf("Expected Canceled error, got %v", ctx.Err())
		}
	})
}
