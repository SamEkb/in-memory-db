package database

import (
	"errors"

	"in-memory-db/internal/database/compute"
)

type mockCompute struct {
	queryResult compute.Query
	parseErr    error
}

func (m *mockCompute) Parse(query string) (compute.Query, error) {
	return m.queryResult, m.parseErr
}

type mockStorage struct {
	data map[string]string
}

func (m *mockStorage) Set(key, value string) {
	m.data[key] = value
}

func (m *mockStorage) Get(key string) (string, error) {
	value, exists := m.data[key]
	if !exists {
		return "", errors.New("not found")
	}
	return value, nil
}

func (m *mockStorage) Del(key string) {
	delete(m.data, key)
}
