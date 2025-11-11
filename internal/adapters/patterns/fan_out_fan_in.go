package patterns

import (
	"context"
	"sync"
	"time"

	"github.com/manuelondina/goroutine-3000/internal/domain"
	"github.com/manuelondina/goroutine-3000/internal/ports"
)

// FanOutFanIn implements the FanOutFanInExecutor interface
type FanOutFanIn struct {
	output ports.OutputPort
}

// NewFanOutFanIn creates a new fan-out/fan-in executor
func NewFanOutFanIn(output ports.OutputPort) ports.FanOutFanInExecutor {
	return &FanOutFanIn{output: output}
}

func (f *FanOutFanIn) Name() string {
	return "Fan-Out/Fan-In Pattern"
}

func (f *FanOutFanIn) Description() string {
	return "Distributes work across multiple goroutines and collects results"
}

func (f *FanOutFanIn) Execute(ctx context.Context, config domain.PatternConfig) (*domain.ExecutionResult, error) {
	startTime := time.Now()

	// Create input data
	input := make([]int, config.NumJobs)
	for i := 0; i < config.NumJobs; i++ {
		input[i] = i + 1
	}

	// Fan-out to workers
	channels := f.FanOut(ctx, input, config.NumWorkers)

	// Fan-in results
	resultChan := f.FanIn(ctx, channels...)

	// Collect results
	count := 0
	for result := range resultChan {
		count++
		f.output.DisplayResult(domain.Result{
			JobID: count,
			Value: result,
			Error: nil,
		})
	}

	duration := time.Since(startTime)

	return &domain.ExecutionResult{
		TotalJobs:      config.NumJobs,
		SuccessfulJobs: count,
		FailedJobs:     0,
		DurationMs:     duration.Milliseconds(),
	}, nil
}

func (f *FanOutFanIn) FanOut(ctx context.Context, input []int, numWorkers int) []<-chan int {
	inputChan := make(chan int)
	resultChans := make([]<-chan int, numWorkers)

	// Start workers
	for i := 0; i < numWorkers; i++ {
		resultChans[i] = f.squareWorker(ctx, inputChan)
	}

	// Send input
	go func() {
		defer close(inputChan)
		for _, num := range input {
			select {
			case <-ctx.Done():
				return
			case inputChan <- num:
			}
		}
	}()

	return resultChans
}

func (f *FanOutFanIn) FanIn(ctx context.Context, channels ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	merged := make(chan int)

	// Start a goroutine for each input channel
	for _, ch := range channels {
		wg.Add(1)
		go func(c <-chan int) {
			defer wg.Done()
			for val := range c {
				select {
				case <-ctx.Done():
					return
				case merged <- val:
				}
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

func (f *FanOutFanIn) squareWorker(ctx context.Context, input <-chan int) <-chan int {
	output := make(chan int)
	go func() {
		defer close(output)
		for num := range input {
			select {
			case <-ctx.Done():
				return
			case output <- num * num:
			}
		}
	}()
	return output
}
