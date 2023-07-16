package consumer

import (
	"context"
	"log"
	"time"

	"github.com/misikdmytro/go-job-runner/internal/config"
	"github.com/misikdmytro/go-job-runner/internal/service"
	"github.com/misikdmytro/go-job-runner/pkg/model"
)

type jobStatusConsumer struct {
	*rabbitMQConsumer[model.JobEventMessage]
	j service.JobService
}

func NewJobStatusConsumer(
	rc config.RabbitMQConfig,
	cc config.ConsumerConfig,
	j service.JobService,
) Consumer {
	return &jobStatusConsumer{
		newRabbitMQConsumer[model.JobEventMessage](rc, cc),
		j,
	}
}

func (c *jobStatusConsumer) Setup() error {
	return c.setup()
}

func (c *jobStatusConsumer) Consume() error {
	return c.consume(c.consumeCallback)
}

func (c *jobStatusConsumer) Shutdown(ctx context.Context) error {
	return c.shutdown()
}

func (c *jobStatusConsumer) consumeCallback(ctx context.Context, m model.JobEventMessage, err error) error {
	if err != nil {
		log.Printf("consuming job statuses error: %v", err)
		return nil
	}

	_, err = c.j.AppendJobStatus(
		ctx,
		m.JobID,
		m.Message,
		time.Unix(m.Timestamp, 0),
		m.Output,
	)

	return err
}
