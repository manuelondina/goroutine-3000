package patterns

import (
	"fmt"
	"sync"
)

// DemoFanOutFanIn demonstrates the fan-out/fan-in pattern where work is distributed
// across multiple goroutines (fan-out) and results are collected (fan-in)
func DemoFanOutFanIn() {
	fmt.Println("=== Fan-Out/Fan-In Pattern Demo ===")

	// Input data
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// Fan-out: distribute work to multiple goroutines
	numWorkers := 3
	inputChan := make(chan int)
	resultChans := make([]<-chan int, numWorkers)

	for i := 0; i < numWorkers; i++ {
		resultChans[i] = squareWorker(inputChan)
	}

	// Send input
	go func() {
		for _, num := range numbers {
			inputChan <- num
		}
		close(inputChan)
	}()

	// Fan-in: merge results from multiple goroutines
	resultChan := fanIn(resultChans...)

	// Collect results
	fmt.Printf("Processing %d numbers with %d workers:\n", len(numbers), numWorkers)
	for result := range resultChan {
		fmt.Printf("  Result: %d\n", result)
	}
	fmt.Println("All computations completed!")
}

func squareWorker(input <-chan int) <-chan int {
	output := make(chan int)
	go func() {
		defer close(output)
		for num := range input {
			output <- num * num
		}
	}()
	return output
}

func fanIn(channels ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	merged := make(chan int)

	// Start a goroutine for each input channel
	for _, ch := range channels {
		wg.Add(1)
		go func(c <-chan int) {
			defer wg.Done()
			for val := range c {
				merged <- val
			}
		}(ch)
	}

	// Close merged channel when all input channels are done
	go func() {
		wg.Wait()
		close(merged)
	}()

	return merged
}
