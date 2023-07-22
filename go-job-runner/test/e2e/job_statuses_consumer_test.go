package e2e_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/misikdmytro/go-job-runner/internal/broker"
	"github.com/misikdmytro/go-job-runner/internal/config"
	"github.com/misikdmytro/go-job-runner/internal/repository"
	"github.com/misikdmytro/go-job-runner/pkg/model"
	"github.com/sethvargo/go-retry"

	amqp "github.com/rabbitmq/amqp091-go"
)

func TestConsumerJobStatusesShouldSaveResultToDB(t *testing.T) {
	cfg, err := config.LoadConfigFrom("../..")
	if err != nil {
		t.Fatal(err)
	}

	ch, close, err := broker.NewRabbitMQChannel(cfg.RabbitMQ)
	if err != nil {
		t.Fatal(err)
	}
	defer close()

	event := model.JobEventMessage{
		JobID:     uuid.NewString(),
		Message:   fmt.Sprintf("test-%s", uuid.NewString()),
		Timestamp: time.Now().UnixMilli(),
		Output: map[string]any{
			"test": "test",
		},
	}

	bytes, err := json.Marshal(event)
	if err != nil {
		t.Fatal(err)
	}

	err = ch.PublishWithContext(
		context.Background(),
		"job-events",
		"event",
		true,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        bytes,
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	db, err := sqlx.Connect("postgres", repository.BuildDataSourceName(cfg.DB))
	if err != nil {
		t.Fatal(err)
	}

	err = retry.Do(
		context.Background(),
		retry.WithMaxRetries(
			5,
			retry.NewConstant(1*time.Second),
		),
		func(ctxt context.Context) error {
			var count int
			err = db.GetContext(ctxt, &count, "SELECT COUNT(*) FROM jobs WHERE job_id = $1 AND message = $2", event.JobID, event.Message)
			if err != nil {
				return err
			}

			if count != 1 {
				return retry.RetryableError(fmt.Errorf("count is not 1"))
			}

			return nil
		},
	)
	if err != nil {
		t.Fatal(err)
	}
}
