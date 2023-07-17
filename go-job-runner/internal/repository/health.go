package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/misikdmytro/go-job-runner/internal/config"
)

type HealthRepository interface {
	Health(context.Context) error
}

type healthRepository struct {
	c config.DBConfig
}

func NewHealthRepository(c config.DBConfig) HealthRepository {
	return &healthRepository{c: c}
}

func (r *healthRepository) Health(c context.Context) error {
	db, err := sqlx.Connect("postgres", BuildDataSourceName(r.c))
	if err != nil {
		return err
	}

	return db.PingContext(c)
}
