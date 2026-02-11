package cache

import (
	"day12/model"
	"fmt"
	"sync"
)

type ResultCache struct {
	mu sync.RWMutex
	results map[int]model.ProcessResult
}

func NewResultCache() *ResultCache{
	return &ResultCache{
		results: make(map[int]model.ProcessResult),
	}
}

func (r *ResultCache) Store(result model.ProcessResult){
	r.mu.Lock()
	defer r.mu.Unlock()
	r.results[result.JobID] = result
} 


func (r *ResultCache) GET(jobID int) (model.ProcessResult, bool){
	r.mu.RLock()
	defer r.mu.RLock()
	result, ok := r.results[jobID]
	return result, ok
}

func (r *ResultCache) Size() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.results)
}

func (r *ResultCache) PrintAll(){
	r.mu.RLock()
	defer r.mu.RUnlock()
	for id, result := range r.results{
		fmt.Printf(" Job: %d -> %s\n",id,result.Output)
	}
}