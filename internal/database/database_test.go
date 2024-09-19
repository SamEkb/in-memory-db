package database

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	"in-memory-db/internal/database/compute"
)

func TestNewDatabase(t *testing.T) {
	logger := zaptest.NewLogger(t)

	storage := &mockStorage{data: make(map[string]string)}
	c := &mockCompute{}

	tests := map[string]struct {
		compute   computeLayer
		storage   storageLayer
		logger    *zap.Logger
		expectErr bool
	}{
		"valid inputs": {
			compute:   c,
			storage:   storage,
			logger:    logger,
			expectErr: false,
		},
		"nil c": {
			compute:   nil,
			storage:   storage,
			logger:    logger,
			expectErr: true,
		},
		"nil storage": {
			compute:   c,
			storage:   nil,
			logger:    logger,
			expectErr: true,
		},
		"nil logger": {
			compute:   c,
			storage:   storage,
			logger:    nil,
			expectErr: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			db, err := NewDatabase(test.compute, test.storage, test.logger)

			if test.expectErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, db)
			}
		})
	}
}

func TestDatabase_HandleQuery(t *testing.T) {
	logger, _ := zap.NewProduction()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			logger.Error("Unable to sync logger", zap.Error(err))
		}
	}(logger)

	storage := &mockStorage{data: map[string]string{"1": "1"}}

	tests := map[string]struct {
		compute   *mockCompute
		queryStr  string
		expected  string
		expectErr bool
	}{
		"successful get query": {
			compute: &mockCompute{
				queryResult: compute.Query{
					CommandId: compute.GetCommandId,
					Arguments: []string{"1"},
				},
			},
			queryStr:  "get 1",
			expected:  "1",
			expectErr: false,
		},
		"successful set query": {
			compute: &mockCompute{
				queryResult: compute.Query{
					CommandId: compute.SetCommandId,
					Arguments: []string{"1", "1"},
				},
			},
			queryStr:  "set 1 1",
			expected:  "Ok",
			expectErr: false,
		},
		"successful delete query": {
			compute: &mockCompute{
				queryResult: compute.Query{
					CommandId: compute.DelCommandId,
					Arguments: []string{"1"},
				},
			},
			queryStr:  "del 1",
			expected:  "Ok",
			expectErr: false,
		},
		"parse error": {
			compute: &mockCompute{
				parseErr: errors.New("parse error"),
			},
			queryStr:  "set 1 value1",
			expected:  "",
			expectErr: true,
		},
		"invalid command ID": {
			compute: &mockCompute{
				queryResult: compute.Query{CommandId: 999},
			},
			queryStr:  "unknown command",
			expected:  "",
			expectErr: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			db, _ := NewDatabase(test.compute, storage, logger)
			result, err := db.HandleQuery(test.queryStr)

			if test.expectErr {
				assert.NotNil(t, err)
				assert.Equal(t, test.expected, result)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, test.expected, result)
			}
		})
	}
}
