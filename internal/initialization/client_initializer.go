package initialization

import (
	"go.uber.org/zap"
	"in-memory-db/internal/configuration"
)

type ClientInitializer struct {
	Logger *zap.Logger
	Config *configuration.Configuration
}

func InitializeClient() (*ClientInitializer, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		panic("Unable to initialize logger")
	}
	config, err := configuration.NewConfiguration()
	if err != nil {
		return &ClientInitializer{}, err
	}

	initializer := &ClientInitializer{
		Logger: logger,
		Config: config,
	}

	return initializer, nil
}
