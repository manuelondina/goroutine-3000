.PHONY: help build run test bench clean fmt vet lint all-demos

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the project
	go build -o bin/goroutine-3000 .

run: build ## Build and run all demonstrations
	./bin/goroutine-3000 all

test: ## Run tests
	go test -v ./...

test-coverage: ## Run tests with coverage
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

bench: ## Run benchmarks
	go test -bench=. -benchmem ./...

clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.out coverage.html
	go clean

fmt: ## Format code
	go fmt ./...

vet: ## Run go vet
	go vet ./...

lint: fmt vet ## Run linters (fmt + vet)

worker-pool: build ## Run worker pool demo
	./bin/goroutine-3000 worker-pool

fan-out-fan-in: build ## Run fan-out/fan-in demo
	./bin/goroutine-3000 fan-out-fan-in

pipeline: build ## Run pipeline demo
	./bin/goroutine-3000 pipeline

stress-test: build ## Run stress test demo
	./bin/goroutine-3000 stress-test

context: build ## Run context demo
	./bin/goroutine-3000 context

error-handling: build ## Run error handling demo
	./bin/goroutine-3000 error-handling

all-demos: run ## Run all demonstrations

install: ## Install binary to $GOPATH/bin
	go install .
