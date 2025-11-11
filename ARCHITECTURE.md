# Architecture Documentation

## Overview

This project follows **Hexagonal Architecture** (also known as Ports and Adapters) combined with **SOLID principles** to achieve:

- ✅ **Clean separation of concerns**
- ✅ **Testability** - Easy to mock and test each layer
- ✅ **Maintainability** - Changes in one layer don't affect others
- ✅ **Extensibility** - Easy to add new patterns or change implementations
- ✅ **Flexibility** - Can swap output adapters (CLI, web, API) without changing business logic

## Architecture Layers

```
┌─────────────────────────────────────────────────────────────┐
│                          main.go                             │
│                    (Dependency Injection)                    │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    Adapters Layer (Input)                    │
│                  internal/adapters/cli/                      │
│                   (CLI Handler)                              │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    Application Layer                         │
│                 internal/application/                        │
│              (Pattern Service - Orchestration)               │
└─────────────────────────────────────────────────────────────┘
                              │
                    ┌─────────┴─────────┐
                    ▼                   ▼
┌──────────────────────────┐  ┌──────────────────────────┐
│    Ports (Interfaces)    │  │  Adapters Layer (Output) │
│     internal/ports/      │  │  internal/adapters/      │
│  - PatternExecutor       │  │  - patterns/             │
│  - OutputPort            │  │  - output/               │
└──────────────────────────┘  └──────────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────────────────────────┐
│                      Domain Layer                            │
│                   internal/domain/                           │
│           (Core Business Types - Pure Go)                    │
└─────────────────────────────────────────────────────────────┘
```

## Directory Structure

```
goroutine-3000/
├── main.go                          # Dependency injection & wiring
├── internal/
│   ├── domain/                      # Domain Layer (Core)
│   │   └── types.go                 # Business entities and value objects
│   │
│   ├── ports/                       # Ports (Interfaces)
│   │   ├── pattern.go               # Pattern executor interfaces
│   │   └── output.go                # Output interfaces
│   │
│   ├── application/                 # Application Services
│   │   └── pattern_service.go       # Pattern orchestration
│   │
│   └── adapters/                    # Adapters (Implementations)
│       ├── cli/                     # CLI adapter (input)
│       │   └── handler.go
│       ├── output/                  # Console output adapter
│       │   └── console.go
│       └── patterns/                # Pattern implementations
│           ├── worker_pool.go
│           ├── pipeline.go
│           ├── fan_out_fan_in.go
│           └── stress.go
│
└── pkg/patterns/                    # Legacy code (to be deprecated)
```

## SOLID Principles Applied

### 1. Single Responsibility Principle (SRP)
Each component has one reason to change:
- **Domain types** (`internal/domain/types.go`): Only change when business concepts change
- **Pattern executors**: Only implement one concurrency pattern
- **Output adapter**: Only handles display logic
- **CLI handler**: Only handles command-line parsing
- **Pattern service**: Only orchestrates pattern execution

**Example:**
```go
// PatternService has single responsibility: orchestrate patterns
type PatternService struct {
    output   ports.OutputPort
    patterns map[string]ports.PatternExecutor
}
```

### 2. Open/Closed Principle (OCP)
Open for extension, closed for modification:
- **Adding new patterns**: Just implement `PatternExecutor` interface and register it
- **Adding new outputs**: Implement `OutputPort` interface (e.g., JSON, Web UI)
- **No need to modify existing code**

**Example:**
```go
// To add a new pattern, just implement the interface:
type NewPattern struct {
    output ports.OutputPort
}

func (n *NewPattern) Execute(ctx context.Context, config domain.PatternConfig) (*domain.ExecutionResult, error) {
    // Implementation
}

// Register it in main.go:
newPattern := patterns.NewPattern(output)
service.RegisterPattern("new-pattern", newPattern)
```

### 3. Liskov Substitution Principle (LSP)
Any implementation of an interface can be substituted:
- All `PatternExecutor` implementations are interchangeable
- All `OutputPort` implementations can replace each other
- Code depends on interfaces, not concrete types

**Example:**
```go
// Any PatternExecutor can be used here
func (s *PatternService) ExecutePattern(ctx context.Context, name string, config domain.PatternConfig) error {
    executor, _ := s.patterns[name]  // Returns PatternExecutor interface
    result, err := executor.Execute(ctx, config)  // Works with any implementation
    // ...
}
```

### 4. Interface Segregation Principle (ISP)
Interfaces are small and focused:
- `PatternExecutor`: Basic pattern execution
- `WorkerPoolExecutor`: Specific to worker pools
- `PipelineExecutor`: Specific to pipelines
- `OutputPort`: Only output operations
- `Logger`: Only logging operations

**Example:**
```go
// Base interface - minimal contract
type PatternExecutor interface {
    Execute(ctx context.Context, config domain.PatternConfig) (*domain.ExecutionResult, error)
    Name() string
    Description() string
}

// Extended interface for specific needs
type WorkerPoolExecutor interface {
    PatternExecutor
    ProcessJobs(ctx context.Context, jobs []domain.Job, numWorkers int) ([]domain.Result, error)
}
```

### 5. Dependency Inversion Principle (DIP)
High-level modules don't depend on low-level modules. Both depend on abstractions:
- Application layer depends on `ports.PatternExecutor` (interface)
- Adapters implement the interfaces
- `main.go` wires everything together

