package configuration

import (
	"fmt"
	"time"

	"in-memory-db/internal/utils"
)

type Network struct {
	Address        string
	MaxConnections int
	IdleTimeout    time.Duration
	MaxMessageSize int
}

const (
	defaultMessageSize    = "4KB"
	defaultAddress        = "127.0.0.1:54321"
	defaultMaxConnections = 5
	defaultIdleTimeout    = "5m"
)

func ConfigureNetwork(conf *NetworkConfig) (*Network, error) {
	address := defaultAddress
	if conf.Address != "" {
		address = conf.Address
	}

	maxConnections := defaultMaxConnections
	if conf.MaxConnections != 0 {
		maxConnections = conf.MaxConnections
	}

	idleTimeout := defaultIdleTimeout
	if conf.IdleTimeout != "" {
		idleTimeout = conf.IdleTimeout
	}

	timeout, err := time.ParseDuration(idleTimeout)
	if err != nil {
		return nil, fmt.Errorf("invalid idle timeout: %v", err)
	}

	messageSize := defaultMessageSize
	if conf.MaxMessageSize != "" {
		messageSize = conf.MaxMessageSize
	}

	size, err := utils.ParseSize(messageSize)
	if err != nil {
		return nil, fmt.Errorf("invalid message size: %v", err)
	}

	networkConfig := &Network{
		Address:        address,
		MaxConnections: maxConnections,
		IdleTimeout:    timeout,
		MaxMessageSize: size,
	}

	return networkConfig, nil
}
