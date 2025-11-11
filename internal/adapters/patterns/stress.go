package patterns

import (
	"context"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/manuelondina/goroutine-3000/internal/domain"
	"github.com/manuelondina/goroutine-3000/internal/ports"
)

// StressTest implements the StressTestExecutor interface
type StressTest struct {
	output ports.OutputPort
}

// NewStressTest creates a new stress test executor
func NewStressTest(output ports.OutputPort) ports.StressTestExecutor {
	return &StressTest{output: output}
}

func (s *StressTest) Name() string {
	return "Stress Test Pattern"
}

func (s *StressTest) Description() string {
	return "Spawns thousands of goroutines to test concurrency limits"
}

func (s *StressTest) Execute(ctx context.Context, config domain.PatternConfig) (*domain.ExecutionResult, error) {
	numGoroutines := config.NumJobs
	if numGoroutines == 0 {
		numGoroutines = 10000
	}

	return s.RunStressTest(ctx, numGoroutines)
}

func (s *StressTest) RunStressTest(ctx context.Context, numGoroutines int) (*domain.ExecutionResult, error) {
	var counter atomic.Int64
	var wg sync.WaitGroup

	startTime := time.Now()
	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)

	s.output.DisplayMessage("System info:")
	s.output.DisplayMessage("  CPU cores: " + string(rune(numCPU+'0')))
	s.output.DisplayMessage("")
	s.output.DisplayMessage("Spawning goroutines...")

	// Spawn thousands of goroutines
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			select {
			case <-ctx.Done():
				return
			default:
				// Simulate some work
				for j := 0; j < 100; j++ {
					counter.Add(1)
				}
			}
		}(i)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	duration := time.Since(startTime)

	return &domain.ExecutionResult{
		TotalJobs:       numGoroutines,
		SuccessfulJobs:  numGoroutines,
		FailedJobs:      0,
		DurationMs:      duration.Milliseconds(),
		OperationsCount: counter.Load(),
	}, nil
}
