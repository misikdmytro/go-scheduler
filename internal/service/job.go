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
}

type jobService struct {
	r repository.WorkerRepository
	b broker.Broker[model.JobLaunchMessage]
}

func NewJobService(r repository.WorkerRepository, b broker.Broker[model.JobLaunchMessage]) JobService {
	return &jobService{r: r, b: b}
}

func (s *jobService) LaunchJob(c context.Context, workerID string, input map[string]any) (string, error) {
	w, err := s.r.Get(c, workerID)
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
