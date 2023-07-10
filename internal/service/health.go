package service

import (
	"github.com/misikdmytro/go-job-runner/internal/exception"
	"github.com/misikdmytro/go-job-runner/internal/repository"
)

type HealthService interface {
	Health() error
}

type healthService struct {
	r repository.HealthRepository
}

func NewHealthService(r repository.HealthRepository) HealthService {
	return &healthService{r: r}
}

func (s *healthService) Health() error {
	err := s.r.Health()
	if err != nil {
		return exception.JobError{
			Code:    exception.UnhealthService,
			Message: "failed to check health",
			Err:     err,
		}
	}

	return nil
}
