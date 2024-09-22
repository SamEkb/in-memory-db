package engine

import (
	"fmt"
	"testing"

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

func initializeEngine(logger *zap.Logger) *Engine {
	return &Engine{hashtable: &Hashtable{
		data: map[string]string{
			"1": "1",
			"2": "2",
		},
	}, logger: logger}
}

func TestNewEngine(t *testing.T) {
	logger := initializeLogger(t)
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			logger.Error("Unable to sync logger", zap.Error(err))
		}
	}(logger)

	engine := initializeEngine(logger)

	assert.NotNil(t, engine, "Engine should be initialized")
	assert.NotNil(t, engine.logger, "Logger should be initialized")
	assert.NotNil(t, engine.hashtable, "Hashtable should be initialized")
}

func TestEngine_Get(t *testing.T) {
	logger := initializeLogger(t)
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			logger.Error("Unable to sync logger", zap.Error(err))
		}
	}(logger)

	engine := initializeEngine(logger)

	tests := map[string]struct {
		key           string
		expectedValue string
		error         error
	}{
		"get existing key": {
			key:           "1",
			expectedValue: "1",
			error:         nil,
		},
		"get non existing key": {
			key:           "4",
			expectedValue: "",
			error:         fmt.Errorf("record with key: %s doesn't exist", "4"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			value, err := engine.Get(test.key)
			assert.Equal(t, value, test.expectedValue)
			assert.Equal(t, err, test.error)
		})
	}
}

func TestEngine_Put(t *testing.T) {
	logger := initializeLogger(t)
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			logger.Error("Unable to sync logger", zap.Error(err))
		}
	}(logger)

	engine := initializeEngine(logger)

	tests := map[string]struct {
		key   string
		value string
	}{
		"put new key and value": {
			key:   "4",
			value: "4",
		},
		"put existing key": {
			key:   "1",
			value: "0",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			engine.Set(test.key, test.value)
			value, err := engine.Get(test.key)
			assert.Equal(t, value, test.value)
			assert.Nil(t, err, "Error should be nil")
		})
	}
}

func TestEngine_Del(t *testing.T) {
	logger := initializeLogger(t)
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			logger.Error("Unable to sync logger", zap.Error(err))
		}
	}(logger)

	engine := initializeEngine(logger)

	tests := map[string]struct {
		key      string
		initial  string
		afterDel string
		exists   bool
	}{
		"delete existing key": {
			key:      "1",
			initial:  "1",
			afterDel: "",
			exists:   true,
		},
		"delete another existing key": {
			key:      "2",
			initial:  "2",
			afterDel: "",
			exists:   true,
		},
		"delete non-existing key": {
			key:      "3",
			initial:  "",
			afterDel: "",
			exists:   false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			value, err := engine.Get(test.key)
			assert.Equal(t, test.initial, value)
			assert.Equal(t, test.exists, err == nil)

			engine.Del(test.key)

			value, err = engine.Get(test.key)
			assert.Equal(t, test.afterDel, value)
			assert.NotNil(t, err)
		})
	}
}
