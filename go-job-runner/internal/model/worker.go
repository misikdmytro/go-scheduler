package model

import "time"

type Worker struct {
	ID          string
	Name        string
	Description string
}

type JobStatus struct {
	ID        int64
	JobID     string
	Message   string
	Timestamp time.Time
	Output    map[string]any
}
