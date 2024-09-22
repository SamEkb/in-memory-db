package storage

import (
	"go.uber.org/zap"
)

type Engine interface {
	Set(string, string)
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
	s.engine.Set(key, value)
}

func (s *Storage) Get(key string) (string, error) {
	result, err := s.engine.Get(key)
	if err != nil {
		return "", err
	}
	return result, nil
}

func (s *Storage) Del(key string) {
	s.engine.Del(key)
}
