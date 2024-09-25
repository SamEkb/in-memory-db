package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"in-memory-db/internal/initialization"
	"in-memory-db/internal/network"

	"go.uber.org/zap"
)

func main() {
	address := flag.String("address", "localhost:3223", "Address of the server to connect to")

	flag.Parse()

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

	if init.Config.Address != "" {
		address = &init.Config.Address
	}

	reader := bufio.NewReader(os.Stdin)
	conn, err := network.NewClient(*address, logger)
	if err != nil {
		panic(fmt.Sprintf("Failed to start client: %v", err))
	}
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
