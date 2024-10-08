package database

import (
	"errors"
	"testing"

	"in-memory-db/internal/database/compute"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
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
	logger := zaptest.NewLogger(t)

	tests := map[string]struct {
		compute   *mockCompute
		storage   *mockStorage
		queryStr  string
		expected  string
		expectErr bool
	}{
		"successful get query": {
			compute: &mockCompute{
				queryResult: compute.Query{
					CommandID: compute.GetCommandID,
					Arguments: []string{"1"},
				},
			},
			storage:   &mockStorage{data: map[string]string{"1": "1"}},
			queryStr:  "get 1",
			expected:  "1",
			expectErr: false,
		},
		"successful set query": {
			compute: &mockCompute{
				queryResult: compute.Query{
					CommandID: compute.SetCommandID,
					Arguments: []string{"1", "1"},
				},
			},
			storage:   &mockStorage{data: map[string]string{}},
			queryStr:  "set 1 1",
			expected:  "Ok",
			expectErr: false,
		},
		"successful delete query": {
			compute: &mockCompute{
				queryResult: compute.Query{
					CommandID: compute.DelCommandID,
					Arguments: []string{"1"},
				},
			},
			storage:   &mockStorage{data: map[string]string{"1": "1"}},
			queryStr:  "del 1",
			expected:  "Ok",
			expectErr: false,
		},
		"parse error": {
			compute: &mockCompute{
				parseErr: errors.New("parse error"),
			},
			storage:   &mockStorage{data: map[string]string{}},
			queryStr:  "set 1 value1",
			expected:  "",
			expectErr: true,
		},
		"invalid command ID": {
			compute: &mockCompute{
				queryResult: compute.Query{CommandID: 999},
			},
			storage:   &mockStorage{data: map[string]string{}},
			queryStr:  "unknown command",
			expected:  "",
			expectErr: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			db, _ := NewDatabase(test.compute, test.storage, logger)
			result, err := db.HandleQuery(test.queryStr)

			if test.expectErr {
				assert.Error(t, err)
				assert.Equal(t, test.expected, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, result)
			}
		})
	}
}
