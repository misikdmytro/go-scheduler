package e2e_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/misikdmytro/go-job-runner/internal/broker"
	"github.com/misikdmytro/go-job-runner/internal/config"
	"github.com/misikdmytro/go-job-runner/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestLaunchJobShouldDoIt(t *testing.T) {
	c := newClient()
	wn, wd := fmt.Sprintf("test-%s", uuid.NewString()), fmt.Sprintf("test-%s", uuid.NewString())
	r, err := c.CreateWorker(
		wn,
		wd,
	)

	if err != nil {
		t.Fatal(err)
	}
	assert.NotEmpty(t, r.ID)

	cfg, err := config.LoadConfigFrom("../..")
	if err != nil {
		t.Fatal(err)
	}

	b, close, err := broker.NewRabbitMQChannel(cfg.RabbitMQ)
	if err != nil {
		t.Fatal(err)
	}
	defer close()

	q, err := b.QueueDeclare(
		fmt.Sprintf("test-%s", uuid.NewString()),
		false,
		true,
		true,
		false,
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	err = b.QueueBind(
		q.Name,
		fmt.Sprintf("worker.%s", wn),
		cfg.Jobs.Exchange,
		false,
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	d, err := b.Consume(
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

	input := map[string]any{
		"test": "test",
	}

	lr, err := c.LaunchJob(r.ID, input)
	if err != nil {
		t.Fatal(err)
	}

	timeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	select {
	case msg := <-d:
		var m model.JobLaunchMessage
		err := json.Unmarshal(msg.Body, &m)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, lr.JobID, m.JobID)
		assert.Equal(t, input, m.Input)
		assert.GreaterOrEqual(t, 5*time.Second, time.Since(time.UnixMilli(m.Timestamp)))
	case <-timeout.Done():
		t.Fatal(timeout.Err())
	}
}

func TestLaunchJobShouldReturnBadRequest(t *testing.T) {
	data := []struct {
		testName string
		workerID string
		input    map[string]any
	}{
		{
			testName: "worker not found",
			workerID: uuid.NewString(),
			input: map[string]any{
				"test": "test",
			},
		},
		{
			testName: "worker id is empty",
			workerID: "",
			input: map[string]any{
				"test": "test",
			},
		},
		{
			testName: "input is nil",
			workerID: uuid.NewString(),
			input:    nil,
		},
	}

	for _, d := range data {
		t.Run(d.testName, func(t *testing.T) {
			c := newClient()
			_, err := c.LaunchJob(d.workerID, d.input)
			assert.Error(t, err)
			assert.Equal(t, "unexpected status code: 400", err.Error())
		})
	}
}
