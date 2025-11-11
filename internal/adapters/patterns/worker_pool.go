package patterns

import (
	"context"
	"sync"
	"time"

	"github.com/manuelondina/goroutine-3000/internal/domain"
	"github.com/manuelondina/goroutine-3000/internal/ports"
)

// WorkerPool implements the WorkerPoolExecutor interface
type WorkerPool struct {
	output ports.OutputPort
}

// NewWorkerPool creates a new worker pool executor
func NewWorkerPool(output ports.OutputPort) ports.WorkerPoolExecutor {
	return &WorkerPool{output: output}
}

func (w *WorkerPool) Name() string {
	return "Worker Pool Pattern"
}

func (w *WorkerPool) Description() string {
	return "Efficiently processes multiple jobs using a fixed number of worker goroutines"
}

func (w *WorkerPool) Execute(ctx context.Context, config domain.PatternConfig) (*domain.ExecutionResult, error) {
	startTime := time.Now()
	
	// Create jobs
	jobs := make([]domain.Job, config.NumJobs)
	for i := 0; i < config.NumJobs; i++ {
		jobs[i] = domain.Job{ID: i + 1, Data: i + 1}
	}

	// Process jobs
	results, err := w.ProcessJobs(ctx, jobs, config.NumWorkers)
	if err != nil {
		return nil, err
	}

	duration := time.Since(startTime)
	
	// Count successes and failures
	successful := 0
	failed := 0
	for _, result := range results {
		if result.Error != nil {
			failed++
		} else {
			successful++
		}
	}

	return &domain.ExecutionResult{
		TotalJobs:      config.NumJobs,
		SuccessfulJobs: successful,
		FailedJobs:     failed,
		DurationMs:     duration.Milliseconds(),
	}, nil
}

func (w *WorkerPool) ProcessJobs(ctx context.Context, jobs []domain.Job, numWorkers int) ([]domain.Result, error) {
	jobsChan := make(chan domain.Job, len(jobs))
	resultsChan := make(chan domain.Result, len(jobs))
	
	var wg sync.WaitGroup

	// Start workers
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go w.worker(ctx, i, jobsChan, resultsChan, &wg)
	}

	// Send jobs
	go func() {
		for _, job := range jobs {
			select {
			case <-ctx.Done():
				close(jobsChan)
				return
			case jobsChan <- job:
			}
		}
		close(jobsChan)
	}()

	// Wait for workers to finish
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// Collect results
	var results []domain.Result
	for result := range resultsChan {
		results = append(results, result)
		w.output.DisplayResult(result)
	}

	return results, nil
}

func (w *WorkerPool) worker(ctx context.Context, id int, jobs <-chan domain.Job, results chan<- domain.Result, wg *sync.WaitGroup) {
	defer wg.Done()
	
	for job := range jobs {
		select {
		case <-ctx.Done():
			return
		default:
			w.output.DisplayProgress(id, job.ID, "processing")
			
			// Simulate work
			time.Sleep(100 * time.Millisecond)
			
			// Process job (example: multiply by 2)
			value := job.Data.(int) * 2
			results <- domain.Result{
				JobID: job.ID,
				Value: value,
				Error: nil,
			}
		}
	}
}
