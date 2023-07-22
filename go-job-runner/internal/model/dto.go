package model

import "time"

type JobStatusDTO struct {
	ID        int64
	JobID     string `db:"job_id"`
	Message   string
	CreatedAt time.Time `db:"created_at"`
	Output    string
}
