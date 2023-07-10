package e2e_test

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateWorkerShouldReturnCreated(t *testing.T) {
	c := client{baseAddress: "http://localhost:4001"}
	r, err := c.CreateWorker(
		fmt.Sprintf("test-%s", uuid.NewString()),
		fmt.Sprintf("test-%s", uuid.NewString()),
	)

	if err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, r.ID)
}
