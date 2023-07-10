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
	jc config.JobsConfig
}

func NewJobLaunchBroker(rc config.RabbitMQConfig, jc config.JobsConfig) Broker[model.JobLaunchMessage] {
	return &jobLaunchBroker{rc: rc, jc: jc}
}

func (b *jobLaunchBroker) Publish(c context.Context, key string, job model.JobLaunchMessage) error {
	conn, err := amqp.Dial(BuildRabbitMQURL(b.rc))
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		b.jc.Exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	bytes, err := json.Marshal(job)
	if err != nil {
		return err
	}

	return ch.Publish(
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
