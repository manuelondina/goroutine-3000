package patterns

import (
	"fmt"
	"sync"
	"time"
)

// DemoWorkerPool demonstrates a worker pool pattern where a fixed number of workers
// process jobs from a queue concurrently
func DemoWorkerPool() {
	fmt.Println("=== Worker Pool Pattern Demo ===")

	const numWorkers = 5
	const numJobs = 20

	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)

	// Start workers
	var wg sync.WaitGroup
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

	// Wait for workers to finish and close results
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	fmt.Printf("Processing %d jobs with %d workers:\n", numJobs, numWorkers)
	for result := range results {
		fmt.Printf("  Result: %d\n", result)
	}
	fmt.Println("All jobs completed!")
}

func worker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		fmt.Printf("  Worker %d processing job %d\n", id, job)
		time.Sleep(100 * time.Millisecond) // Simulate work
		results <- job * 2
	}
}
