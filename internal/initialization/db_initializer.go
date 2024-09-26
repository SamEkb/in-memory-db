package initialization

import (
	"in-memory-db/internal/database"
	"in-memory-db/internal/database/compute"
	"in-memory-db/internal/database/storage"
	"in-memory-db/internal/database/storage/engine"

	"go.uber.org/zap"
)

func InitializeDatabase(logger *zap.Logger) (*database.Database, error) {
	eng := engine.NewEngine(logger)
	sto := storage.NewStorage(eng, logger)

	com, err := compute.NewCompute(logger)
	if err != nil {
		logger.Error("Failed to create new compute", zap.Error(err))
		return nil, err
	}

	db, err := database.NewDatabase(com, sto, logger)
	if err != nil {
		logger.Error("Failed to create new database", zap.Error(err))
		return nil, err
	}

	return db, nil
}
