package network

import (
	"fmt"
	"net"

	"go.uber.org/zap"
)

type TcpClient struct {
	connection net.Conn
	logger     *zap.Logger
}

func NewClient(address string, logger *zap.Logger) (*TcpClient, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return &TcpClient{}, err
	}

	client := &TcpClient{
		logger:     logger,
		connection: conn,
	}

	return client, nil
}

func (c *TcpClient) Send(message []byte) ([]byte, error) {
	_, err := c.connection.Write(message)
	if err != nil {
		return []byte{}, err
	}

	reply := make([]byte, 1024)

	resp, err := c.connection.Read(reply)

	fmt.Println(resp)

	return reply[:resp], nil
}

func (c *TcpClient) Close() {
	_ = c.connection.Close()
}
