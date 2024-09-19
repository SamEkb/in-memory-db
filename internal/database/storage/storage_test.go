package storage

import (
	"fmt"
	"testing"

	"in-memory-db/internal/database/storage/engine"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func initializeLogger(t *testing.T) *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatalf("Unable to create logger: %v", err)
	}
	return logger
}

func initializeStorage(t *testing.T) *Storage {
	logger := initializeLogger(t)
	eng := engine.NewEngine(logger)

	storage := NewStorage(eng, logger)

	return storage
}

func TestNewStorage(t *testing.T) {
	storage := initializeStorage(t)

	assert.NotNil(t, storage, "Storage shouldn't be nil")
	assert.NotNil(t, storage.logger, "Logger shouldn't be nil")
	assert.NotNil(t, storage.engine, "Engine shouldn't be nil")
}

func TestStorage_Set(t *testing.T) {
	storage := initializeStorage(t)

	tests := []struct {
		name   string
		key    string
		value  string
		exists bool
	}{
		{
			name:   "set key and value",
			key:    "1",
			value:  "1",
			exists: true,
		},
		{
			name:   "set another key and value",
			key:    "2",
			value:  "2",
			exists: true,
		},
		{
			name:   "overwrite existing key",
			key:    "1",
			value:  "new_value",
			exists: true,
		},
		{
			name:   "set empty key",
			key:    "",
			value:  "value",
			exists: false,
		},
		{
			name:   "set nil value",
			key:    "3",
			value:  "",
			exists: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.key != "" {
				storage.Set(test.key, test.value)
			}

			value, err := storage.Get(test.key)
			if test.exists {
				assert.Equal(t, test.value, value)
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
			}
		})
	}
}

func TestStorage_Get(t *testing.T) {
	storage := initializeStorage(t)

	storage.Set("1", "1")
	storage.Set("2", "2")

	tests := map[string]struct {
		key           string
		expectedValue string
		error         error
		expectErr     bool
	}{
		"get existing key": {
			key:           "1",
			expectedValue: "1",
			error:         nil,
			expectErr:     false,
		},
		"get another existing key": {
			key:           "2",
			expectedValue: "2",
			error:         nil,
			expectErr:     false,
		},
		"get non-existing key": {
			key:           "3",
			expectedValue: "3",
			error:         fmt.Errorf("record with key: %s doesnt exist", "3"),
			expectErr:     true,
		},
		"get empty key": {
			key:           "",
			expectedValue: "",
			error:         fmt.Errorf("record with key: %s doesnt exist", ""),
			expectErr:     true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			value, err := storage.Get(test.key)
			if test.expectErr {
				assert.Equal(t, err, test.error)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, test.expectedValue, value)
			}
		})
	}
}

func TestStorage_Del(t *testing.T) {
	storage := initializeStorage(t)

	storage.Set("1", "1")
	storage.Set("2", "2")

	tests := map[string]struct {
		key   string
		value string
	}{
		"delete existing key": {
			key:   "1",
			value: "1",
		},
		"delete another existing key": {
			key:   "2",
			value: "2",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			value, err := storage.Get(test.key)
			assert.Nil(t, err, "Expected no error when getting existing key")
			assert.Equal(t, test.value, value, "Expected existing value before deletion")

			storage.Del(test.key)

			value, err = storage.Get(test.key)
			assert.Equal(t, "", value)
			assert.NotNil(t, err)
		})
	}
}
