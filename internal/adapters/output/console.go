package output

import (
	"fmt"

	"github.com/manuelondina/goroutine-3000/internal/domain"
	"github.com/manuelondina/goroutine-3000/internal/ports"
)

// ConsoleOutput implements the OutputPort interface for console output
type ConsoleOutput struct {
	verbose bool
}

// NewConsoleOutput creates a new console output adapter
func NewConsoleOutput(verbose bool) ports.OutputPort {
	return &ConsoleOutput{verbose: verbose}
}

func (c *ConsoleOutput) DisplayHeader(patternName string) {
	fmt.Printf("=== %s ===\n", patternName)
}

func (c *ConsoleOutput) DisplayProgress(workerID int, jobID int, message string) {
	if c.verbose {
		if workerID > 0 {
			fmt.Printf("  Worker %d processing job %d: %s\n", workerID, jobID, message)
		} else {
			fmt.Printf("  Processing job %d: %s\n", jobID, message)
		}
	}
}

func (c *ConsoleOutput) DisplayResult(result domain.Result) {
	if result.Error != nil {
		fmt.Printf("  Job %d failed: %v\n", result.JobID, result.Error)
	} else {
		fmt.Printf("  Result: %v\n", result.Value)
	}
}

func (c *ConsoleOutput) DisplayExecutionResult(result *domain.ExecutionResult) {
	fmt.Println("\nExecution Summary:")
	fmt.Printf("  Total Jobs: %d\n", result.TotalJobs)
	fmt.Printf("  Successful: %d\n", result.SuccessfulJobs)
	fmt.Printf("  Failed: %d\n", result.FailedJobs)
	fmt.Printf("  Duration: %d ms\n", result.DurationMs)
	if result.OperationsCount > 0 {
		fmt.Printf("  Total Operations: %d\n", result.OperationsCount)
		opsPerSec := float64(result.OperationsCount) / (float64(result.DurationMs) / 1000.0)
		fmt.Printf("  Operations/second: %.2f\n", opsPerSec)
	}
}

func (c *ConsoleOutput) DisplayError(err error) {
	fmt.Printf("Error: %v\n", err)
}

func (c *ConsoleOutput) DisplayMessage(message string) {
	fmt.Println(message)
}
