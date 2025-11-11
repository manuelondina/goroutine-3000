package patterns

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// DemoErrorHandling demonstrates error handling patterns in concurrent operations
func DemoErrorHandling() {
	fmt.Println("=== Error Handling Pattern Demo ===")

	// Demo 1: Collecting errors from multiple goroutines
	fmt.Println("\n1. Collecting errors with error channel:")
	demoErrorChannel()

	// Demo 2: First error wins pattern
	fmt.Println("\n2. First error wins pattern:")
	demoFirstError()

	fmt.Println("\nError handling demonstrations completed!")
}

func demoErrorChannel() {
	const numWorkers = 5
	errorChan := make(chan error, numWorkers)
	var wg sync.WaitGroup

	// Spawn workers that may fail
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			if err := riskyOperation(id); err != nil {
				errorChan <- fmt.Errorf("worker %d: %w", id, err)
			}
		}(i)
	}

	// Close error channel when all workers are done
	go func() {
		wg.Wait()
		close(errorChan)
	}()

	// Collect and report errors
	var errors []error
	for err := range errorChan {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		fmt.Printf("  Encountered %d error(s):\n", len(errors))
		for _, err := range errors {
			fmt.Printf("    - %v\n", err)
		}
	} else {
		fmt.Println("  All workers completed successfully!")
	}
}

func demoFirstError() {
	const numWorkers = 5
	errorChan := make(chan error, 1) // Buffer of 1 for first error
	var wg sync.WaitGroup

	// Spawn workers
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			if err := riskyOperation(id); err != nil {
				select {
				case errorChan <- fmt.Errorf("worker %d: %w", id, err):
					// First error sent
				default:
					// Another error already sent
				}
			}
		}(i)
	}

	// Wait for completion
	go func() {
		wg.Wait()
		close(errorChan)
	}()

	// Get first error if any
	if err := <-errorChan; err != nil {
		fmt.Printf("  First error: %v\n", err)
		fmt.Println("  Ignoring subsequent errors...")
	} else {
		fmt.Println("  All workers completed successfully!")
	}
}

func riskyOperation(id int) error {
	// Seed random with current time + id for variability
	rand.Seed(time.Now().UnixNano() + int64(id))
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)

	// Randomly fail ~30% of the time
	if rand.Float64() < 0.3 {
		return fmt.Errorf("operation failed")
	}
	return nil
}
