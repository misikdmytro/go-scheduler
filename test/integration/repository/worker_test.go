package repository_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/misikdmytro/go-job-runner/internal/config"
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
