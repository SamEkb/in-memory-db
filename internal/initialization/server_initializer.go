package initialization

import (
	"fmt"

	"in-memory-db/internal/database"
	"in-memory-db/internal/database/compute"
	"in-memory-db/internal/database/storage"
	"in-memory-db/internal/database/storage/engine"

	"go.uber.org/zap"
)

func InitializeServer() (*database.Database, *zap.Logger, error) {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{
		"server.log",
		"stderr",
	}

	logger, err := config.Build()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}

	eng := engine.NewEngine(logger)
	sto := storage.NewStorage(eng, logger)

	com, err := compute.NewCompute(logger)
	if err != nil {
		logger.Error("Failed to create new compute", zap.Error(err))
		return nil, nil, err
	}

	db, err := database.NewDatabase(com, sto, logger)
	if err != nil {
		logger.Error("Failed to create new database", zap.Error(err))
		return nil, nil, err
	}
	return db, logger, nil
}
