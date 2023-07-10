package service

import (
	"context"
	"fmt"

	"github.com/misikdmytro/go-job-runner/internal/exception"
	"github.com/misikdmytro/go-job-runner/internal/model"
	"github.com/misikdmytro/go-job-runner/internal/repository"
)

type WorkerService interface {
	Get(context.Context, string) (*model.Worker, error)
	Create(context.Context, string, string) (string, error)
	Delete(context.Context, string) (bool, error)
}

type workerService struct {
	r repository.WorkerRepository
}

func NewWorkerService(r repository.WorkerRepository) WorkerService {
	return &workerService{r: r}
}

func (s *workerService) Get(c context.Context, id string) (*model.Worker, error) {
	w, err := s.r.Get(c, id)
	if err != nil {
		return nil, exception.JobError{
			Code:    exception.UnknownError,
			Message: fmt.Sprintf("failed to get worker with id %s", id),
			Err:     err,
		}
	}

	return w, nil
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

func (s *workerService) Delete(c context.Context, id string) (bool, error) {
	ok, err := s.r.Delete(c, id)
	if err != nil {
		return false, exception.JobError{
			Code:    exception.UnknownError,
			Message: fmt.Sprintf("failed to delete worker with id %s", id),
			Err:     err,
		}
	}

	return ok, nil
}
