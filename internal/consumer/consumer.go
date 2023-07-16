package consumer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/misikdmytro/go-job-runner/internal/broker"
	"github.com/misikdmytro/go-job-runner/internal/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

var ErrConsumerClosed = fmt.Errorf("consumer closed")

type Consumer interface {
	Setup() error
	Consume() error
	Shutdown(context.Context) error
}

type consumeCallback[T any] func(context.Context, T, error) error

type rabbitMQConsumer[T any] struct {
	rc     config.RabbitMQConfig
	cc     config.ConsumerConfig
	close  func() error
	cancel func()
	b      *amqp.Channel
}

func newRabbitMQConsumer[T any](
	rc config.RabbitMQConfig,
	cc config.ConsumerConfig,
) *rabbitMQConsumer[T] {
	return &rabbitMQConsumer[T]{
		rc:     rc,
		cc:     cc,
		cancel: func() {},
		close:  func() error { return nil },
	}
}

func (r *rabbitMQConsumer[T]) setup() error {
	b, close, err := broker.NewRabbitMQChannel(r.rc)
	if err != nil {
		return err
	}

	r.b = b
	r.close = close

	err = b.ExchangeDeclare(
		r.cc.Exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	q, err := b.QueueDeclare(
		r.cc.Queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	err = b.QueueBind(
		q.Name,
		r.cc.RoutingKey,
		r.cc.Exchange,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *rabbitMQConsumer[T]) consume(callback consumeCallback[T]) error {
	msgs, err := r.b.Consume(
		r.cc.Queue,
		r.cc.Consumer,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	ctxt, c := context.WithCancel(context.Background())
	defer c()

	r.cancel = c

	for {
		select {
		case <-ctxt.Done():
			return ErrConsumerClosed
		case msg := <-msgs:
			var m T
			if err := json.Unmarshal(msg.Body, &m); err != nil {
				msg.Nack(false, false)
				callback(ctxt, m, err)
				continue
			}

			if err := callback(ctxt, m, nil); err != nil {
				msg.Nack(false, true)
				callback(ctxt, m, err)
				continue
			}

			if err := msg.Ack(false); err != nil {
				msg.Nack(false, false)
				callback(ctxt, m, err)
			}
		}
	}
}

func (r *rabbitMQConsumer[T]) shutdown() error {
	r.cancel()
	return r.close()
}
