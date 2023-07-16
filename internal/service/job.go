package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/misikdmytro/go-job-runner/internal/broker"
	"github.com/misikdmytro/go-job-runner/internal/exception"
	"github.com/misikdmytro/go-job-runner/internal/repository"
	"github.com/misikdmytro/go-job-runner/pkg/model"
)

type JobService interface {
	LaunchJob(context.Context, string, map[string]any) (string, error)
	AppendJobStatus(context.Context, string, string, time.Time, map[string]any) (int64, error)
}

type jobService struct {
	wr repository.WorkerRepository
	jr repository.JobRepository
	b  broker.Broker[model.JobLaunchMessage]
}

func NewJobService(
	wr repository.WorkerRepository,
	jr repository.JobRepository,
	b broker.Broker[model.JobLaunchMessage],
) JobService {
	return &jobService{wr: wr, jr: jr, b: b}
}

func (s *jobService) LaunchJob(c context.Context, workerID string, input map[string]any) (string, error) {
	w, err := s.wr.Get(c, workerID)
	if err != nil {
		return "", err
	}

	if w == nil {
		return "", exception.JobError{
			Code:    exception.WorkerNotFound,
			Message: fmt.Sprintf("worker '%s' not found", workerID),
		}
	}

	jobID := uuid.NewString()
	err = s.b.Publish(
		c,
		fmt.Sprintf("worker.%s", w.Name),
		model.JobLaunchMessage{
			JobID:     jobID,
			Timestamp: time.Now().UnixMilli(),
			Input:     input,
		},
	)

	return jobID, err
}

func (s *jobService) AppendJobStatus(c context.Context, jobID, message string, timestamp time.Time, output map[string]any) (int64, error) {
	id, err := s.jr.AppendStatus(c, jobID, message, timestamp, output)
	if err != nil {
		return 0, exception.JobError{
			Code:    exception.FailedToAppendJobStatus,
			Err:     err,
			Message: fmt.Sprintf("failed to append job status for job '%s'", jobID),
		}
	}

	return id, nil
}
