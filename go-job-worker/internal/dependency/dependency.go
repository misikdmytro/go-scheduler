package dependency

import "github.com/misikdmytro/go-job-worker/internal/config"

type Dependency struct {
	CFG config.Config
}

func NewDependency() (*Dependency, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	return &Dependency{
		CFG: cfg,
	}, nil
}
