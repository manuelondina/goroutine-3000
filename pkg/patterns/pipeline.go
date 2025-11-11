package patterns

import (
	"fmt"
)

// DemoPipeline demonstrates a pipeline pattern where data flows through
// multiple stages, each processing the data concurrently
func DemoPipeline() {
	fmt.Println("=== Pipeline Pattern Demo ===")

	// Input numbers
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// Stage 1: Generate numbers
	inputChan := generator(numbers)

	// Stage 2: Square the numbers
	squaredChan := square(inputChan)

	// Stage 3: Add 1 to each number
	addedChan := addOne(squaredChan)

	// Stage 4: Filter even numbers
	filteredChan := filterEven(addedChan)

	// Consume results
	fmt.Printf("Pipeline: numbers -> square -> add 1 -> filter even\n")
	fmt.Println("Results:")
	for result := range filteredChan {
		fmt.Printf("  %d\n", result)
	}
	fmt.Println("Pipeline completed!")
}

func generator(numbers []int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range numbers {
			out <- n
		}
	}()
	return out
}

func square(input <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range input {
			out <- n * n
		}
	}()
	return out
}

func addOne(input <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range input {
			out <- n + 1
		}
	}()
	return out
}

func filterEven(input <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range input {
			if n%2 == 0 {
				out <- n
			}
		}
	}()
	return out
}
