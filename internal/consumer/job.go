package consumer

import (
	"context"
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

func (c *jobStatusConsumer) Consume(ctxt context.Context) error {
	return c.consume(ctxt, c.consumeCallback)
}

func (c *jobStatusConsumer) Err() chan error {
	return c.err
}

func (c *jobStatusConsumer) Close() error {
	return c.close()
}

func (c *jobStatusConsumer) consumeCallback(ctx context.Context, m model.JobEventMessage) error {
	_, err := c.j.AppendJobStatus(
		ctx,
		m.JobID,
		m.Message,
		time.Unix(m.Timestamp, 0),
		m.Output,
	)

	return err
}
