package ports

import "github.com/manuelondina/goroutine-3000/internal/domain"

// OutputPort defines the interface for output operations
type OutputPort interface {
	DisplayHeader(patternName string)
	DisplayProgress(workerID int, jobID int, message string)
	DisplayResult(result domain.Result)
	DisplayExecutionResult(result *domain.ExecutionResult)
	DisplayError(err error)
	DisplayMessage(message string)
}

// Logger defines the interface for logging operations
type Logger interface {
	Info(message string)
	Error(message string, err error)
	Debug(message string)
}
