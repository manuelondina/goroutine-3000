package patterns

import (
	"fmt"
	"testing"
)

func TestErrorHandling(t *testing.T) {
	t.Run("error channel collects errors", func(t *testing.T) {
		errorChan := make(chan error, 3)

		// Simulate some errors
		go func() {
			errorChan <- fmt.Errorf("error 1")
			errorChan <- fmt.Errorf("error 2")
			errorChan <- fmt.Errorf("error 3")
			close(errorChan)
		}()

		var errors []error
		for err := range errorChan {
			errors = append(errors, err)
		}

		if len(errors) != 3 {
			t.Errorf("Expected 3 errors, got %d", len(errors))
		}
	})

	t.Run("first error pattern works", func(t *testing.T) {
		errorChan := make(chan error, 1)
		sent := 0

		// Try to send multiple errors
		for i := 1; i <= 5; i++ {
			select {
			case errorChan <- fmt.Errorf("error %d", i):
				sent++
			default:
				// Channel full, skip
			}
		}
		close(errorChan)

		// Should only have 1 error due to buffer size of 1
		if sent != 1 {
			t.Errorf("Expected to send 1 error, sent %d", sent)
		}

		// Should only receive 1 error
		var receivedErrors []error
		for err := range errorChan {
			receivedErrors = append(receivedErrors, err)
		}

		if len(receivedErrors) != 1 {
			t.Errorf("Expected 1 error, got %d", len(receivedErrors))
		}
	})

	t.Run("no errors when all succeed", func(t *testing.T) {
		errorChan := make(chan error, 5)

		// Simulate successful operations
		go func() {
			// No errors sent
			close(errorChan)
		}()

		var errors []error
		for err := range errorChan {
			errors = append(errors, err)
		}

		if len(errors) != 0 {
			t.Errorf("Expected 0 errors, got %d", len(errors))
		}
	})
}
