package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/manuelondina/goroutine-3000/internal/application"
	"github.com/manuelondina/goroutine-3000/internal/domain"
)

// Handler handles CLI commands
// Follows Single Responsibility Principle - only handles CLI interaction
type Handler struct {
	service *application.PatternService
}

// NewHandler creates a new CLI handler
func NewHandler(service *application.PatternService) *Handler {
	return &Handler{service: service}
}

// Handle processes command line arguments and executes the appropriate action
func (h *Handler) Handle(args []string) error {
	if len(args) < 2 {
		h.printUsage()
		return fmt.Errorf("insufficient arguments")
	}

	command := args[1]
	ctx := context.Background()

	// Default configuration
	config := domain.PatternConfig{
		NumWorkers:    5,
		NumJobs:       20,
		TimeoutMs:     0,
		EnableLogging: true,
	}

	switch command {
	case "worker-pool":
		return h.service.ExecutePattern(ctx, "worker-pool", config)
	case "fan-out-fan-in":
		config.NumWorkers = 3
		config.NumJobs = 10
		return h.service.ExecutePattern(ctx, "fan-out-fan-in", config)
	case "pipeline":
		config.NumJobs = 10
		return h.service.ExecutePattern(ctx, "pipeline", config)
	case "stress-test":
		config.NumJobs = 10000
		return h.service.ExecutePattern(ctx, "stress-test", config)
	case "all":
		return h.service.ExecuteAll(ctx, config)
	default:
		h.printUsage()
		return fmt.Errorf("unknown command: %s", command)
	}
}

func (h *Handler) printUsage() {
	fmt.Println("goroutine-3000 - Exploit and test goroutines to max capabilities")
	fmt.Println("\nUsage:")
	fmt.Println("  goroutine-3000 <command>")
	fmt.Println("\nCommands:")
	fmt.Println("  worker-pool       - Demonstrate worker pool pattern")
	fmt.Println("  fan-out-fan-in    - Demonstrate fan-out/fan-in pattern")
	fmt.Println("  pipeline          - Demonstrate pipeline pattern")
	fmt.Println("  stress-test       - Stress test with thousands of goroutines")
	fmt.Println("  all               - Run all demonstrations")
}

// Run is the main entry point for the CLI
func (h *Handler) Run() {
	if err := h.Handle(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
