package broker

import (
	"context"
	"encoding/json"

	"github.com/misikdmytro/go-job-worker/internal/config"
	"github.com/misikdmytro/go-job-worker/pkg/model"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Broker[T any] interface {
	Publish(context.Context, string, T) error
}

type jobEventsBroker struct {
	rc config.RabbitMQConfig
	jc config.PublisherConfig
}

func NewJobEventsBroker(rc config.RabbitMQConfig, jc config.PublisherConfig) Broker[model.JobEventMessage] {
	return &jobEventsBroker{
		rc: rc,
		jc: jc,
	}
}

func (b *jobEventsBroker) Publish(ctxt context.Context, routingKey string, msg model.JobEventMessage) error {
	ch, close, err := NewRabbitMQChannel(b.rc)
	if err != nil {
		return err
	}
	defer close()

	bytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return ch.PublishWithContext(
		ctxt,
		b.jc.Exchange,
		routingKey,
		true,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        bytes,
		},
	)
}
