package dependency

import (
	"github.com/gin-gonic/gin"
	"github.com/misikdmytro/go-job-runner/internal/broker"
	"github.com/misikdmytro/go-job-runner/internal/config"
	"github.com/misikdmytro/go-job-runner/internal/handler"
	"github.com/misikdmytro/go-job-runner/internal/repository"
	"github.com/misikdmytro/go-job-runner/internal/server"
	"github.com/misikdmytro/go-job-runner/internal/service"
	"github.com/misikdmytro/go-job-runner/pkg/model"
)

type Dependency struct {
	CFG config.Config

	WR repository.WorkerRepository
	WS service.WorkerService
	WH handler.WorkerHandler

	JB broker.Broker[model.JobLaunchMessage]
	JS service.JobService
	JH handler.JobHandler

	HR repository.HealthRepository
	HS service.HealthService
	HH handler.HealthHandler

	E *gin.Engine
}

func NewDependency() (*Dependency, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	wr := repository.NewWorkerRepository(cfg.DB)
	ws := service.NewWorkerService(wr)
	wh := handler.NewWorkerHandler(ws)

	jb := broker.NewJobLaunchBroker(cfg.RabbitMQ, cfg.Jobs)
	js := service.NewJobService(wr, jb)
	jh := handler.NewJobHandler(js)

	hr := repository.NewHealthRepository(cfg.DB)
	hs := service.NewHealthService(hr)
	hh := handler.NewHealthHandler(hs)

	e := server.NewEngine(wh, jh, hh)

	return &Dependency{
		CFG: cfg,

		WR: wr,
		WS: ws,
		WH: wh,

		JB: jb,
		JS: js,
		JH: jh,

		E: e,
	}, nil
}
