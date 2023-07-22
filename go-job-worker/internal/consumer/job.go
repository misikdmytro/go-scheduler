package consumer

import (
	"context"
	"log"
	"time"

	"github.com/misikdmytro/go-job-worker/internal/broker"
	"github.com/misikdmytro/go-job-worker/internal/config"
	"github.com/misikdmytro/go-job-worker/internal/helper"
	"github.com/misikdmytro/go-job-worker/internal/model"
	"github.com/misikdmytro/go-job-worker/internal/service"
	pkgmodel "github.com/misikdmytro/go-job-worker/pkg/model"
)

type jobLaunchConsumer struct {
	*rabbitMQConsumer[pkgmodel.JobLaunchMessage]
	j service.WorkerService
	b broker.Broker[pkgmodel.JobEventMessage]
}

func NewJobLaunchConsumer(
	rc config.RabbitMQConfig,
	cc config.SubscriberConfig,
	j service.WorkerService,
	b broker.Broker[pkgmodel.JobEventMessage],
) Consumer {
	return &jobLaunchConsumer{
		rabbitMQConsumer: newRabbitMQConsumer[pkgmodel.JobLaunchMessage](rc, cc),
		j:                j,
		b:                b,
	}
}

func (j *jobLaunchConsumer) Consume() error {
	return j.consume(j.consumeCallback)
}

func (j *jobLaunchConsumer) Shutdown(ctxt context.Context) error {
	return j.shutdown()
}

func (j *jobLaunchConsumer) consumeCallback(ctxt context.Context, msg pkgmodel.JobLaunchMessage, err error) error {
	if err != nil {
		log.Printf("job launch message error: %v", err)
		j.b.Publish(
			ctxt,
			"event",
			pkgmodel.JobEventMessage{
				JobID:     msg.JobID,
				Message:   "random worker error",
				Timestamp: time.Now().UnixMilli(),
				Output: map[string]any{
					"error": err.Error(),
				},
			},
		)

		return nil
	}

	err = j.b.Publish(
		ctxt,
		"event",
		pkgmodel.JobEventMessage{
			JobID:     msg.JobID,
			Message:   "random worker started",
			Timestamp: time.Now().UnixMilli(),
			Output:    msg.Input,
		},
	)
	if err != nil {
		return err
	}

	args, err := helper.To[model.JobArgs](msg.Input)
	if err != nil {
		return err
	}

	err = j.j.Launch(ctxt, args)
	if err != nil {
		return err
	}

	err = j.b.Publish(
		ctxt,
		"event",
		pkgmodel.JobEventMessage{
			JobID:     msg.JobID,
			Message:   "random worker finished",
			Timestamp: time.Now().UnixMilli(),
		},
	)
	if err != nil {
		return err
	}

	return nil
}
