package main

import (
	"github.com/manuelondina/goroutine-3000/internal/adapters/cli"
	"github.com/manuelondina/goroutine-3000/internal/adapters/output"
	"github.com/manuelondina/goroutine-3000/internal/adapters/patterns"
	"github.com/manuelondina/goroutine-3000/internal/application"
)

// main is the entry point that wires all dependencies together
// Follows Dependency Inversion Principle - wiring happens here
func main() {
	// Initialize output adapter (infrastructure layer)
	consoleOutput := output.NewConsoleOutput(true)

	// Initialize application service (application layer)
	patternService := application.NewPatternService(consoleOutput)

	// Initialize and register pattern executors (adapters layer)
	// Follows Open/Closed Principle - easy to add new patterns
	workerPool := patterns.NewWorkerPool(consoleOutput)
	fanOutFanIn := patterns.NewFanOutFanIn(consoleOutput)
	pipeline := patterns.NewPipeline(consoleOutput)
	stressTest := patterns.NewStressTest(consoleOutput)

	patternService.RegisterPattern("worker-pool", workerPool)
	patternService.RegisterPattern("fan-out-fan-in", fanOutFanIn)
	patternService.RegisterPattern("pipeline", pipeline)
	patternService.RegisterPattern("stress-test", stressTest)

	// Initialize and run CLI handler (adapter layer)
	cliHandler := cli.NewHandler(patternService)
	cliHandler.Run()
}
