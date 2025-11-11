package patterns

import (
	"context"
	"time"

	"github.com/manuelondina/goroutine-3000/internal/domain"
	"github.com/manuelondina/goroutine-3000/internal/ports"
)

// Pipeline implements the PipelineExecutor interface
type Pipeline struct {
	output ports.OutputPort
}

// NewPipeline creates a new pipeline executor
func NewPipeline(output ports.OutputPort) ports.PipelineExecutor {
	return &Pipeline{output: output}
}

func (p *Pipeline) Name() string {
	return "Pipeline Pattern"
}

func (p *Pipeline) Description() string {
	return "Chains multiple processing stages together with concurrent data flow"
}

func (p *Pipeline) Execute(ctx context.Context, config domain.PatternConfig) (*domain.ExecutionResult, error) {
	startTime := time.Now()

	// Create input data
	input := make([]int, config.NumJobs)
	for i := 0; i < config.NumJobs; i++ {
		input[i] = i + 1
	}

	// Build and run pipeline
	outputChan := p.BuildPipeline(ctx, input)

	// Collect results
	count := 0
	for result := range outputChan {
		count++
		p.output.DisplayResult(domain.Result{
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

func (p *Pipeline) BuildPipeline(ctx context.Context, input []int) <-chan int {
	// Stage 1: Generate numbers
	generatorChan := p.generator(ctx, input)

	// Stage 2: Square the numbers
	squaredChan := p.square(ctx, generatorChan)

	// Stage 3: Add 1 to each number
	addedChan := p.addOne(ctx, squaredChan)

	// Stage 4: Filter even numbers
	filteredChan := p.filterEven(ctx, addedChan)

	return filteredChan
}

func (p *Pipeline) generator(ctx context.Context, numbers []int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range numbers {
			select {
			case <-ctx.Done():
				return
			case out <- n:
			}
		}
	}()
	return out
}

func (p *Pipeline) square(ctx context.Context, input <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range input {
			select {
			case <-ctx.Done():
				return
			case out <- n * n:
			}
		}
	}()
	return out
}

func (p *Pipeline) addOne(ctx context.Context, input <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range input {
			select {
			case <-ctx.Done():
				return
			case out <- n + 1:
			}
		}
	}()
	return out
}

func (p *Pipeline) filterEven(ctx context.Context, input <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range input {
			select {
			case <-ctx.Done():
				return
			default:
				if n%2 == 0 {
					out <- n
				}
			}
		}
	}()
	return out
}
