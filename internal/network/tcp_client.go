package network

import (
	"fmt"
	"io"
	"net"

	"go.uber.org/zap"
	"in-memory-db/internal/initialization"
)

type TcpClient struct {
	connection net.Conn
	App        *initialization.App
}

func NewClient(address string) (*TcpClient, error) {
	app, err := initialization.NewApp()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize app: %v", err)
	}

	conn, err := net.Dial("tcp", address)
	if err != nil {
		return &TcpClient{}, err
	}

	client := &TcpClient{
		App:        app,
		connection: conn,
	}

	return client, nil
}

func (c *TcpClient) Send(message []byte) ([]byte, error) {
	_, err := c.connection.Write(message)
	if err != nil {
		c.App.Logger.Error("Failed to send message", zap.Error(err))
		return []byte{}, err
	}

	reply := make([]byte, 1024)

	resp, err := c.connection.Read(reply)
	if err != nil {
		if err == io.EOF {
			c.App.Logger.Error("Connection closed by server", zap.Error(err))
		} else {
			c.App.Logger.Error("Failed to read reply from server", zap.Error(err))
		}
		return nil, err
	}

	c.App.Logger.Info("Received response from server", zap.Int("bytes", resp))

	fmt.Println(resp)

	return reply[:resp], nil
}

func (c *TcpClient) Close() {
	_ = c.connection.Close()
}