**Example:**
```go
// High-level service depends on abstraction (OutputPort), not concrete implementation
type PatternService struct {
    output   ports.OutputPort  // Interface, not *ConsoleOutput
}

// Concrete implementation
type ConsoleOutput struct {
    verbose bool
}

func (c *ConsoleOutput) DisplayHeader(patternName string) {
    fmt.Printf("=== %s ===\n", patternName)
}

// Dependency injection in main.go
consoleOutput := output.NewConsoleOutput(true)
service := application.NewPatternService(consoleOutput)
```

## Hexagonal Architecture Benefits

### 1. **Domain Isolation**
The domain layer (`internal/domain/`) contains only pure business logic with no external dependencies. This makes it:
- Easy to test
- Easy to reason about
- Independent of frameworks and UI

### 2. **Port-Adapter Pattern**
- **Ports** (`internal/ports/`): Define contracts (interfaces)
- **Adapters** (`internal/adapters/`): Implement the contracts
- **Application** (`internal/application/`): Uses ports to orchestrate

### 3. **Testability**
Each layer can be tested independently:
```go
// Test pattern executor with mock output
mockOutput := &MockOutputPort{}
workerPool := patterns.NewWorkerPool(mockOutput)
result, err := workerPool.Execute(ctx, config)
```

### 4. **Flexibility**
Easy to swap implementations:
- Replace `ConsoleOutput` with `JSONOutput` or `WebOutput`
- Replace `CLIHandler` with `HTTPHandler` or `GRPCHandler`
- No changes to business logic required

## Adding a New Pattern

1. **Implement the interface** in `internal/adapters/patterns/`:
```go
type MyNewPattern struct {
    output ports.OutputPort
}

func NewMyNewPattern(output ports.OutputPort) ports.PatternExecutor {
    return &MyNewPattern{output: output}
}

func (m *MyNewPattern) Execute(ctx context.Context, config domain.PatternConfig) (*domain.ExecutionResult, error) {
    // Implementation
}

func (m *MyNewPattern) Name() string {
    return "My New Pattern"
}

func (m *MyNewPattern) Description() string {
    return "Description of my pattern"
}
```

2. **Register in main.go**:
```go
myPattern := patterns.NewMyNewPattern(consoleOutput)
patternService.RegisterPattern("my-pattern", myPattern)
```

3. **Add CLI command** in `internal/adapters/cli/handler.go`:
```go
case "my-pattern":
    return h.service.ExecutePattern(ctx, "my-pattern", config)
```

That's it! No changes to existing code.

## Adding a New Output Adapter

1. **Implement OutputPort interface**:
```go
type JSONOutput struct{}

func (j *JSONOutput) DisplayHeader(patternName string) {
    fmt.Printf("{\"header\": \"%s\"}\n", patternName)
}
// ... implement other methods
```

2. **Use it in main.go**:
```go
jsonOutput := output.NewJSONOutput()
patternService := application.NewPatternService(jsonOutput)
```

## Testing Strategy

### Unit Tests
- Test each adapter independently with mocks
- Test application service with mock ports
- Test domain types in isolation

### Integration Tests
- Test full stack with real implementations
- Test different pattern executors together

### Example Test
```go
func TestWorkerPool_Execute(t *testing.T) {
    mockOutput := &MockOutputPort{}
    wp := patterns.NewWorkerPool(mockOutput)
    
    config := domain.PatternConfig{
        NumWorkers: 3,
        NumJobs:    10,
    }
    
    result, err := wp.Execute(context.Background(), config)
    
    assert.NoError(t, err)
    assert.Equal(t, 10, result.SuccessfulJobs)
}
```

## Comparison: Before vs After

### Before (Monolithic)
```
pkg/patterns/
├── worker_pool.go        # Mixed concerns: logic + I/O
├── pipeline.go           # Direct fmt.Println calls
├── fan_out_fan_in.go     # Hard to test
└── stress.go             # Tightly coupled

main.go                   # Direct calls to demos
```

**Issues:**
- Business logic mixed with I/O
- Hard to test (can't mock output)
- Hard to extend (must modify existing code)
- No clear boundaries
- Violates multiple SOLID principles

### After (Hexagonal + SOLID)
```
internal/
├── domain/               # Pure business types
├── ports/                # Interfaces (contracts)
├── application/          # Orchestration
└── adapters/             # Implementations
    ├── cli/              # CLI adapter
    ├── output/           # Output adapters
    └── patterns/         # Pattern implementations

main.go                   # Dependency injection
```

**Benefits:**
- ✅ Clear separation of concerns
- ✅ Easy to test (mockable interfaces)
- ✅ Easy to extend (implement interfaces)
- ✅ Clear boundaries between layers
- ✅ Follows all SOLID principles
- ✅ Business logic independent of frameworks

## Key Takeaways

1. **Hexagonal Architecture** provides clear boundaries between core business logic and external concerns
2. **SOLID principles** ensure the code is maintainable, extensible, and testable
3. **Dependency Injection** in `main.go` wires everything together without coupling
4. **Interfaces (Ports)** define contracts that adapters implement
5. **Domain layer** contains only pure business logic with zero external dependencies
6. **Easy to extend**: Add new patterns or outputs without modifying existing code
7. **Easy to test**: Every layer can be tested independently with mocks

This architecture scales well and makes the codebase professional, maintainable, and a joy to work with! 🚀
