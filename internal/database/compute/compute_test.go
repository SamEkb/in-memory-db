package compute

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

func TestNewCompute(t *testing.T) {
	logger := zaptest.NewLogger(t)
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			t.Fatalf("Failed to sync logger: %v", err)
		}
	}(logger)

	compute, err := NewCompute(logger)

	assert.Nil(t, err, "Error should be nil")
	assert.NotNil(t, compute.logger, "Logger shouldn't be nil")
}

func TestCompute_Parse(t *testing.T) {
	logger := zaptest.NewLogger(t)
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			t.Fatalf("Failed to sync logger: %v", err)
		}
	}(logger)

	compute, err := NewCompute(logger)

	assert.Nil(t, err, "Error should be nil")

	tests := map[string]struct {
		query         string
		expectedQuery Query
		expectError   bool
	}{
		"parse valid query": {
			query:         "set a a",
			expectedQuery: Query{CommandID: SetCommandID, Arguments: []string{"a", "a"}},
			expectError:   false,
		},
		"parse empty query": {
			query:         "",
			expectedQuery: Query{},
			expectError:   true,
		},
		"parse invalid command query": {
			query:         "sset a a",
			expectedQuery: Query{},
			expectError:   true,
		},
		"parse invalid arguments number query": {
			query:         "set a",
			expectedQuery: Query{},
			expectError:   true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			query, err := compute.Parse(test.query)
			assert.Equal(t, query, test.expectedQuery)
			if test.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
