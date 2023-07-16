package repository_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/misikdmytro/go-job-runner/internal/config"
	"github.com/misikdmytro/go-job-runner/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestAppendStatus(t *testing.T) {
	cfg, err := config.LoadConfigFrom("../../..")
	if err != nil {
		t.Fatal(err)
	}

	r := repository.NewJobRepository(cfg.DB)

	jobID := uuid.NewString()
	message := fmt.Sprintf("test-%s", uuid.NewString())
	timestamp := time.Now()
	output := map[string]any{
		"test": "test",
	}

	id, err := r.AppendStatus(context.Background(), jobID, message, timestamp, output)
	assert.NoError(t, err)
	assert.NotEmpty(t, id)

	db, err := sqlx.Connect("postgres", repository.BuildDataSourceName(cfg.DB))
	if err != nil {
		t.Fatal(err)
	}

	var count int
	err = db.Get(&count, "SELECT COUNT(*) FROM jobs WHERE job_id = $1 AND message = $2", jobID, message)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 1, count)
}
