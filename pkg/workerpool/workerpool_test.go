package workerpool

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// TestBasicWorkerPool tests if the worker pool processes tasks correctly.
func TestBasicWorkerPool(t *testing.T) {
	numWorkers := 3
	bufferSize := 10
	wp := New(numWorkers, bufferSize)

	// Create a task that does nothing and always succeeds.
	taskID := 1
	task := Task{
		ID: taskID,
		Func: func(ctx context.Context) error {
			// Simulate some work
			time.Sleep(100 * time.Millisecond)
			return nil
		},
	}

	// Start the worker pool
	wp.Start(context.Background())

	// Submit the task
	wp.Submit(task)

	// Wait for the worker pool to process the task
	wp.Stop()

	// Collect any errors (if any)
	errors := wp.CollectErrors()
	if len(errors) > 0 {
		t.Fatalf("Expected no errors, but got: %v", errors)
	}
}

// TestWorkerPoolErrorHandling tests if the worker pool correctly reports errors.
func TestWorkerPoolErrorHandling(t *testing.T) {
	numWorkers := 3
	bufferSize := 10
	wp := New(numWorkers, bufferSize)

	// Create a task that always fails.
	taskID := 2
	task := Task{
		ID: taskID,
		Func: func(ctx context.Context) error {
			// Simulate a failure
			return fmt.Errorf("task %d failed", taskID)
		},
	}

	// Start the worker pool
	wp.Start(context.Background())

	// Submit the task
	wp.Submit(task)

	// Wait for the worker pool to process the task
	wp.Stop()

	// Collect any errors (should contain the error from the task)
	errors := wp.CollectErrors()
	if len(errors) == 0 {
		t.Fatal("Expected an error, but got none")
	}

	// Just check that the error is not nil
	if errors[0] == nil {
		t.Fatal("Expected a non-nil error, but got nil")
	}
}

// TestWorkerPoolContextCancellation tests if workers stop processing when the context is canceled.
func TestWorkerPoolContextCancellation(t *testing.T) {
	numWorkers := 2
	bufferSize := 5
	wp := New(numWorkers, bufferSize)

	// Create a task that simulates a long-running operation
	task := Task{
		ID: 3,
		Func: func(ctx context.Context) error {
			select {
			case <-time.After(500 * time.Millisecond):
				return nil
			case <-ctx.Done():
				// Simulate cancellation
				return fmt.Errorf("task %d was cancelled", 3)
			}
		},
	}

	// Create a context with a timeout that will cancel the task
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// Start the worker pool
	wp.Start(ctx)

	// Submit the task
	wp.Submit(task)

	// Wait for the worker pool to process the task (it should be cancelled due to timeout)
	wp.Stop()

	// Collect any errors (should contain the cancellation error)
	errors := wp.CollectErrors()
	if len(errors) == 0 {
		t.Fatal("Expected an error, but got none")
	}

	// Just check that the error is not nil
	if errors[0] == nil {
		t.Fatal("Expected a non-nil error, but got nil")
	}
}

// TestWorkerPoolShutdown tests if the pool shuts down correctly.
func TestWorkerPoolShutdown(t *testing.T) {
	numWorkers := 2
	bufferSize := 5
	wp := New(numWorkers, bufferSize)

	// Create a task that simulates some work
	task := Task{
		ID: 4,
		Func: func(ctx context.Context) error {
			time.Sleep(50 * time.Millisecond)
			return nil
		},
	}

	// Start the worker pool
	wp.Start(context.Background())

	// Submit multiple tasks
	for i := 0; i < 5; i++ {
		task.ID = i
		wp.Submit(task)
	}

	// Stop the worker pool (this should wait for all tasks to finish)
	wp.Stop()

	// Collect any errors (there should be none)
	errors := wp.CollectErrors()
	if len(errors) > 0 {
		t.Fatalf("Expected no errors, but got: %v", errors)
	}
}
