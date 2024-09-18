package database

import (
	"errors"
	"fmt"

	"go.uber.org/zap"
	"in-memory-db/internal/database/compute"
)

type storageLayer interface {
	Set(string, string)
	Get(string) (string, error)
	Del(string)
}

type computeLayer interface {
	Parse(string) (compute.Query, error)
}

type Database struct {
	compute computeLayer
	storage storageLayer
	logger  *zap.Logger
}

func NewDatabase(compute computeLayer, storage storageLayer, logger *zap.Logger) (*Database, error) {
	if compute == nil {
		return nil, errors.New("compute is empty")
	}

	if storage == nil {
		return nil, errors.New("storage is empty")
	}

	if logger == nil {
		return nil, errors.New("logger is empty")
	}

	return &Database{
		compute: compute,
		storage: storage,
		logger:  logger,
	}, nil
}

func (d *Database) HandleQuery(queryStr string) (string, error) {
	query, err := d.compute.Parse(queryStr)
	if err != nil {
		d.logger.Error("failed to parse query string", zap.String("query", queryStr), zap.Error(err))
		return "", err
	}

	commandId := query.GetCommandId()

	switch commandId {
	case compute.GetCommandId:
		return d.handleGetQuery(query)
	case compute.SetCommandId:
		return d.handleSetQuery(query)
	case compute.DelCommandId:
		return d.handleDelQuery(query)
	default:
		d.logger.Error("Invalid command ID", zap.Int("commandId", commandId), zap.Error(err))
		return "", fmt.Errorf("non-existent command ID: %d", commandId)
	}
}

func (d *Database) handleGetQuery(query compute.Query) (string, error) {
	args := query.GetArguments()
	if len(args) == 0 {
		d.logger.Error("No arguments provided for GET query")
		return "", errors.New("no arguments provided")
	}

	key := args[0]
	result, err := d.storage.Get(key)
	if err != nil {
		d.logger.Error("Get operation failed", zap.String("key", key), zap.Error(err))
		return "", fmt.Errorf("get failed for key %s: %w", args[0], err)
	}

	return result, nil
}

func (d *Database) handleSetQuery(query compute.Query) (string, error) {
	args := query.GetArguments()
	if len(args) == 0 {
		d.logger.Error("No arguments provided for SET query")
		return "", errors.New("no arguments provided")
	}

	d.storage.Set(args[0], args[1])

	return "Ok", nil
}

func (d *Database) handleDelQuery(query compute.Query) (string, error) {
	args := query.GetArguments()
	if len(args) == 0 {
		d.logger.Error("No arguments provided for DEL query")
		return "", errors.New("no arguments provided")
	}

	d.storage.Del(args[0])

	return "Ok", nil
}
