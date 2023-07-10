package service

import (
	"context"

	"github.com/misikdmytro/go-job-runner/internal/exception"
	"github.com/misikdmytro/go-job-runner/internal/repository"
)

type WorkerService interface {
	Create(context.Context, string, string) (string, error)
}

type workerService struct {
	r repository.WorkerRepository
}

func NewWorkerService(r repository.WorkerRepository) WorkerService {
	return &workerService{r: r}
}

func (s *workerService) Create(c context.Context, name, description string) (string, error) {
	id, err := s.r.Create(c, name, description)
	if err != nil {
		return "", exception.JobError{
			Code:    exception.UnknownError,
			Message: "failed to create worker",
			Err:     err,
		}
	}

	return id, nil
}
