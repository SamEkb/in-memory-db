package configuration

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigureNetwork(t *testing.T) {
	tests := map[string]struct {
		conf *NetworkConfig
	}{
		"configure network with config values": {
			conf: &NetworkConfig{
				Address:        "1.2.3.4",
				MaxConnections: 1,
				IdleTimeout:    "5m",
				MaxMessageSize: "5kb",
			},
		},
		"configure network with default values": {
			conf: &NetworkConfig{
				Address:        "",
				MaxConnections: 0,
				IdleTimeout:    "",
				MaxMessageSize: "",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			network, err := ConfigureNetwork(test.conf)
			assert.NoError(t, err)
			assert.NotNil(t, network)
		})
	}
}
