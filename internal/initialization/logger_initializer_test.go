package initialization

import (
	"testing"

	"in-memory-db/internal/configuration"

	"github.com/stretchr/testify/assert"
)

func TestInitializeLogger(t *testing.T) {
	conf := &configuration.LoggingConfig{
		Level:  "debug",
		Output: defaultOutput,
	}

	_, err := InitializeLogger(conf)
	assert.NoError(t, err)

	tests := map[string]struct {
		conf        *configuration.LoggingConfig
		expectError bool
	}{
		"initialize logger successfully": {
			conf: &configuration.LoggingConfig{
				Level:  "debug",
				Output: defaultOutput,
			},
			expectError: false,
		},
		"initialize with invalid logging level": {
			conf: &configuration.LoggingConfig{
				Level:  "invalid",
				Output: defaultOutput,
			},
			expectError: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			logger, err := InitializeLogger(test.conf)
			if test.expectError {
				assert.Nil(t, logger)
				assert.Error(t, err)
			} else {
				assert.NotNil(t, logger)
				assert.NoError(t, err)
			}
		})
	}
}
