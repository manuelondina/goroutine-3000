package application

import (
	"context"
	"fmt"

	"github.com/manuelondina/goroutine-3000/internal/domain"
	"github.com/manuelondina/goroutine-3000/internal/ports"
)

// PatternService orchestrates pattern execution
// Follows Single Responsibility Principle - only manages pattern orchestration
type PatternService struct {
	output   ports.OutputPort
	patterns map[string]ports.PatternExecutor
}

// NewPatternService creates a new pattern service
// Follows Dependency Inversion Principle - depends on abstractions (ports)
func NewPatternService(output ports.OutputPort) *PatternService {
	return &PatternService{
		output:   output,
		patterns: make(map[string]ports.PatternExecutor),
	}
}

// RegisterPattern registers a pattern executor
// Follows Open/Closed Principle - open for extension via registration
func (s *PatternService) RegisterPattern(name string, executor ports.PatternExecutor) {
	s.patterns[name] = executor
}

// ExecutePattern executes a registered pattern by name
func (s *PatternService) ExecutePattern(ctx context.Context, name string, config domain.PatternConfig) error {
	executor, exists := s.patterns[name]
	if !exists {
		return fmt.Errorf("pattern '%s' not found", name)
	}

	s.output.DisplayHeader(executor.Name())
	s.output.DisplayMessage(executor.Description())
	s.output.DisplayMessage("")

	result, err := executor.Execute(ctx, config)
	if err != nil {
		s.output.DisplayError(err)
		return err
	}

	s.output.DisplayExecutionResult(result)
	return nil
}

// ExecuteAll executes all registered patterns
func (s *PatternService) ExecuteAll(ctx context.Context, config domain.PatternConfig) error {
	s.output.DisplayMessage("=== Running All Patterns ===")
	s.output.DisplayMessage("")

	for name := range s.patterns {
		if err := s.ExecutePattern(ctx, name, config); err != nil {
			return err
		}
		s.output.DisplayMessage("")
	}

	return nil
}

// ListPatterns returns the names of all registered patterns
func (s *PatternService) ListPatterns() []string {
	names := make([]string, 0, len(s.patterns))
	for name := range s.patterns {
		names = append(names, name)
	}
	return names
}

// GetPattern returns a pattern executor by name
func (s *PatternService) GetPattern(name string) (ports.PatternExecutor, error) {
	executor, exists := s.patterns[name]
	if !exists {
		return nil, fmt.Errorf("pattern '%s' not found", name)
	}
	return executor, nil
}
