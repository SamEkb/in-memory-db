package storage

import (
	"fmt"

	"go.uber.org/zap"
)

type Engine interface {
	Put(string, string)
	Get(string) (string, error)
	Del(string)
}

type Storage struct {
	engine Engine
	logger *zap.Logger
}

func NewStorage(engine Engine, logger *zap.Logger) *Storage {
	return &Storage{engine: engine, logger: logger}
}

func (s *Storage) Set(key string, value string) {
	s.engine.Put(key, value)
}

func (s *Storage) Get(key string) (string, error) {
	result, err := s.engine.Get(key)
	if err != nil {
		s.logger.Error("Get operation failed: record doesn't exist", zap.String("key", key), zap.Error(err))
		return "", fmt.Errorf("record with key: %s doesnt exist", key)
	}

	return result, nil
}

func (s *Storage) Del(key string) {
	s.engine.Del(key)
}
