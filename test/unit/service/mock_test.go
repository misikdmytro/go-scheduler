package service_test

import (
	"context"
	"time"

	"github.com/misikdmytro/go-job-runner/internal/broker"
	"github.com/misikdmytro/go-job-runner/internal/model"
	"github.com/misikdmytro/go-job-runner/internal/repository"
	publicmodel "github.com/misikdmytro/go-job-runner/pkg/model"

	"github.com/stretchr/testify/mock"
)

type workerRepoMock struct {
	mock.Mock
}

func (m *workerRepoMock) Get(c context.Context, id string) (*model.Worker, error) {
	args := m.Called(c, id)
	return args.Get(0).(*model.Worker), args.Error(1)
}

func (m *workerRepoMock) Create(c context.Context, name, topic string) (string, error) {
	args := m.Called(c, name, topic)
	return args.String(0), args.Error(1)
}

func (m *workerRepoMock) Delete(c context.Context, id string) (bool, error) {
	args := m.Called(c, id)
	return args.Bool(0), args.Error(1)
}

var _ repository.WorkerRepository = (*workerRepoMock)(nil)

type brokerMock[T any] struct {
	mock.Mock
}

func (m *brokerMock[T]) Publish(c context.Context, topic string, data T) error {
	args := m.Called(c, topic, data)
	return args.Error(0)
}

var _ broker.Broker[publicmodel.JobLaunchMessage] = (*brokerMock[publicmodel.JobLaunchMessage])(nil)

type jobRepoMock struct {
	mock.Mock
}

func (r *jobRepoMock) AppendStatus(c context.Context, jobID, message string, timestamp time.Time, output map[string]any) (int64, error) {
	args := r.Called(c, jobID, message, timestamp, output)
	return args.Get(0).(int64), args.Error(1)
}

var _ repository.JobRepository = (*jobRepoMock)(nil)
