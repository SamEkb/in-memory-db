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
	defer server.Close()

	fmt.Println("Server is running...")

	server.AcceptConnections()

	fmt.Println("Server is shutting down...")
}
