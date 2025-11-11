package ports

import (
	"context"

	"github.com/manuelondina/goroutine-3000/internal/domain"
)

// PatternExecutor defines the interface for executing concurrency patterns
type PatternExecutor interface {
	Execute(ctx context.Context, config domain.PatternConfig) (*domain.ExecutionResult, error)
	Name() string
	Description() string
}

// WorkerPoolExecutor defines the interface for worker pool pattern
type WorkerPoolExecutor interface {
	PatternExecutor
	ProcessJobs(ctx context.Context, jobs []domain.Job, numWorkers int) ([]domain.Result, error)
}

// PipelineExecutor defines the interface for pipeline pattern
type PipelineExecutor interface {
	PatternExecutor
	BuildPipeline(ctx context.Context, input []int) <-chan int
}

// FanOutFanInExecutor defines the interface for fan-out/fan-in pattern
type FanOutFanInExecutor interface {
	PatternExecutor
	FanOut(ctx context.Context, input []int, numWorkers int) []<-chan int
	FanIn(ctx context.Context, channels ...<-chan int) <-chan int
}

// StressTestExecutor defines the interface for stress testing
type StressTestExecutor interface {
	PatternExecutor
	RunStressTest(ctx context.Context, numGoroutines int) (*domain.ExecutionResult, error)
}
