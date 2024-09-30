package configuration

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfiguration(t *testing.T) {
	tests := map[string]struct {
		setupConfig      func()
		expectError      bool
		expectedErrorMsg string
	}{
		"successful config loading": {
			setupConfig: func() {
				validConfig := `
network:
  address: "127.0.0.1:8080"
  max_connections: 100
  max_message_size: "4KB"
  idle_timeout: "5m"
engine:
  type: "in-memory"
logging:
  level: "info"
  output: "stdout"
`
				err := os.WriteFile(config, []byte(validConfig), 0600)
				assert.NoError(t, err)
			},
			expectError: false,
		},
		"config file not found": {
			setupConfig: func() {
				err := os.Remove(config)
				if err != nil {
					return
				}
			},
			expectError:      true,
			expectedErrorMsg: "no such file or directory",
		},
		"invalid YAML format": {
			setupConfig: func() {
				invalidConfig := `
invalid_yaml_content
`
				err := os.WriteFile(config, []byte(invalidConfig), 0600)
				assert.NoError(t, err)
			},
			expectError:      true,
			expectedErrorMsg: "yaml: unmarshal errors",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if test.setupConfig != nil {
				test.setupConfig()
			}

			configFile, err := NewConfiguration()

			if test.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), test.expectedErrorMsg)
				assert.Nil(t, configFile)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, configFile)
				assert.Equal(t, "127.0.0.1:8080", configFile.Network.Address)
				assert.Equal(t, 100, configFile.Network.MaxConnections)
				assert.Equal(t, "4KB", configFile.Network.MaxMessageSize)
			}

			err = os.Remove(config)
			if err != nil {
				return
			}
		})
	}
}
