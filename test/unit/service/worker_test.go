package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/misikdmytro/go-job-runner/internal/repository"
	"github.com/misikdmytro/go-job-runner/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type repoMock struct {
	mock.Mock
}

func (m *repoMock) Create(c context.Context, name, topic string) (string, error) {
	args := m.Called(c, name, topic)
	return args.String(0), args.Error(1)
}

var _ repository.WorkerRepository = (*repoMock)(nil)

func TestCreateShouldReturnError(t *testing.T) {
	r := &repoMock{}
	r.On("Create", mock.Anything, mock.Anything, mock.Anything).Return("", fmt.Errorf("test error"))
	s := service.NewWorkerService(r)

	_, err := s.Create(context.Background(), "name", "desc")
	assert.Error(t, err)
}

func TestCreateShouldDoIt(t *testing.T) {
	r := &repoMock{}
	id := uuid.NewString()
	r.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(id, nil)
	s := service.NewWorkerService(r)

	res, err := s.Create(context.Background(), "name", "desc")

	assert.NoError(t, err)
	assert.Equal(t, id, res)
}
