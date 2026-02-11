package sink

import (
	"context"
	"day12/cache"
	"day12/model"
	"fmt"
	"sync"
)

func Collect(ctx context.Context, processResult <- chan model.ProcessResult, c *cache.ResultCache, wg *sync.WaitGroup ) {
	defer wg.Done()

	for{
		select { 
		case <- ctx.Done():
			fmt.Println("Sink: Context canclled.....Stopping!")
			return
		case result, ok := <-processResult:
			if !ok {
				fmt.Println("Sink: All results Collected!")
				return
			}
			c.Store(result)
			fmt.Printf("Sink: stored result for job #%d\n", result.JobID)
		}
	
	}
}