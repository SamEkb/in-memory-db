package initialization

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestInitializeDatabase(t *testing.T) {

	tests := map[string]struct {
		logger      *zap.Logger
		expectError bool
	}{
		"initialize logger successfully": {
			logger:      zap.NewNop(),
			expectError: false,
		},
		"failed to initialize logger": {
			logger:      nil,
			expectError: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			db, err := InitializeDatabase(test.logger)
			if test.expectError {
				assert.Error(t, err)
				assert.Nil(t, db)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, db)
			}
		})
	}
}
