# goroutine-3000
Overclocking goroutines to its max capability

A comprehensive Go project that demonstrates and tests goroutines to their maximum capabilities. This project showcases various concurrency patterns, stress testing, and best practices for working with goroutines in Go.

## Features

- **Worker Pool Pattern**: Efficient task processing with a fixed number of workers
- **Fan-Out/Fan-In Pattern**: Parallel work distribution and result collection
- **Pipeline Pattern**: Multi-stage concurrent data processing
- **Stress Testing**: Spawn and manage thousands of goroutines simultaneously
- **Context-Based Cancellation**: Proper timeout and cancellation handling
- **Error Handling**: Patterns for managing errors in concurrent operations
- **Comprehensive Tests**: Unit tests and benchmarks for all patterns

## Installation

```bash
# Clone the repository
git clone https://github.com/manuelondina/goroutine-3000.git
cd goroutine-3000

# Build the project
make build

# Or install directly
go install github.com/manuelondina/goroutine-3000@latest
```

## Usage

### Running All Demonstrations

```bash
make run
# or
./bin/goroutine-3000 all
```

### Running Individual Patterns

```bash
# Worker Pool Pattern
make worker-pool
# or
./bin/goroutine-3000 worker-pool

# Fan-Out/Fan-In Pattern
make fan-out-fan-in
# or
./bin/goroutine-3000 fan-out-fan-in

# Pipeline Pattern
make pipeline
# or
./bin/goroutine-3000 pipeline

# Stress Test (10,000 goroutines)
make stress-test
# or
./bin/goroutine-3000 stress-test

# Context-Based Cancellation
make context
# or
./bin/goroutine-3000 context

# Error Handling Patterns
make error-handling
# or
./bin/goroutine-3000 error-handling
```

## Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run benchmarks
make bench
```

## Patterns Explained

### Worker Pool

Efficiently processes multiple jobs using a fixed number of worker goroutines. Ideal for rate-limiting and controlled resource usage.

### Fan-Out/Fan-In

Distributes work across multiple goroutines (fan-out) and collects results into a single channel (fan-in). Perfect for parallel processing of independent tasks.

### Pipeline

Chains multiple processing stages together, with each stage running concurrently. Data flows through the pipeline from one stage to the next.

### Stress Test

Demonstrates the ability to spawn and manage thousands of goroutines simultaneously, showcasing Go's lightweight concurrency model.

### Context-Based Cancellation

Shows proper timeout, deadline, and cancellation handling using Go's context package. Essential for building robust concurrent applications.

### Error Handling

Demonstrates patterns for collecting and managing errors from multiple goroutines, including "first error wins" and "collect all errors" approaches.

## Makefile Commands

- `make help` - Show all available commands
- `make build` - Build the project
- `make run` - Run all demonstrations
- `make test` - Run tests
- `make bench` - Run benchmarks
- `make clean` - Clean build artifacts
- `make fmt` - Format code
- `make vet` - Run go vet
- `make lint` - Run all linters

## Architecture

This project follows **Hexagonal Architecture** (Ports and Adapters) with **SOLID principles**:

- ✅ Clean separation of concerns
- ✅ Easy to test with mockable interfaces
- ✅ Easy to extend with new patterns
- ✅ Dependency injection for flexibility

**See [ARCHITECTURE.md](ARCHITECTURE.md) for detailed documentation.**

## Project Structure

```
goroutine-3000/
├── main.go                          # Dependency injection & wiring
├── internal/
│   ├── domain/                      # Core business types (pure Go)
│   │   └── types.go
│   ├── ports/                       # Interfaces (contracts)
│   │   ├── pattern.go
│   │   └── output.go
│   ├── application/                 # Business orchestration
│   │   └── pattern_service.go
│   └── adapters/                    # Implementations
│       ├── cli/                     # CLI adapter
│       ├── output/                  # Output adapters
│       └── patterns/                # Pattern implementations
├── pkg/patterns/                    # Legacy code (deprecated)
├── Makefile                         # Build and test automation
├── go.mod                           # Go module definition
├── ARCHITECTURE.md                  # Architecture documentation
└── README.md                        # This file
```

## Performance

The stress test demonstrates Go's ability to handle massive concurrency:

- Spawns 10,000 goroutines simultaneously
- Performs atomic operations across all goroutines
- Completes in milliseconds on modern hardware
- Efficient memory usage due to goroutine's lightweight nature

## Requirements

- Go 1.18 or higher

## Contributing

Contributions are welcome! Feel free to submit issues or pull requests.

## License

See LICENSE file for details.
