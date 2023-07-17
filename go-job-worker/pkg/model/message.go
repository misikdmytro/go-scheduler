package model

type JobLaunchMessage struct {
	JobID     string         `json:"job_id"`
	Timestamp int64          `json:"timestamp"`
	Input     map[string]any `json:"input"`
}

type JobEventMessage struct {
	JobID     string         `json:"job_id"`
	Message   string         `json:"message"`
	Timestamp int64          `json:"timestamp"`
	Output    map[string]any `json:"output"`
}
