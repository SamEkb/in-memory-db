package initialization

import (
	"errors"
	"testing"

	"in-memory-db/internal/configuration"
	"in-memory-db/internal/database"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNewApp(t *testing.T) {
	tests := map[string]struct {
		patchConfig        func()
		patchLogger        func()
		patchDatabase      func()
		patchNetworkConfig func()
		expectError        bool
	}{
		"successful app initialization": {
			patchConfig: func() {
				monkey.Patch(configuration.NewConfiguration, func() (*configuration.Config, error) {
					return &configuration.Config{
						Network: &configuration.NetworkConfig{
							Address:        "127.0.0.1:8080",
							MaxConnections: 10,
						},
						Logging: &configuration.LoggingConfig{
							Level:  "info",
							Output: "stdout",
						},
					}, nil
				})
			},
			patchLogger: func() {
				monkey.Patch(InitializeLogger, func(*configuration.LoggingConfig) (*zap.Logger, error) {
					return zap.NewNop(), nil
				})
			},
			patchDatabase: func() {
				monkey.Patch(InitializeDatabase, func(*zap.Logger) (*database.Database, error) {
					return &database.Database{}, nil
				})
			},
			patchNetworkConfig: func() {
				monkey.Patch(configuration.ConfigureNetwork, func(*configuration.NetworkConfig) (*configuration.Network, error) {
					return &configuration.Network{}, nil
				})
			},
			expectError: false,
		},
		"failed to load configuration": {
			patchConfig: func() {
				monkey.Patch(configuration.NewConfiguration, func() (*configuration.Config, error) {
					return nil, errors.New("config load error")
				})
			},
			expectError: true,
		},
		"failed to build logger": {
			patchConfig: func() {
				monkey.Patch(configuration.NewConfiguration, func() (*configuration.Config, error) {
					return &configuration.Config{
						Logging: &configuration.LoggingConfig{
							Level:  "info",
							Output: "stdout",
						},
					}, nil
				})
			},
			patchLogger: func() {
				monkey.Patch(InitializeLogger, func(*configuration.LoggingConfig) (*zap.Logger, error) {
					return nil, errors.New("logger error")
				})
			},
			expectError: true,
		},
		"failed to initialize database": {
			patchConfig: func() {
				monkey.Patch(configuration.NewConfiguration, func() (*configuration.Config, error) {
					return &configuration.Config{}, nil
				})
			},
			patchLogger: func() {
				monkey.Patch(InitializeLogger, func(*configuration.LoggingConfig) (*zap.Logger, error) {
					return zap.NewNop(), nil
				})
			},
			patchDatabase: func() {
				monkey.Patch(InitializeDatabase, func(*zap.Logger) (*database.Database, error) {
					return nil, errors.New("db error")
				})
			},
			expectError: true,
		},
		"failed to initialize network": {
			patchConfig: func() {
				monkey.Patch(configuration.NewConfiguration, func() (*configuration.Config, error) {
					return &configuration.Config{}, nil
				})
			},
			patchLogger: func() {
				monkey.Patch(InitializeLogger, func(*configuration.LoggingConfig) (*zap.Logger, error) {
					return zap.NewNop(), nil
				})
			},
			patchDatabase: func() {
				monkey.Patch(InitializeDatabase, func(*zap.Logger) (*database.Database, error) {
					return &database.Database{}, nil
				})
			},
			patchNetworkConfig: func() {
				monkey.Patch(configuration.ConfigureNetwork, func(config *configuration.NetworkConfig) (*configuration.Network, error) {
					return nil, errors.New("network error")
				})
			},
			expectError: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			defer monkey.UnpatchAll()

			if test.patchConfig != nil {
				test.patchConfig()
			}
			if test.patchLogger != nil {
				test.patchLogger()
			}
			if test.patchDatabase != nil {
				test.patchDatabase()
			}
			if test.patchNetworkConfig != nil {
				test.patchNetworkConfig()
			}

			app, err := NewApp()

			if test.expectError {
				assert.Nil(t, app)
				assert.Error(t, err)
			} else {
				assert.NotNil(t, app)
				assert.NoError(t, err)
			}
		})
	}
}
