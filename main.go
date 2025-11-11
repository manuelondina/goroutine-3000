package main

import (
	"fmt"
	"os"

	"github.com/manuelondina/goroutine-3000/pkg/patterns"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "worker-pool":
		patterns.DemoWorkerPool()
	case "fan-out-fan-in":
		patterns.DemoFanOutFanIn()
	case "pipeline":
		patterns.DemoPipeline()
	case "stress-test":
		patterns.DemoStressTest()
	case "context":
		patterns.DemoContext()
	case "error-handling":
		patterns.DemoErrorHandling()
	case "all":
		fmt.Println("=== Running All Goroutine Demos ===")
		fmt.Println()
		patterns.DemoWorkerPool()
		fmt.Println()
		patterns.DemoFanOutFanIn()
		fmt.Println()
		patterns.DemoPipeline()
		fmt.Println()
		patterns.DemoStressTest()
		fmt.Println()
		patterns.DemoContext()
		fmt.Println()
		patterns.DemoErrorHandling()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("goroutine-3000 - Exploit and test goroutines to max capabilities")
	fmt.Println("\nUsage:")
	fmt.Println("  goroutine-3000 <command>")
	fmt.Println("\nCommands:")
	fmt.Println("  worker-pool       - Demonstrate worker pool pattern")
	fmt.Println("  fan-out-fan-in    - Demonstrate fan-out/fan-in pattern")
	fmt.Println("  pipeline          - Demonstrate pipeline pattern")
	fmt.Println("  stress-test       - Stress test with thousands of goroutines")
	fmt.Println("  context           - Demonstrate context-based cancellation")
	fmt.Println("  error-handling    - Demonstrate error handling in goroutines")
	fmt.Println("  all               - Run all demonstrations")
}
