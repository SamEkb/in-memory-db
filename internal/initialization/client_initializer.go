package initialization

import (
	"go.uber.org/zap"
	"in-memory-db/internal/configuration"
)

type ClientInitializer struct {
	Logger *zap.Logger
	Config *configuration.NetworkConfig
}

func InitializeClient() (*ClientInitializer, error) {
	conf, err := configuration.NewConfiguration()
	if err != nil {
		return &ClientInitializer{}, err
	}

	logger, err := InitializeLogger(conf.Logging)

	initializer := &ClientInitializer{
		Logger: logger,
		Config: conf.Network,
	}

	return initializer, nil
}
