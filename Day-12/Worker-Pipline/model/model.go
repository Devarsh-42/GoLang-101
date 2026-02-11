package model

type Job struct {
	ID int
	URL string
}

type FetchResult struct {
	JobID int
	Data string
}


type ProcessResult struct {
	JobID int
	Output string
}