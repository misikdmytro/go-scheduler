package service_test

import (
	"context"

	"github.com/misikdmytro/go-job-runner/internal/broker"
	"github.com/misikdmytro/go-job-runner/internal/model"
	"github.com/misikdmytro/go-job-runner/internal/repository"
	publicmodel "github.com/misikdmytro/go-job-runner/pkg/model"

	"github.com/stretchr/testify/mock"
)

type repoMock struct {
	mock.Mock
}

func (m *repoMock) Get(c context.Context, id string) (*model.Worker, error) {
	args := m.Called(c, id)
	return args.Get(0).(*model.Worker), args.Error(1)
}

func (m *repoMock) Create(c context.Context, name, topic string) (string, error) {
	args := m.Called(c, name, topic)
	return args.String(0), args.Error(1)
}

func (m *repoMock) Delete(c context.Context, id string) (bool, error) {
	args := m.Called(c, id)
	return args.Bool(0), args.Error(1)
}

var _ repository.WorkerRepository = (*repoMock)(nil)

type brokerMock[T any] struct {
	mock.Mock
}

func (m *brokerMock[T]) Publish(c context.Context, topic string, data T) error {
	args := m.Called(c, topic, data)
	return args.Error(0)
}

var _ broker.Broker[publicmodel.JobLaunchMessage] = (*brokerMock[publicmodel.JobLaunchMessage])(nil)
