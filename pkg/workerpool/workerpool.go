package workerpool

import (
	"fmt"
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
}
