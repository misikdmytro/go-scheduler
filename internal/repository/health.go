package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/misikdmytro/go-job-runner/internal/config"
)

type HealthRepository interface {
	Health() error
}

type healthRepository struct {
	c config.DBConfig
}

func NewHealthRepository(c config.DBConfig) HealthRepository {
	return &healthRepository{c: c}
}

func (r *healthRepository) Health() error {
	db, err := sqlx.Connect("postgres", buildDataSourceName(r.c))
	if err != nil {
		return err
	}

	return db.Ping()
}
