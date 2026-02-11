# Worker Pipeline

A concurrent **Worker Pipeline** built in Go that demonstrates goroutines, channels, `sync.WaitGroup`, `context` cancellation, and thread-safe caching — all wired together in a multi-stage data pipeline.

## Architecture

```
┌──────────┐    ┌──────────┐    ┌───────────┐    ┌──────┐    ┌───────────┐
│ Producer │───>│ Fetcher  │───>│ Processor │───>│ Sink │───>│   Cache   │
│ (1 goroutine) │ (N workers) │  (N workers)│ (1 goroutine) │ (thread-safe) │
└──────────┘    └──────────┘    └───────────┘    └──────┘    └───────────┘
     jobs ch       fetchResults ch   processResults ch          in-memory map
```

### Pipeline Stages

| Stage | Package | Description |
|-------|---------|-------------|
| **Producer** | `producer/` | Generates `Job` structs (with incremental IDs) and sends them into the `jobs` channel. |
| **Fetcher** | `fetcher/` | Multiple concurrent workers read jobs, simulate network I/O (sleep), and emit `FetchResult` values. |
| **Processor** | `processor/` | Multiple concurrent workers read fetch results, simulate CPU-bound work, transform data (uppercase), and emit `ProcessResult` values. |
| **Sink** | `sink/` | A single goroutine that reads processed results and stores them in a thread-safe cache. |
| **Cache** | `cache/` | Thread-safe in-memory store (`sync.RWMutex` + `map`) for final results. |

## Data Models

Defined in `model/model.go`:

```go
type Job struct {
    ID  int
    URL string
}

type FetchResult struct {
    JobID int
    Data  string
}

type ProcessResult struct {
    JobID  int
    Output string
}
```

## Key Concepts Demonstrated

- **Fan-out / Fan-in** — Multiple fetcher and processor workers run concurrently, reading from a shared input channel and writing to a shared output channel.
- **Buffered Channels** — Each stage is connected via buffered channels to decouple producers from consumers.
- **Graceful Shutdown via `context.Context`** — Every stage listens for `ctx.Done()` and exits cleanly on timeout or cancellation.
- **`sync.WaitGroup`** — Used to wait for all workers in a stage to finish before closing the downstream channel.
- **`defer close(ch)`** — Each stage closes its output channel when done, signaling downstream stages to stop.
- **Thread-safe Cache** — Uses `sync.RWMutex` to safely read/write results from concurrent goroutines.

## Configuration

In `main.go`, you can tune these values:

```go
numJobs       := 10               // Total jobs to produce
numFetchers   := 3                // Concurrent fetcher workers
numProcessors := 3                // Concurrent processor workers
timeout       := 10 * time.Second // Pipeline-wide timeout
```

## Getting Started

### Prerequisites

- **Go 1.21+** installed

### Run

```bash
cd Worker-Pipline
go run main.go
```

### Sample Output

```
Producer: created job #1
Producer: created job #2
...
Fetcher worker 1: fetching job #1 (will take 50ms)
Fetcher worker 2: fetching job #2 (will take 50ms)
...
Processor worker 1: processing job #1 (will take 67ms)
...
Sink: stored result for job #1
Sink: stored result for job #2
...
Sink: All results Collected!

------- Pipeline complete ------
Total results in cache: 10
Results:
 Job: 1 -> RAW-DATA-FOR-JOB-1
 Job: 2 -> RAW-DATA-FOR-JOB-2
...
```

## Project Structure

```
Worker-Pipline/
├── main.go            # Entry point — wires all stages together
├── go.mod             # Go module definition
├── model/
│   └── model.go       # Data types: Job, FetchResult, ProcessResult
├── producer/
│   └── producer.go    # Stage 1: Job generation
├── fetcher/
│   └── fetcher.go     # Stage 2: Simulated network I/O
├── processor/
│   └── processor.go   # Stage 3: Simulated CPU processing
├── sink/
│   └── sink.go        # Stage 4: Collects results into cache
└── cache/
    └── cache.go       # Thread-safe in-memory result store
```

## License

This project is for educational/learning purposes as part of the **GoLang-101** series.
