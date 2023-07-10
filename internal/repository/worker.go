package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/misikdmytro/go-job-runner/internal/config"

	_ "github.com/lib/pq"
)

type WorkerRepository interface {
	Create(context.Context, string, string) (string, error)
}

type workerRepository struct {
	c config.DBConfig
}

func NewWorkerRepository(c config.DBConfig) WorkerRepository {
	return &workerRepository{c: c}
}

func (r *workerRepository) Create(ctx context.Context, name, description string) (string, error) {
	db, err := sqlx.Connect("postgres", BuildDataSourceName(r.c))
	if err != nil {
		return "", err
	}

	id := uuid.NewString()
	_, err = db.ExecContext(ctx, "CALL create_worker($1, $2, $3)", id, name, description)
	if err != nil {
		return "", err
	}

	return id, nil
}
