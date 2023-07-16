package e2e_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/misikdmytro/go-job-runner/internal/broker"
	"github.com/misikdmytro/go-job-runner/internal/config"
	"github.com/misikdmytro/go-job-runner/internal/consumer"
	"github.com/misikdmytro/go-job-runner/internal/repository"
	"github.com/misikdmytro/go-job-runner/internal/service"
	"github.com/misikdmytro/go-job-runner/pkg/model"
	"github.com/sethvargo/go-retry"

	amqp "github.com/rabbitmq/amqp091-go"
)

func TestConsumerJobStatusesShouldSaveResultToDB(t *testing.T) {
	cfg, err := config.LoadConfigFrom("../..")
	if err != nil {
		t.Fatal(err)
	}

	ctxt, cancel := context.WithCancel(context.Background())
	defer cancel()

	wr := repository.NewWorkerRepository(cfg.DB)
	jr := repository.NewJobRepository(cfg.DB)
	b := broker.NewJobLaunchBroker(cfg.RabbitMQ, cfg.Jobs)

	js := service.NewJobService(wr, jr, b)

	c := consumer.NewJobStatusConsumer(cfg.RabbitMQ, cfg.JobEventsConsumer, js)
	defer c.Shutdown(ctxt)

	err = c.Setup()
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		if err := c.Consume(); !errors.Is(err, consumer.ErrConsumerClosed) {
			panic(err)
		}
	}()

	ch, close, err := broker.NewRabbitMQChannel(cfg.RabbitMQ)
	if err != nil {
		t.Fatal(err)
	}
	defer close()

	event := model.JobEventMessage{
		JobID:     uuid.NewString(),
		Message:   fmt.Sprintf("test-%s", uuid.NewString()),
		Timestamp: time.Now().Unix(),
		Output: map[string]any{
			"test": "test",
		},
	}

	bytes, err := json.Marshal(event)
	if err != nil {
		t.Fatal(err)
	}

	err = ch.PublishWithContext(
		ctxt,
		cfg.JobEventsConsumer.Exchange,
		cfg.JobEventsConsumer.RoutingKey,
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
		ctxt,
		retry.WithMaxRetries(
			10,
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
