package main

import (
	"context"
	"day12/cache"
	"day12/fetcher"
	"day12/model"
	"day12/processor"
	"day12/producer"
	"day12/sink"
	"fmt"
	"sync"
	"time"
)

func main() {

	// Configuration 
	numJobs := 10               
	numFetchers := 3            
	numProcessors := 3          
	timeout := 10 * time.Second 

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Here i am Creating a buffered channel for each pipline stage 
	jobs := make(chan model.Job, numJobs)
	fetchResults := make(chan model.FetchResult, numJobs)
	processResults := make(chan model.ProcessResult, numJobs)


	resultCache := cache.NewResultCache()

	// Starting the pipeline stages 
	go producer.Producer(ctx, numJobs, jobs) // Stage 1 -> Producer

	fetcher.StartWorkers(ctx,numFetchers,jobs,fetchResults) // stage 2 -> fetcher -> StartWorkers spwans multiple fetchers concurrently

	processor.StartWorkers(ctx,numProcessors,fetchResults,processResults) // Stage 3 -> Processor -> same as fetcher but simulates CPU I/O work

	var sinkWg sync.WaitGroup
	sinkWg.Add(1) // Stage 4: Sink reads process results and stores them in the cache
	go sink.Collect(ctx,processResults,resultCache, &sinkWg)
	sinkWg.Wait()

	fmt.Println("\n------- Pipeline complete ------")
	fmt.Printf("Total results in cache: %d\n", resultCache.Size())
	fmt.Println("Results:")
	resultCache.PrintAll()
	

}