package processor

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"
	"day12/model"
)

func StartWorkers(ctx context.Context, numWorkers int, fetchResults <-chan model.FetchResult, processResults chan<- model.ProcessResult) {
	var wg sync.WaitGroup

	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(ctx, i, fetchResults, processResults, &wg)
	}

	go func() { // 
		wg.Wait()
		close(processResults)
		fmt.Println("Processor: all workers done.")
	}()
}

// worker for single processor goroutine 

func worker(ctx context.Context, id int, fetchResults <-chan model.FetchResult, processResults chan<- model.ProcessResult, wg *sync.WaitGroup) {
	defer wg.Done()

	for fetchResult := range fetchResults {
		select {
		case <-ctx.Done(): // For Graceful shutdown before processing
			fmt.Printf("Processor worker %d: context cancelled, stopping.\n", id)
			return
		default:
		}

		sleepTime := time.Duration(30+rand.Intn(70)) * time.Millisecond
		fmt.Printf("Processor worker %d: processing job #%d (will take %v)\n", id, fetchResult.JobID, sleepTime)
		time.Sleep(sleepTime)

		result := model.ProcessResult{
			JobID:  fetchResult.JobID,
			Output: strings.ToUpper(fetchResult.Data), // Processed Data will be all Upeer case
		}

		select {
		case <-ctx.Done(): // Graceful shutdown while sending data downstream
			fmt.Printf("Processor worker %d: context cancelled while sending result.\n", id)
			return
		case processResults <- result:

		}
	}
}
