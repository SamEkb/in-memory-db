package network

import (
	"errors"
	"net"
	"testing"

	"in-memory-db/internal/configuration"
	"in-memory-db/internal/initialization"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNewClient(t *testing.T) {
	t.Run("successful client initialization", func(t *testing.T) {
		monkey.Patch(initialization.NewApp, func() (*initialization.App, error) {
			return &initialization.App{
				Logger: zap.NewNop(),
				Config: &initialization.Config{
					Network: &configuration.Network{
						Address:        "127.0.0.1:0",
						MaxConnections: 10,
					},
				},
			}, nil
		})

		monkey.Patch(net.Dial, func(network, address string) (net.Conn, error) {
			return &mockConn{}, nil
		})
		defer monkey.UnpatchAll()

		client, err := NewClient("127.0.0.1:0")
		assert.NotNil(t, client)
		assert.NoError(t, err)
	})

	t.Run("failed to initialize client", func(t *testing.T) {
		monkey.Patch(initialization.NewApp, func() (*initialization.App, error) {
			return nil, errors.New("failed to initialize app")
		})

		defer monkey.UnpatchAll()

		client, err := NewClient("127.0.0.1:0")
		assert.Nil(t, client)
		assert.Error(t, err)
	})

	t.Run("failed to dial", func(t *testing.T) {
		monkey.Patch(net.Dial, func(network, address string) (net.Conn, error) {
			return &mockConn{}, nil
		})

		defer monkey.UnpatchAll()

		client, err := NewClient("127.0.0.1:0")
		assert.Nil(t, client)
		assert.Error(t, err)
	})
}

func TestTcpClient_Send(t *testing.T) {
	t.Run("successful send and receive", func(t *testing.T) {
		mockConnection := &mockConn{
			writeFunc: func(b []byte) (int, error) {
				return len(b), nil
			},
			readFunc: func(b []byte) (int, error) {
				copy(b, "response from server")
				return len("response from server"), nil
			},
		}

		client := &TcpClient{
			connection: mockConnection,
			App: &initialization.App{
				Logger: zap.NewNop(),
			},
		}

		response, err := client.Send([]byte("test message"))
		assert.NoError(t, err, "Expected no error on send and receive")
		assert.Equal(t, "response from server", string(response), "Expected correct server response")
	})

	t.Run("failed to write message", func(t *testing.T) {
		mockConnection := &mockConn{
			writeFunc: func(b []byte) (int, error) {
				return len(b), errors.New("failed to write message")
			},
			readFunc: func(b []byte) (int, error) {
				return len("response from server"), nil
			},
		}

		client := &TcpClient{
			connection: mockConnection,
			App: &initialization.App{
				Logger: zap.NewNop(),
			},
		}

		response, err := client.Send([]byte("test message"))
		assert.Nil(t, response)
		assert.Error(t, err)
	})

	t.Run("failed to read message", func(t *testing.T) {
		mockConnection := &mockConn{
			writeFunc: func(b []byte) (int, error) {
				return len(b), nil
			},
			readFunc: func(b []byte) (int, error) {
				return len("response from server"), errors.New("failed to read message")
			},
		}

		client := &TcpClient{
			connection: mockConnection,
			App: &initialization.App{
				Logger: zap.NewNop(),
			},
		}

		response, err := client.Send([]byte("test message"))
		assert.Nil(t, response)
		assert.Error(t, err)
	})
}

func TestTcpClient_Close(t *testing.T) {
	tests := map[string]struct {
		mockConnection *mockConn
		expectError    bool
		expectedError  string
	}{
		"successfully close client": {
			mockConnection: &mockConn{
				closeFunc: func() error {
					return nil
				},
			},
			expectError: false,
		},
		"failed to close client": {
			mockConnection: &mockConn{
				closeFunc: func() error {
					return errors.New("close error")
				},
			},
			expectError:   true,
			expectedError: "close error",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &TcpClient{
				connection: test.mockConnection,
				App: &initialization.App{
					Logger: zap.NewNop(),
				},
			}

			err := client.Close()
			if test.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
