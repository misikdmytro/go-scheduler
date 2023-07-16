package consumer

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/misikdmytro/go-job-runner/internal/broker"
	"github.com/misikdmytro/go-job-runner/internal/config"
)

type Consumer interface {
	Consume(context.Context) error
	Err() chan error
	Close() error
}

type consumeCallback[T any] func(context.Context, T) error

type rabbitMQConsumer[T any] struct {
	rc     config.RabbitMQConfig
	cc     config.ConsumerConfig
	err    chan error
	cancel func()
}

func newRabbitMQConsumer[T any](
	rc config.RabbitMQConfig,
	cc config.ConsumerConfig,
) *rabbitMQConsumer[T] {
	return &rabbitMQConsumer[T]{
		rc:     rc,
		cc:     cc,
		err:    make(chan error),
		cancel: func() {},
	}
}

func (r *rabbitMQConsumer[T]) consume(ctx context.Context, callback consumeCallback[T]) error {
	b, close, err := broker.BuildRabbitMQChannel(r.rc)
	if err != nil {
		return err
	}
	defer close()

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

	msgs, err := b.Consume(
		q.Name,
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

	cancel, c := context.WithCancel(ctx)
	r.cancel = c

	defer c()
	for {
		select {
		case <-cancel.Done():
			r.err <- cancel.Err()
		case msg := <-msgs:
			var m T
			if err := json.Unmarshal(msg.Body, &m); err != nil {
				if nErr := msg.Nack(false, false); nErr != nil {
					r.err <- errors.Join(nErr, err)
				}

				r.err <- err
			}

			if err := callback(cancel, m); err != nil {
				if nErr := msg.Nack(false, true); nErr != nil {
					r.err <- errors.Join(nErr, err)
				}

				r.err <- err
			}

			if err := msg.Ack(false); err != nil {
				if nErr := msg.Nack(false, false); nErr != nil {
					r.err <- errors.Join(nErr, err)
				}

				r.err <- err
			}
		}
	}
}

func (r *rabbitMQConsumer[T]) close() error {
	r.cancel()
	return nil
}
