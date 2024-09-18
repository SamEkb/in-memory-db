package main

import (
	"bufio"
	"fmt"
	"os"

	"go.uber.org/zap"
	"in-memory-db/internal/database"
	"in-memory-db/internal/database/compute"
	"in-memory-db/internal/database/storage"
	"in-memory-db/internal/database/storage/engine"
)

func main() {
	db, logger, err := initialize()
	if err != nil {
		panic(fmt.Sprintf("Initialization error: %v", err))
	}

	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			logger.Error("Unable to sync logger", zap.Error(err))
		}
	}(logger)

	reader := bufio.NewReader(os.Stdin)
	for {
		query, err := reader.ReadString('\n')
		if err != nil {
			logger.Error("Failed to read query", zap.Error(err))
			continue
		}

		res, err := db.HandleQuery(query)
		if err != nil {
			logger.Error("Failed to handle query", zap.String("query", query), zap.Error(err))
			continue
		}

		if len(res) != 0 {
			fmt.Println(res)
		}
	}
}

func initialize() (*database.Database, *zap.Logger, error) {
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
