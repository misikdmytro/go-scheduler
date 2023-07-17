package e2e_test

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetShouldBeOK(t *testing.T) {
	c := newClient()

	name := fmt.Sprintf("test-%s", uuid.NewString())
	desc := fmt.Sprintf("test-%s", uuid.NewString())

	r, err := c.CreateWorker(
		name,
		desc,
	)

	if err != nil {
		t.Fatal(err)
	}

	res, err := c.GetWorker(r.ID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, r.ID, res.Worker.ID)
	assert.Equal(t, name, res.Worker.Name)
	assert.Equal(t, desc, res.Worker.Description)
}

func TestGetShouldReturnNotFound(t *testing.T) {
	c := newClient()

	_, err := c.GetWorker(uuid.NewString())

	assert.Error(t, err)
	assert.Equal(t, "unexpected status code: 404", err.Error())
}
