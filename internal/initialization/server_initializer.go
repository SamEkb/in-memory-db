package initialization

import (
	"fmt"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"path/filepath"

	"in-memory-db/internal/configuration"
	"in-memory-db/internal/database"
	"in-memory-db/internal/database/compute"
	"in-memory-db/internal/database/storage"
	"in-memory-db/internal/database/storage/engine"

	"go.uber.org/zap"
)

type ServerInitializer struct {
	DB     *database.Database
	Logger *zap.Logger
	Config *configuration.Configuration
}

func InitializeServer() (*ServerInitializer, error) {
	conf, err := configuration.NewConfiguration()
	if err != nil {
		return &ServerInitializer{}, err
	}

	logDir := filepath.Dir(conf.Logging.Output)
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err = os.MkdirAll(logDir, os.ModePerm)
		if err != nil {
			return &ServerInitializer{}, fmt.Errorf("failed to create log directory: %v", err)
		}
	}

	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(getLogLevel(conf.Logging.Level))
	config.OutputPaths = []string{
		conf.Logging.Output,
		"stderr",
	}

	logger, err := config.Build()
	if err != nil {
		return &ServerInitializer{}, fmt.Errorf("failed to build logger: %v", err)
	}

	eng := engine.NewEngine(logger)
	sto := storage.NewStorage(eng, logger)

	com, err := compute.NewCompute(logger)
	if err != nil {
		logger.Error("Failed to create new compute", zap.Error(err))
		return &ServerInitializer{}, err
	}

	db, err := database.NewDatabase(com, sto, logger)
	if err != nil {
		logger.Error("Failed to create new database", zap.Error(err))
		return &ServerInitializer{}, err
	}

	initializer := &ServerInitializer{
		DB:     db,
		Logger: logger,
		Config: conf,
	}

	return initializer, nil
}

func getLogLevel(level string) zapcore.Level {
	parsedLevel, err := zapcore.ParseLevel(level)
	if err != nil {
		log.Fatalf("Invalid log level: %v", err)
	}
	return parsedLevel
}
