package workerpool

import (
	"context"
	"fmt"
	"sync"
)

// Task represents a single unit of work, which is a function to be executed.
type Task struct {
	ID   int
	Func func(context.Context) error // Function to be executed by the worker, accepting context
}

// WorkerPool manages a pool of workers to process tasks.
type WorkerPool struct {
	numWorkers int
	tasks      chan Task
	errCh      chan error // Error channel to report errors from workers
	wg         sync.WaitGroup
}

// New creates and initializes a new WorkerPool.
func New(numWorkers, bufferSize int) *WorkerPool {
	return &WorkerPool{
		numWorkers: numWorkers,
		tasks:      make(chan Task, bufferSize),
		errCh:      make(chan error, bufferSize), // Buffered channel for errors
	}
}

// Start initializes the worker goroutines.
func (wp *WorkerPool) Start(ctx context.Context) {
	for i := 1; i <= wp.numWorkers; i++ {
		wp.wg.Add(1)
		go func(workerID int) {
			defer wp.wg.Done()
			for {
				select {
				case task, ok := <-wp.tasks:
					if !ok {
						// Channel closed and all tasks have been processed
						return
					}
					// Execute the task with context, handling possible cancellation or timeout
					if err := task.Func(ctx); err != nil {
						// Send the error to the error channel
						wp.errCh <- fmt.Errorf("worker %d failed to process task %d: %v", workerID, task.ID, err)
					} else {
						fmt.Printf("Worker %d successfully processed task %d\n", workerID, task.ID)
					}
				case <-ctx.Done():
					// Handle the cancellation or timeout of the context
					fmt.Printf("Worker %d was cancelled due to context timeout\n", workerID)
					return
				}
			}
		}(i)
	}
}

// Submit adds a task to the pool for processing.
func (wp *WorkerPool) Submit(task Task) {
	wp.tasks <- task
}

// Stop signals the workers to stop processing and waits for them to finish.
func (wp *WorkerPool) Stop() {
	close(wp.tasks) // Close task channel to signal workers to stop
	wp.wg.Wait()    // Wait for all workers to finish

	// Close the error channel after all workers are done
	close(wp.errCh)
}

// CollectErrors listens for errors from the error channel.
func (wp *WorkerPool) CollectErrors() []error {
	var errs []error
	for err := range wp.errCh {
		errs = append(errs, err)
	}
	return errs
}
