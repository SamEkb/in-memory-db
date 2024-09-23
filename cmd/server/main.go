package main

import (
	"fmt"

	"in-memory-db/internal/initialization"
	"in-memory-db/internal/network"
)

func main() {
	db, logger, err := initialization.InitializeServer()
	if err != nil {
		panic(fmt.Sprintf("Initialization error: %v", err))
	}

	server, err := network.NewServer("127.0.0.1:3223", logger, db)
	if err != nil {
		panic(fmt.Sprintf("Failed to start server: %v", err))
	}
	defer server.Close()

	fmt.Println("Server is running...")

	server.AcceptConnections()

	fmt.Println("Server is shutting down...")
}
