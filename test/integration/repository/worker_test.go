package repository_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/misikdmytro/go-job-runner/internal/config"
	"github.com/misikdmytro/go-job-runner/internal/model"
	"github.com/misikdmytro/go-job-runner/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestCreateShouldDoIt(t *testing.T) {
	cfg, err := config.LoadConfigFrom("../../..")
	if err != nil {
		t.Fatal(err)
	}

	r := repository.NewWorkerRepository(cfg.DB)
	name := fmt.Sprintf("test-%s", uuid.NewString())
	desc := fmt.Sprintf("test-%s", uuid.NewString())

	id, err := r.Create(context.Background(), name, desc)
	assert.NoError(t, err)
	assert.NotEmpty(t, id)

	db, err := sqlx.Connect("postgres", repository.BuildDataSourceName(cfg.DB))
	if err != nil {
		t.Fatal(err)
	}

	var count int
	err = db.Get(&count, "SELECT COUNT(*) FROM workers WHERE id = $1 AND name = $2 AND description = $3", id, name, desc)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 1, count)
}

func TestGetShouldReturnEmptyModel(t *testing.T) {
	cfg, err := config.LoadConfigFrom("../../..")
	if err != nil {
		t.Fatal(err)
	}

	r := repository.NewWorkerRepository(cfg.DB)
	id := uuid.NewString()

	w, err := r.Get(context.Background(), id)
	assert.NoError(t, err)
	assert.Nil(t, w)
}

func TestGetShouldReturnModel(t *testing.T) {
	cfg, err := config.LoadConfigFrom("../../..")
	if err != nil {
		t.Fatal(err)
	}

	db, err := sqlx.Connect("postgres", repository.BuildDataSourceName(cfg.DB))
	if err != nil {
		t.Fatal(err)
	}

	id := uuid.NewString()
	name := fmt.Sprintf("test-%s", uuid.NewString())
	desc := fmt.Sprintf("test-%s", uuid.NewString())

	_, err = db.Exec("INSERT INTO workers (id, name, description) VALUES ($1, $2, $3)", id, name, desc)
	if err != nil {
		t.Fatal(err)
	}

	r := repository.NewWorkerRepository(cfg.DB)

	w, err := r.Get(context.Background(), id)
	assert.NoError(t, err)
	assert.Equal(t, model.Worker{
		ID:          id,
		Name:        name,
		Description: desc,
	}, *w)
}

func TestDeleteShouldReturnOk(t *testing.T) {
	cfg, err := config.LoadConfigFrom("../../..")
	if err != nil {
		t.Fatal(err)
	}

	db, err := sqlx.Connect("postgres", repository.BuildDataSourceName(cfg.DB))
	if err != nil {
		t.Fatal(err)
	}

	id := uuid.NewString()
	name := fmt.Sprintf("test-%s", uuid.NewString())
	desc := fmt.Sprintf("test-%s", uuid.NewString())

	_, err = db.Exec("INSERT INTO workers (id, name, description) VALUES ($1, $2, $3)", id, name, desc)
	if err != nil {
		t.Fatal(err)
	}

	r := repository.NewWorkerRepository(cfg.DB)

	ok, err := r.Delete(context.Background(), id)
	assert.NoError(t, err)
	assert.True(t, ok)
}

func TestDeleteShouldReturnFalse(t *testing.T) {
	cfg, err := config.LoadConfigFrom("../../..")
	if err != nil {
		t.Fatal(err)
	}

	r := repository.NewWorkerRepository(cfg.DB)

	ok, err := r.Delete(context.Background(), uuid.NewString())
	assert.NoError(t, err)
	assert.False(t, ok)
}
