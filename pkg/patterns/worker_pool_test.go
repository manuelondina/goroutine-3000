package patterns

import (
	"sync"
	"testing"
	"time"
)

func TestWorkerPool(t *testing.T) {
	t.Run("processes all jobs", func(t *testing.T) {
		const numWorkers = 3
		const numJobs = 10

		jobs := make(chan int, numJobs)
		results := make(chan int, numJobs)
		var wg sync.WaitGroup

		// Start workers
		for w := 1; w <= numWorkers; w++ {
			wg.Add(1)
			go worker(w, jobs, results, &wg)
		}

		// Send jobs
		go func() {
			for j := 1; j <= numJobs; j++ {
				jobs <- j
			}
			close(jobs)
		}()

		// Close results when done
		go func() {
			wg.Wait()
			close(results)
		}()

		// Collect results
		count := 0
		for range results {
			count++
		}

		if count != numJobs {
			t.Errorf("Expected %d results, got %d", numJobs, count)
		}
	})

	t.Run("workers process jobs concurrently", func(t *testing.T) {
		const numWorkers = 5
		const numJobs = 20

		jobs := make(chan int, numJobs)
		results := make(chan int, numJobs)
		var wg sync.WaitGroup

		start := time.Now()

		// Start workers
		for w := 1; w <= numWorkers; w++ {
			wg.Add(1)
			go worker(w, jobs, results, &wg)
		}

		// Send jobs
		go func() {
			for j := 1; j <= numJobs; j++ {
				jobs <- j
			}
			close(jobs)
		}()

		// Close results when done
		go func() {
			wg.Wait()
			close(results)
		}()

		// Drain results
		for range results {
		}

		duration := time.Since(start)

		// With 5 workers processing 20 jobs (each taking 100ms),
		// it should take ~400ms (20/5 * 100ms), not 2000ms (20 * 100ms)
		maxExpected := 800 * time.Millisecond // Give some buffer
		if duration > maxExpected {
			t.Errorf("Processing took too long: %v (expected < %v)", duration, maxExpected)
		}
	})
}
