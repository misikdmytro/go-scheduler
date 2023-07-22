package dependency

import (
	"github.com/misikdmytro/go-job-worker/internal/broker"
	"github.com/misikdmytro/go-job-worker/internal/config"
	"github.com/misikdmytro/go-job-worker/internal/consumer"
	"github.com/misikdmytro/go-job-worker/internal/service"
	"github.com/misikdmytro/go-job-worker/pkg/model"
)

type Dependency struct {
	CFG config.Config

	WS service.WorkerService

	JB broker.Broker[model.JobEventMessage]
	JC consumer.Consumer
}

func NewDependency() (*Dependency, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	ws := service.NewWorkerService()

	jb := broker.NewJobEventsBroker(cfg.RabbitMQ, cfg.JobEvents)
	jc := consumer.NewJobLaunchConsumer(cfg.RabbitMQ, cfg.Jobs, ws, jb)

	return &Dependency{
		CFG: cfg,

		WS: ws,

		JB: jb,
		JC: jc,
	}, nil
}
