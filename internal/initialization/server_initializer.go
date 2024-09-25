package initialization

import (
	"fmt"

	"go.uber.org/zap"
	"in-memory-db/internal/configuration"
	"in-memory-db/internal/database"
)

type ServerInitializer struct {
	DB      *database.Database
	Logger  *zap.Logger
	Network *Network
}

type Network struct {
	Address        string
	MaxConnections int
	IdleTimeout    string
	MaxMessageSize string
}

const defaultMessageSize = "4KB"
const defaultAddress = "127.0.0.1:54321"
const defaultMaxConnections = 5
const defaultIdleTimeout = "5m"

func InitializeServer() (*ServerInitializer, error) {
	address := defaultAddress
	maxMessageSize := defaultMessageSize
	maxConnections := defaultMaxConnections
	idleTimeout := defaultIdleTimeout
	conf, err := configuration.NewConfiguration()
	if err != nil {
		return &ServerInitializer{}, err
	}

	logger, err := InitializeLogger(conf.Logging)
	if err != nil {
		return &ServerInitializer{}, fmt.Errorf("failed to build logger: %v", err)
	}

	db, err := InitializeDatabase(logger)
	if err != nil {
		return &ServerInitializer{}, fmt.Errorf("failed to initialize database %v", err)
	}

	if conf.Network.Address != "" {
		address = conf.Network.Address
	}
	if conf.Network.MaxConnection != 0 {
		maxConnections = conf.Network.MaxConnection
	}
	if conf.Network.IdleTimeout != "" {
		idleTimeout = conf.Network.IdleTimeout
	}
	if conf.Network.MaxMessageSize != "" {
		maxMessageSize = conf.Network.MaxMessageSize
	}

	network := &Network{
		Address:        address,
		MaxConnections: maxConnections,
		IdleTimeout:    idleTimeout,
		MaxMessageSize: maxMessageSize,
	}

	initializer := &ServerInitializer{
		DB:      db,
		Logger:  logger,
		Network: network,
	}

	return initializer, nil
}
