package producer

import (
	"fmt"
	"context"
	"day12/model"
)

func Producer(ctx context.Context, numJobs int, jobs chan<- model.Job){
	defer close(jobs)

	for i := 1; i <= numJobs; i++ {
		job := model.Job{ID: i}

		select {
		case <-ctx.Done():
			fmt.Println("Producer: context cancelled, early stopping")
			return
		case jobs <- job:
			fmt.Printf("Producer: created job #%d\n", job.ID)
		}
	}

	fmt.Println("Producer: all jobs created.")
}