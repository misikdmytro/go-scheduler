package consumer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/misikdmytro/go-job-runner/internal/broker"
	"github.com/misikdmytro/go-job-runner/internal/config"
)

var ErrConsumerClosed = fmt.Errorf("consumer closed")

type Consumer interface {
	Consume() error
	Shutdown(context.Context) error
}

type consumeCallback[T any] func(context.Context, T, error) error

type rabbitMQConsumer[T any] struct {
	rc     config.RabbitMQConfig
	cc     config.ConsumerConfig
	close  func() error
	cancel func()
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

func (r *rabbitMQConsumer[T]) consume(callback consumeCallback[T]) error {
	b, close, err := broker.NewRabbitMQChannel(r.rc)
	if err != nil {
		return err
	}
	defer close()

	msgs, err := b.Consume(
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
