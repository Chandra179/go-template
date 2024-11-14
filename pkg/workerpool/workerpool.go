package workerpool

import (
	"context"
	"sync"
)

type TaskFunc func(ctx context.Context) (interface{}, error)

func worker(ctx context.Context, taskCh <-chan TaskFunc, resultCh chan<- interface{}, errCh chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()

	for task := range taskCh {
		select {
		case <-ctx.Done():
			return
		default:
			result, err := task(ctx)
			if err != nil {
				select {
				case errCh <- err:
				default:
				}
				return
			}
			select {
			case resultCh <- result:
			case <-ctx.Done():
				return
			}
		}
	}
}

func workerPool(ctx context.Context, numWorkers int, tasks []TaskFunc) (chan interface{}, error) {
	taskCh := make(chan TaskFunc)
	resultCh := make(chan interface{}, len(tasks)) // Buffered channel with capacity equal to the number of tasks
	errCh := make(chan error, 1)

	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(ctx, taskCh, resultCh, errCh, &wg)
	}

	go func() {
		defer close(taskCh)
		for _, task := range tasks {
			select {
			case taskCh <- task:
			case <-ctx.Done():
				return
			}
		}
	}()

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	select {
	case err := <-errCh:
		return resultCh, err
	case <-ctx.Done():
		return resultCh, nil
	}
}
