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

func (r *workerRepository) Create(ctx context.Context, name, topic string) (string, error) {
	db, err := sqlx.Connect("postgres", buildDataSourceName(r.c))
	if err != nil {
		return "", err
	}

	id := uuid.NewString()
	if err := db.GetContext(ctx, &id, "CALL create_worker($1, $2, $3)", id, name, topic); err != nil {
		return "", nil
	}

	return id, nil
}
