package main

import (
	"bufio"
	"fmt"
	"os"

	"in-memory-db/internal/initialization"

	"go.uber.org/zap"
)

func main() {
	db, logger, err := initialization.Initialize()
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
