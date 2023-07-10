package model

type JobStatus int

const (
	JobStarted JobStatus = iota
	JobFinished
	JobFailed
	JobCanceled
)

type JobLaunchMessage struct {
	JobID string         `json:"job_id"`
	Input map[string]any `json:"input"`
}

type JobEventMessage struct {
	JobID     string         `json:"job_id"`
	Status    JobStatus      `json:"status"`
	Timestamp int64          `json:"timestamp"`
	Output    map[string]any `json:"output"`
}
