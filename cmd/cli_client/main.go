package main

import (
	"bufio"
	"fmt"
	"os"

	"in-memory-db/internal/initialization"
	"in-memory-db/internal/network"

	"go.uber.org/zap"
)

func main() {
	init, err := initialization.InitializeClient()
	if err != nil {
		panic(fmt.Sprintf("Initialization error: %v", err))
	}
	logger := init.Logger

	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			logger.Error("Unable to sync logger", zap.Error(err))
		}
	}(logger)

	reader := bufio.NewReader(os.Stdin)
	conn, err := network.NewClient(init.Config.Network.Address, logger)
	defer conn.Close()
	for {
		query, err := reader.ReadString('\n')
		if err != nil {
			logger.Error("Failed to read query", zap.Error(err))
			continue
		}

		request := []byte(query)

		response, err := conn.Send(request)
		if err != nil {
			logger.Error("Failed to send request", zap.String("query", query), zap.Error(err))
			continue
		}

		if len(response) != 0 {
			fmt.Println(string(response))
		}
	}
}
