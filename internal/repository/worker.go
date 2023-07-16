package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/misikdmytro/go-job-runner/internal/config"
	"github.com/misikdmytro/go-job-runner/internal/model"

	_ "github.com/lib/pq"
)

type WorkerRepository interface {
	Get(context.Context, string) (*model.Worker, error)
	Create(context.Context, string, string) (string, error)
	Delete(context.Context, string) (bool, error)
}

type workerRepository struct {
	c config.DBConfig
}

func NewWorkerRepository(c config.DBConfig) WorkerRepository {
	return &workerRepository{c: c}
}

func (r *workerRepository) Get(ctx context.Context, id string) (*model.Worker, error) {
	db, err := sqlx.Connect("postgres", BuildDataSourceName(r.c))
	if err != nil {
		return nil, err
	}

	var w model.Worker
	if err = db.GetContext(ctx, &w, "SELECT * FROM get_worker($1)", id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &w, nil
}

func (r *workerRepository) Create(ctx context.Context, name, description string) (string, error) {
	db, err := sqlx.Connect("postgres", BuildDataSourceName(r.c))
	if err != nil {
		return "", err
	}

	var id string
	err = db.GetContext(ctx, &id, "SELECT create_worker ($1, $2)", name, description)
	return id, err
}

func (r *workerRepository) Delete(ctx context.Context, id string) (bool, error) {
	db, err := sqlx.Connect("postgres", BuildDataSourceName(r.c))
	if err != nil {
		return false, err
	}

	var deletedID *string
	err = db.GetContext(ctx, &deletedID, "SELECT delete_worker($1)", id)
	return deletedID != nil && *deletedID == id, err
}
