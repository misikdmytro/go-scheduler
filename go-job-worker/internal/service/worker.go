package service

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/misikdmytro/go-job-worker/internal/model"
)

type WorkerService interface {
	Launch(context.Context, model.JobArgs) error
}

type workerService struct {
}

func NewWorkerService() WorkerService {
	return &workerService{}
}

func (w *workerService) Launch(ctx context.Context, args model.JobArgs) error {
	log.Printf("job args: %v", args)

	// create random
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	if r.Float64() < 0.5 {
		return nil
	} else {
		return fmt.Errorf("random fail")
	}
}
