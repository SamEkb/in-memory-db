package initialization

import (
	"go.uber.org/zap"
	"in-memory-db/internal/database"
	"in-memory-db/internal/database/compute"
	"in-memory-db/internal/database/storage"
	"in-memory-db/internal/database/storage/engine"
)

func Initialize() (*database.Database, *zap.Logger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		panic("Unable to initialize logger")
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
