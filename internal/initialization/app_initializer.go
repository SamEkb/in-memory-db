package initialization

import (
	"fmt"

	"in-memory-db/internal/configuration"
	"in-memory-db/internal/database"

	"go.uber.org/zap"
)

type App struct {
	Logger *zap.Logger
	DB     *database.Database
	Config *Config
}

type Config struct {
	Network *configuration.Network
	Logging *configuration.LoggingConfig
}

func NewApp() (*App, error) {
	conf, err := configuration.NewConfiguration()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %v", err)
	}

	logger, err := InitializeLogger(conf.Logging)
	if err != nil {
		return nil, fmt.Errorf("failed to build logger: %v", err)
	}

	db, err := InitializeDatabase(logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %v", err)
	}

	networkConfig, err := configuration.ConfigureNetwork(conf.Network)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize network: %v", err)
	}

	return &App{
		Logger: logger,
		DB:     db,
		Config: &Config{
			Network: networkConfig,
			Logging: conf.Logging,
		},
	}, nil
}
