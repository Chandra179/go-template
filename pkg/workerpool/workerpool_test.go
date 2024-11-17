package workerpool

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestWorkerPool_New(t *testing.T) {
	tests := []struct {
		name       string
		numWorkers int
		bufferSize int
	}{
		{"valid configuration", 5, 10},
		{"zero buffer size", 3, 0},
		{"minimum workers", 1, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wp := New(tt.numWorkers, tt.bufferSize)
			if wp == nil {
				t.Fatal("expected non-nil WorkerPool")
			}
			if wp.numWorkers != tt.numWorkers {
				t.Errorf("expected %d workers, got %d", tt.numWorkers, wp.numWorkers)
			}
			if cap(wp.tasks) != tt.bufferSize {
				t.Errorf("expected task buffer size %d, got %d", tt.bufferSize, cap(wp.tasks))
			}
		})
	}
}

func TestWorkerPool_SuccessfulTasks(t *testing.T) {
	wp := New(3, 10)
	ctx := context.Background()

	// Start the worker pool
	wp.Start(ctx)

	// Create successful tasks
	numTasks := 5
	taskResults := make(chan int, numTasks)

	for i := 0; i < numTasks; i++ {
		taskID := i
		wp.Submit(Task{
			ID: taskID,
			Func: func(ctx context.Context) error {
				taskResults <- taskID
				return nil
			},
		})
	}

	// Stop the pool and wait for completion
	wp.Stop()

	// Check results
	close(taskResults)
	completedTasks := make(map[int]bool)
	for id := range taskResults {
		completedTasks[id] = true
	}

	if len(completedTasks) != numTasks {
		t.Errorf("expected %d completed tasks, got %d", numTasks, len(completedTasks))
	}

	// Verify no errors were reported
	errors := wp.CollectErrors()
	if len(errors) != 0 {
		t.Errorf("expected no errors, got %d errors", len(errors))
	}
}

func TestWorkerPool_ErrorHandling(t *testing.T) {
	wp := New(2, 5)
	ctx := context.Background()

	wp.Start(ctx)

	expectedError := errors.New("task error")
	numTasks := 3

	// Submit tasks that will fail
	for i := 0; i < numTasks; i++ {
		wp.Submit(Task{
			ID: i,
			Func: func(ctx context.Context) error {
				return expectedError
			},
		})
	}

	wp.Stop()
	errors := wp.CollectErrors()

	if len(errors) != numTasks {
		t.Errorf("expected %d errors, got %d", numTasks, len(errors))
	}

	for _, err := range errors {
		if !strings.Contains(err.Error(), expectedError.Error()) {
			t.Errorf("expected error containing %q, got %q", expectedError, err)
		}
	}
}

func TestWorkerPool_ContextCancellation(t *testing.T) {
	wp := New(2, 5)
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	wp.Start(ctx)

	// Submit a task that will block
	wp.Submit(Task{
		ID: 1,
		Func: func(ctx context.Context) error {
			select {
			case <-time.After(100 * time.Millisecond):
				return nil
			case <-ctx.Done():
				return ctx.Err()
			}
		},
	})

	wp.Stop()
	errors := wp.CollectErrors()

	if len(errors) != 1 {
		t.Errorf("expected 1 error due to context cancellation, got %d", len(errors))
	}

	if len(errors) > 0 && !strings.Contains(errors[0].Error(), "context deadline exceeded") {
		t.Errorf("expected deadline exceeded error, got %v", errors[0])
	}
}

func TestWorkerPool_ConcurrentTasks(t *testing.T) {
	wp := New(4, 10)
	ctx := context.Background()

	wp.Start(ctx)

	numTasks := 20
	completedTasks := make(chan int, numTasks)

	// Submit tasks that simulate work with random delays
	for i := 0; i < numTasks; i++ {
		taskID := i
		wp.Submit(Task{
			ID: taskID,
			Func: func(ctx context.Context) error {
				time.Sleep(time.Duration(taskID%3) * time.Millisecond)
				completedTasks <- taskID
				return nil
			},
		})
	}

	wp.Stop()
	close(completedTasks)

	completed := make(map[int]bool)
	for id := range completedTasks {
		completed[id] = true
	}

	if len(completed) != numTasks {
		t.Errorf("expected %d completed tasks, got %d", numTasks, len(completed))
	}

	// Verify all tasks were processed
	for i := 0; i < numTasks; i++ {
		if !completed[i] {
			t.Errorf("task %d was not completed", i)
		}
	}
}
