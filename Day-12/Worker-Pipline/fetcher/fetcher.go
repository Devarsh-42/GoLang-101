package fetcher

import (
	"context"
	"day12/model"
	"fmt"
	"sync"
	"time"
)

// Optional Step (Used in Newtowrk intensive Applications)
// Each worker reads jobs from the jobs channel  simulates I/O work & and sends the result to the results channel

func StartWorkers(ctx context.Context, numWorkers int, jobs <-chan model.Job, results chan<- model.FetchResult) {
	var wg sync.WaitGroup

	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(ctx, i, jobs, results, &wg)
	}

	go func() {
		wg.Wait()
		close(results)
		fmt.Println("Fetcher: all workers done.")
	}()

}

func worker(ctx context.Context, id int, jobs <-chan model.Job, results chan<- model.FetchResult, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		select {
		case <-ctx.Done():
			fmt.Printf("Fetcher worker %d: context cancelled, stopping.\n", id)
			return
		default:
		}

		sleepTime := time.Duration(50 * time.Millisecond)
		fmt.Printf("Fetcher worker %d: fetching job #%d (will take %v)\n", id, job.ID, sleepTime)
		time.Sleep(sleepTime)

		result := model.FetchResult{
			JobID: job.ID,
			Data:  fmt.Sprintf("raw-data-for-job-%d", job.ID),
		}

		select {
		case <-ctx.Done():
			fmt.Printf("Fetcher worker %d: context cancelled while sending result.\n", id)
			return
		case results <- result:
		}
	}

}
