package patterns

import (
	"testing"
)

func TestFanOutFanIn(t *testing.T) {
	t.Run("squares all numbers", func(t *testing.T) {
		input := make(chan int, 5)
		numbers := []int{1, 2, 3, 4, 5}
		expected := map[int]bool{1: true, 4: true, 9: true, 16: true, 25: true}

		// Send numbers
		go func() {
			for _, n := range numbers {
				input <- n
			}
			close(input)
		}()

		// Process with worker
		output := squareWorker(input)

		// Collect results
		results := make(map[int]bool)
		for result := range output {
			results[result] = true
		}

		// Check all expected results are present
		for expectedVal := range expected {
			if !results[expectedVal] {
				t.Errorf("Expected result %d not found", expectedVal)
			}
		}

		if len(results) != len(expected) {
			t.Errorf("Expected %d results, got %d", len(expected), len(results))
		}
	})

	t.Run("fanIn merges multiple channels", func(t *testing.T) {
		ch1 := make(chan int)
		ch2 := make(chan int)
		ch3 := make(chan int)

		// Send data to channels
		go func() {
			ch1 <- 1
			ch1 <- 2
			close(ch1)
		}()
		go func() {
			ch2 <- 3
			ch2 <- 4
			close(ch2)
		}()
		go func() {
			ch3 <- 5
			ch3 <- 6
			close(ch3)
		}()

		// Merge channels
		merged := fanIn(ch1, ch2, ch3)

		// Collect results
		results := make(map[int]bool)
		for val := range merged {
			results[val] = true
		}

		// Verify all values received
		expected := []int{1, 2, 3, 4, 5, 6}
		if len(results) != len(expected) {
			t.Errorf("Expected %d results, got %d", len(expected), len(results))
		}

		for _, val := range expected {
			if !results[val] {
				t.Errorf("Expected value %d not found in results", val)
			}
		}
	})
}
