package network

import (
	"errors"
	"in-memory-db/internal/database"
	"net"
	"sync"
	"testing"
	"time"

	"in-memory-db/internal/configuration"
	"in-memory-db/internal/initialization"
	"in-memory-db/internal/synchronization"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNewServer(t *testing.T) {
	t.Run("successful server initialization", func(t *testing.T) {
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

		monkey.Patch(net.Listen, func(network, address string) (net.Listener, error) {
			return &mockListener{}, nil
		})

		defer monkey.UnpatchAll()

		server, err := initializeMockServer(false)
		assert.NotNil(t, server)
		assert.Nil(t, err)
	})

	t.Run("failed to initialize app", func(t *testing.T) {
		monkey.Patch(initialization.NewApp, func() (*initialization.App, error) {
			return nil, errors.New("failed to initialize app")
		})
		defer monkey.UnpatchAll()

		server, err := initializeMockServer(true)
		assert.Nil(t, server)
		assert.NotNil(t, err)
	})

	t.Run("failed to listen on address", func(t *testing.T) {
		monkey.Patch(net.Listen, func(network, address string) (net.Listener, error) {
			return nil, errors.New("failed to listen on address")
		})
		defer monkey.UnpatchAll()

		server, err := initializeMockServer(false)
		assert.Nil(t, server)
		assert.NotNil(t, err)
	})
}

func TestTcpServer_AcceptConnections(t *testing.T) {
	mockConn := &mockConn{}
	mockListener := &mockListener{acceptConn: mockConn}

	app := &initialization.App{
		Logger: zap.NewNop(),
		Config: &initialization.Config{
			Network: &configuration.Network{
				IdleTimeout:    5 * time.Second,
				MaxConnections: 1,
			},
		},
	}

	sem := synchronization.NewSemaphore(app.Config.Network.MaxConnections)

	server := &TcpServer{
		listener:  mockListener,
		app:       app,
		semaphore: sem,
	}

	var wg sync.WaitGroup

	blocked := false
	wg.Add(2)

	go func() {
		defer wg.Done()
		server.semaphore.Acquire()
		time.Sleep(200 * time.Millisecond)
		server.semaphore.Release()
	}()

	go func() {
		defer wg.Done()
		time.Sleep(50 * time.Millisecond)
		start := time.Now()
		server.semaphore.Acquire()
		elapsed := time.Since(start)
		if elapsed > 150*time.Millisecond {
			blocked = true
		}
		server.semaphore.Release()
	}()

	wg.Wait()

	assert.True(t, blocked, "Second go was blocked")
}

func TestTcpServer_Close(t *testing.T) {
	tests := map[string]struct {
		mockListener  *mockListener
		expectError   bool
		expectedError string
	}{
		"successfully close server": {
			mockListener: &mockListener{
				closeFunc: func() error {
					return nil
				},
			},
			expectError: false,
		},
		"failed to close server": {
			mockListener: &mockListener{
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
			server := &TcpServer{
				listener: test.mockListener,
				app: &initialization.App{
					Logger: zap.NewNop(),
				},
			}

			err := server.Close()

			if test.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func initializeMockApp(shouldFail bool) (*initialization.MockApp, error) {
	if shouldFail {
		return nil, errors.New("failed to initialize app")
	}

	mockLogger := zap.NewNop()
	mockDB := &database.MockDatabase{}
	mockConfig := &initialization.Config{
		Network: &configuration.Network{
			Address:        "127.0.0.1:8080",
			MaxMessageSize: 1024,
			IdleTimeout:    5 * time.Minute,
			MaxConnections: 10,
		},
		Logging: &configuration.LoggingConfig{
			Level: "info",
		},
	}

	return initialization.NewMockApp(mockLogger, mockDB, mockConfig), nil
}

func initializeMockServer(shouldFail bool) (*TcpServer, error) {
	mockSemaphore := &initialization.MockSemaphore{}
	mockApp, err := initializeMockApp(shouldFail)
	if err != nil {
		return nil, err
	}

	return NewServer(mockSemaphore, mockApp)
}
