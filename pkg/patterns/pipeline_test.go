package patterns

import (
	"testing"
)

func TestPipeline(t *testing.T) {
	t.Run("generator produces all numbers", func(t *testing.T) {
		numbers := []int{1, 2, 3, 4, 5}
		output := generator(numbers)

		results := []int{}
		for val := range output {
			results = append(results, val)
		}

		if len(results) != len(numbers) {
			t.Errorf("Expected %d numbers, got %d", len(numbers), len(results))
		}

		for i, val := range results {
			if val != numbers[i] {
				t.Errorf("Expected %d at position %d, got %d", numbers[i], i, val)
			}
		}
	})

	t.Run("square stage works correctly", func(t *testing.T) {
		input := make(chan int)
		go func() {
			input <- 2
			input <- 3
			input <- 4
			close(input)
		}()

		output := square(input)
		expected := []int{4, 9, 16}
		results := []int{}

		for val := range output {
			results = append(results, val)
		}

		if len(results) != len(expected) {
			t.Errorf("Expected %d results, got %d", len(expected), len(results))
		}

		for i, val := range results {
			if val != expected[i] {
				t.Errorf("Expected %d at position %d, got %d", expected[i], i, val)
			}
		}
	})

	t.Run("addOne stage works correctly", func(t *testing.T) {
		input := make(chan int)
		go func() {
			input <- 1
			input <- 5
			input <- 10
			close(input)
		}()

		output := addOne(input)
		expected := []int{2, 6, 11}
		results := []int{}

		for val := range output {
			results = append(results, val)
		}

		if len(results) != len(expected) {
			t.Errorf("Expected %d results, got %d", len(expected), len(results))
		}

		for i, val := range results {
			if val != expected[i] {
				t.Errorf("Expected %d at position %d, got %d", expected[i], i, val)
			}
		}
	})

	t.Run("filterEven stage filters correctly", func(t *testing.T) {
		input := make(chan int)
		go func() {
			for i := 1; i <= 10; i++ {
				input <- i
			}
			close(input)
		}()

		output := filterEven(input)
		results := []int{}

		for val := range output {
			results = append(results, val)
		}

		// Should only have even numbers: 2, 4, 6, 8, 10
		expected := []int{2, 4, 6, 8, 10}
		if len(results) != len(expected) {
			t.Errorf("Expected %d results, got %d", len(expected), len(results))
		}

		for i, val := range results {
			if val != expected[i] {
				t.Errorf("Expected %d at position %d, got %d", expected[i], i, val)
			}
			if val%2 != 0 {
				t.Errorf("Expected even number, got odd: %d", val)
			}
		}
	})

	t.Run("full pipeline works end-to-end", func(t *testing.T) {
		numbers := []int{1, 2, 3, 4, 5}

		// Pipeline: square -> add 1 -> filter even
		output := generator(numbers)
		output = square(output)
		output = addOne(output)
		output = filterEven(output)

		results := []int{}
		for val := range output {
			results = append(results, val)
		}

		// 1->1->2 (even), 2->4->5 (odd), 3->9->10 (even), 4->16->17 (odd), 5->25->26 (even)
		expected := []int{2, 10, 26}
		if len(results) != len(expected) {
			t.Errorf("Expected %d results, got %d", len(expected), len(results))
		}

		for i, val := range results {
			if val != expected[i] {
				t.Errorf("Expected %d at position %d, got %d", expected[i], i, val)
			}
		}
	})
}
