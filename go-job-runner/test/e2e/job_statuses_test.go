package e2e_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestJobStatusesEmpty(t *testing.T) {
	c := newClient()
	r, err := c.GetJobStatuses(uuid.NewString())

	if err != nil {
		t.Fatal(err)
	}

	assert.Empty(t, r)
}
