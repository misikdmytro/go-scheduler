package e2e_test

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateWorkerShouldReturnCreated(t *testing.T) {
	c := newClient()
	r, err := c.CreateWorker(
		fmt.Sprintf("test-%s", uuid.NewString()),
		fmt.Sprintf("test-%s", uuid.NewString()),
	)

	if err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, r.ID)
}

func TestCreateWorkerShouldReturnBadRequest(t *testing.T) {
	d := []struct {
		testName string
		name     string
		desc     string
	}{
		{
			testName: "empty name",
			name:     "",
			desc:     "desc",
		},
		{
			testName: "empty desc",
			name:     "name",
			desc:     "",
		},
		{
			testName: "empty name and desc",
			name:     "",
			desc:     "",
		},
		{
			testName: "too long name",
			name:     "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890",
			desc:     "desc",
		},
		{
			testName: "too long desc",
			name:     "name",
			desc:     "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890",
		},
	}

	for _, d := range d {
		t.Run(d.testName, func(t *testing.T) {
			c := newClient()
			_, err := c.CreateWorker(d.name, d.desc)

			assert.Error(t, err)
			assert.Equal(t, err.Error(), "unexpected status code: 400")
		})
	}
}
