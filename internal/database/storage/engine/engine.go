package engine

import (
	"fmt"

	"go.uber.org/zap"
)

type Engine struct {
	hashtable *Hashtable
	logger    *zap.Logger
}

func NewEngine(logger *zap.Logger) *Engine {
	hashtable := NewHashtable()

	return &Engine{
		hashtable: hashtable,
		logger:    logger,
	}
}

func (e *Engine) Set(key string, value string) {
	e.hashtable.Put(key, value)
}

func (e *Engine) Get(key string) (string, error) {
	result, ok := e.hashtable.Get(key)
	if !ok {
		return "", fmt.Errorf("record with key: %s doesn't exist", key)
	}
	return result, nil
}

func (e *Engine) Del(key string) {
	e.hashtable.Del(key)
}
