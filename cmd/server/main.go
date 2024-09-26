package main

import (
	"fmt"

	"in-memory-db/internal/network"
)

func main() {
	server, err := network.NewServer()
	if err != nil {
		panic(fmt.Sprintf("Failed to start server: %v", err))
	}
	defer func(server *network.TcpServer) {
		err := server.Close()
		if err != nil {
			panic(fmt.Sprintf("Failed to close server: %v", err))
		}
	}(server)

	fmt.Println("Server is running...")

	server.AcceptConnections()

	fmt.Println("Server is shutting down...")
}
