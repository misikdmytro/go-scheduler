package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/misikdmytro/go-job-runner/internal/config"
)

type JobRepository interface {
	AppendStatus(context.Context, string, string, time.Time, map[string]any) (int64, error)
}

type jobRepository struct {
	c config.DBConfig
}

func NewJobRepository(c config.DBConfig) JobRepository {
	return &jobRepository{c: c}
}

func (r *jobRepository) AppendStatus(ctx context.Context, jobID, message string, timestamp time.Time, output map[string]any) (int64, error) {
	db, err := sqlx.Connect("postgres", BuildDataSourceName(r.c))
	if err != nil {
		return 0, err
	}

	b, err := json.Marshal(output)
	if err != nil {
		return 0, err
	}

	var id int64
	err = db.GetContext(ctx, &id, "SELECT append_job_status ($1, $2, $3, $4)", jobID, message, timestamp, b)
	return id, err
}
