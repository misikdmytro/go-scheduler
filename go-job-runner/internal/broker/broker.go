package broker

import (
	"context"
	"encoding/json"

	"github.com/misikdmytro/go-job-runner/internal/config"
	"github.com/misikdmytro/go-job-runner/pkg/model"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Broker[T any] interface {
	Publish(context.Context, string, T) error
}

type jobLaunchBroker struct {
	rc config.RabbitMQConfig
	jc config.PublisherConfig
}

func NewJobLaunchBroker(rc config.RabbitMQConfig, jc config.PublisherConfig) Broker[model.JobLaunchMessage] {
	return &jobLaunchBroker{rc: rc, jc: jc}
}

func (b *jobLaunchBroker) Publish(ctxt context.Context, key string, job model.JobLaunchMessage) error {
	ch, close, err := NewRabbitMQChannel(b.rc)
	if err != nil {
		return err
	}
	defer close()

	bytes, err := json.Marshal(job)
	if err != nil {
		return err
	}

	return ch.PublishWithContext(
		ctxt,
		b.jc.Exchange,
		key,
		true,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        bytes,
		},
	)
}
