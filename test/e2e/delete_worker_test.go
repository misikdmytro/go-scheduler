package e2e_test

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDeleteShouldReturnDeleted(t *testing.T) {
	c := newClient()

	r, err := c.CreateWorker(
		fmt.Sprintf("test-%s", uuid.NewString()),
		fmt.Sprintf("test-%s", uuid.NewString()),
	)

	if err != nil {
		t.Fatal(err)
	}

	res, err := c.DeleteWorker(r.ID)

	assert.NoError(t, err)
	assert.True(t, res.Deleted)
}

func TestDeleteShouldReturnNotDeleted(t *testing.T) {
	c := newClient()

	res, err := c.DeleteWorker(uuid.NewString())

	assert.NoError(t, err)
	assert.False(t, res.Deleted)
}
