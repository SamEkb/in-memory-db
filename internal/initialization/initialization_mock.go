package initialization

import (
	"in-memory-db/internal/database"

	"go.uber.org/zap"
)

type MockSemaphore struct {
	acquireCalled bool
	releaseCalled bool
}

func (m *MockSemaphore) Acquire() {
	m.acquireCalled = true
}

func (m *MockSemaphore) Release() {
	m.releaseCalled = true
}

type MockApp struct {
	mockLogger *zap.Logger
	mockDB     database.IDatabase
	mockConfig *Config
}

func (m *MockApp) GetLogger() *zap.Logger {
	return m.mockLogger
}

func (m *MockApp) GetDataBase() database.IDatabase {
	return m.mockDB
}

func (m *MockApp) GetConfig() *Config {
	return m.mockConfig
}

func NewMockApp(logger *zap.Logger, db database.IDatabase, config *Config) *MockApp {
	return &MockApp{
		mockLogger: logger,
		mockDB:     db,
		mockConfig: config,
	}
}
