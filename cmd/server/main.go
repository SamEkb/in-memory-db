package main

import (
	"fmt"
	"in-memory-db/internal/initialization"
	"in-memory-db/internal/network"
	"in-memory-db/internal/synchronization"
)

func main() {
	app, err := initialization.NewApp()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize server: %v", err))
	}

	semaphore := synchronization.NewSemaphore(app.Config.Network.MaxConnections)
	server, err := network.NewServer(semaphore, app)
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
