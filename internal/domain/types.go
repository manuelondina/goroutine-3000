package domain

import "context"

// Job represents a unit of work to be processed
type Job struct {
	ID   int
	Data interface{}
}

// Result represents the output of a processed job
type Result struct {
	JobID int
	Value interface{}
	Error error
}

// PatternConfig holds configuration for pattern execution
type PatternConfig struct {
	NumWorkers    int
	NumJobs       int
	TimeoutMs     int
	EnableLogging bool
}

// ExecutionResult contains metrics about pattern execution
type ExecutionResult struct {
	TotalJobs       int
	SuccessfulJobs  int
	FailedJobs      int
	DurationMs      int64
	OperationsCount int64
}

// PatternContext wraps context with additional pattern-specific data
type PatternContext struct {
	Ctx    context.Context
	Cancel context.CancelFunc
	Config PatternConfig
}
