package service_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/misikdmytro/go-job-runner/internal/model"
	"github.com/misikdmytro/go-job-runner/internal/service"
	publicmodel "github.com/misikdmytro/go-job-runner/pkg/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLaunchJob(t *testing.T) {
	data := []struct {
		testName   string
		worker     *model.Worker
		workerErr  error
		publishErr error
		fails      bool
	}{
		{
			testName:  "worker error",
			workerErr: fmt.Errorf("test error"),
			fails:     true,
		},
		{
			testName: "worker nil",
			worker:   nil,
			fails:    true,
		},
		{
			testName:   "publish error",
			worker:     &model.Worker{Name: uuid.NewString()},
			publishErr: fmt.Errorf("test error"),
			fails:      true,
		},
		{
			testName: "ok",
			worker:   &model.Worker{Name: uuid.NewString()},
			fails:    false,
		},
	}

	for _, d := range data {
		t.Run(d.testName, func(t *testing.T) {
			workerID := uuid.NewString()

			wr := &workerRepoMock{}
			wr.On("Get", mock.Anything, workerID).Return(d.worker, d.workerErr)

			b := &brokerMock[publicmodel.JobLaunchMessage]{}
			if d.worker != nil {
				b.On("Publish", mock.Anything, fmt.Sprintf("worker.%s", d.worker.Name), mock.Anything).Return(d.publishErr)
			}

			jr := &jobRepoMock{}

			s := service.NewJobService(wr, jr, b)
			jobID, err := s.LaunchJob(context.Background(), workerID, nil)

			if d.fails {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, jobID)
			}
		})
	}
}

func TestAppendJobStatus(t *testing.T) {
	data := []struct {
		testName string
		jobID    string
		message  string
		output   map[string]any
		fails    bool
	}{
		{
			testName: "ok",
			jobID:    uuid.NewString(),
			message:  "test message",
			output:   map[string]any{"test": "test"},
			fails:    false,
		},
	}

	for _, d := range data {
		t.Run(d.testName, func(t *testing.T) {
			jr := &jobRepoMock{}
			jr.On("AppendStatus", mock.Anything, d.jobID, mock.Anything, mock.Anything, mock.Anything).Return(int64(1), nil)

			s := service.NewJobService(nil, jr, nil)
			_, err := s.AppendJobStatus(context.Background(), d.jobID, d.message, time.Now(), d.output)

			if d.fails {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				jr.AssertCalled(t, "AppendStatus", mock.Anything, d.jobID, d.message, mock.Anything, d.output)
			}
		})
	}
}
