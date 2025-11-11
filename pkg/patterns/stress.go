package patterns

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// DemoStressTest demonstrates the ability to spawn and manage thousands of goroutines
func DemoStressTest() {
	fmt.Println("=== Goroutine Stress Test Demo ===")

	const numGoroutines = 10000
	var counter atomic.Int64
	var wg sync.WaitGroup

	startTime := time.Now()
	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)

	fmt.Printf("System info:\n")
	fmt.Printf("  CPU cores: %d\n", numCPU)
	fmt.Printf("  GOMAXPROCS: %d\n", runtime.GOMAXPROCS(0))
	fmt.Printf("\nSpawning %d goroutines...\n", numGoroutines)

	// Spawn thousands of goroutines
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// Simulate some work
			for j := 0; j < 100; j++ {
				counter.Add(1)
			}
		}(i)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	duration := time.Since(startTime)

	fmt.Printf("\nStress test completed!\n")
	fmt.Printf("  Total goroutines: %d\n", numGoroutines)
	fmt.Printf("  Total operations: %d\n", counter.Load())
	fmt.Printf("  Time taken: %v\n", duration)
	fmt.Printf("  Operations/second: %.2f\n", float64(counter.Load())/duration.Seconds())
	fmt.Printf("  Current goroutines: %d\n", runtime.NumGoroutine())
}
