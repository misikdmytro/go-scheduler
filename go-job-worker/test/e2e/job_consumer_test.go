package e2e_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/misikdmytro/go-job-worker/internal/broker"
	"github.com/misikdmytro/go-job-worker/internal/config"
	"github.com/misikdmytro/go-job-worker/pkg/model"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestWorkerShouldTrackJobProgress(t *testing.T) {
	cfg, err := config.LoadConfigFrom("../..")
	if err != nil {
		t.Fatal(err)
	}

	ch, close, err := broker.NewRabbitMQChannel(cfg.RabbitMQ)
	if err != nil {
		t.Fatal(err)
	}
	defer close()

	q, err := ch.QueueDeclare(
		fmt.Sprintf("test-%s", uuid.New().String()),
		false,
		true,
		true,
		false,
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	err = ch.QueueBind(
		q.Name,
		"event",
		cfg.JobEvents.Exchange,
		false,
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	l := model.JobLaunchMessage{
		JobID:     uuid.New().String(),
		Timestamp: time.Now().UnixMilli(),
		Input: map[string]any{
			"key":   "value",
			"value": "key",
		},
	}
	b, err := json.Marshal(l)
	if err != nil {
		t.Fatal(err)
	}
	ch.PublishWithContext(
		context.Background(),
		"jobs",
		"worker.random",
		true,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        b,
		},
	)

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		true,
		false,
		false,
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

outer:
	for {
		select {
		case <-c.Done():
			t.Fatal("timeout")
		case msg := <-msgs:
			var event model.JobEventMessage
			if err := json.Unmarshal(msg.Body, &event); err != nil {
				t.Fatal(err)
			}

			if event.JobID == l.JobID {
				assert.GreaterOrEqual(t, event.Timestamp, l.Timestamp)
				break outer
			}
		}
	}
}
