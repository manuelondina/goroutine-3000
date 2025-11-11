package patterns

import (
	"context"
	"fmt"
	"time"
)

// DemoContext demonstrates context-based cancellation and timeout handling
func DemoContext() {
	fmt.Println("=== Context-Based Cancellation Demo ===")

	// Demo 1: Context with timeout
	fmt.Println("\n1. Context with timeout:")
	ctx1, cancel1 := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel1()

	done1 := make(chan bool)
	go longRunningTask("Task-1", ctx1, done1)
	<-done1

	// Demo 2: Context with cancellation
	fmt.Println("\n2. Context with manual cancellation:")
	ctx2, cancel2 := context.WithCancel(context.Background())

	done2 := make(chan bool)
	go longRunningTask("Task-2", ctx2, done2)

	// Cancel after 1 second
	time.Sleep(1 * time.Second)
	fmt.Println("  Cancelling Task-2...")
	cancel2()
	<-done2

	// Demo 3: Context with deadline
	fmt.Println("\n3. Context with deadline:")
	deadline := time.Now().Add(1500 * time.Millisecond)
	ctx3, cancel3 := context.WithDeadline(context.Background(), deadline)
	defer cancel3()

	done3 := make(chan bool)
	go longRunningTask("Task-3", ctx3, done3)
	<-done3

	fmt.Println("\nAll context demonstrations completed!")
}

func longRunningTask(name string, ctx context.Context, done chan<- bool) {
	defer func() { done <- true }()

	ticker := time.NewTicker(300 * time.Millisecond)
	defer ticker.Stop()

	for i := 1; i <= 10; i++ {
		select {
		case <-ctx.Done():
			fmt.Printf("  %s: Cancelled at iteration %d (reason: %v)\n", name, i, ctx.Err())
			return
		case <-ticker.C:
			fmt.Printf("  %s: Working... iteration %d\n", name, i)
		}
	}
	fmt.Printf("  %s: Completed all iterations\n", name)
}
