package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"in-memory-db/internal/network"

	"go.uber.org/zap"
)

func main() {
	address := flag.String("address", "localhost:54321", "Address of the server to connect to")
	flag.Parse()

	reader := bufio.NewReader(os.Stdin)
	client, err := network.NewClient(*address)
	if err != nil {
		panic(fmt.Sprintf("Failed to start client: %v", err))
	}
	defer func(client *network.TcpClient) {
		err := client.Close()
		if err != nil {
			panic(fmt.Sprintf("Failed to close client: %v", err))
		}
	}(client)
	for {
		query, err := reader.ReadString('\n')
		if err != nil {
			client.App.Logger.Error("Failed to read query", zap.Error(err))
			continue
		}

		request := []byte(query)

		response, err := client.Send(request)
		if err != nil {
			client.App.Logger.Error("Failed to send request", zap.String("query", query), zap.Error(err))
			continue
		}

		if len(response) != 0 {
			fmt.Println(string(response))
		}
	}
}
