package workerpool

import (
	"fmt"
<<<<<<< HEAD
	"log"
	"sync"
	"time"
)

// WorkerPool represents the worker pool, with a task channel size that can be buffered or unbuffered.
type WorkerPool struct {
	taskCh   chan func() error
	resultCh chan error
	wg       sync.WaitGroup
}

// NewWorkerPool creates a new worker pool with specified number of workers and an optional buffer size for the task channel.
func NewWorkerPool(numWorkers int, bufferSize int) *WorkerPool {
	// Initialize the channel with the provided buffer size, 0 means unbuffered
	taskCh := make(chan func() error, bufferSize)
	pool := &WorkerPool{
		taskCh:   taskCh,
		resultCh: make(chan error),
	}

	// Start workers
	for i := 0; i < numWorkers; i++ {
		pool.wg.Add(1)
		go pool.worker(i)
	}
	return pool
}

// AddTask adds a task to the worker pool. If the task channel is closed, it logs an error.
func (p *WorkerPool) AddTask(task func() error) {
	select {
	case p.taskCh <- task:
		// Successfully added task
	default:
		log.Println("Warning: Task channel is closed or full, could not add task.")
	}
}

// Results returns the result channel for processing errors or results.
func (p *WorkerPool) Results() <-chan error {
	return p.resultCh
}

// Close stops accepting new tasks, waits for workers to complete, and closes resultCh
func (p *WorkerPool) Close(timeout time.Duration) error {
	close(p.taskCh)
	done := make(chan struct{})
	go func() {
		p.wg.Wait()
		close(done)
	}()
	select {
	case <-done:
		close(p.resultCh)
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("worker pool shutdown timed out after %v", timeout)
	}
}

// worker processes tasks and sends results to resultCh
func (p *WorkerPool) worker(workerID int) {
	defer p.wg.Done()
	for task := range p.taskCh {
		if task == nil {
			log.Printf("Worker %d received a nil task, skipping\n", workerID)
			continue
		}
		err := task()
		select {
		case p.resultCh <- err:
		default:
			// If resultCh is unbuffered and no listener, skip
			log.Printf("Worker %d processed task with error: %v\n", workerID, err)
		}
	}
	log.Printf("Worker %d exiting\n", workerID)
=======
	"sync"
)

// Task represents a single unit of work, which is a function to be executed.
type Task struct {
	ID   int
	Func func() error // Function to be executed by the worker
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
func (wp *WorkerPool) Start() {
	for i := 1; i <= wp.numWorkers; i++ {
		wp.wg.Add(1)
		go func(workerID int) {
			defer wp.wg.Done()
			for task := range wp.tasks {
				if err := task.Func(); err != nil {
					// Send the error to the error channel
					wp.errCh <- fmt.Errorf("worker %d failed to process task %d: %v", workerID, task.ID, err)
				} else {
					fmt.Printf("Worker %d successfully processed task %d\n", workerID, task.ID)
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
>>>>>>> 4055eda (update worker pool)
}
