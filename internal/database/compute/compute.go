package compute

import (
	"errors"
	"strings"

	"go.uber.org/zap"
)

type Compute struct {
	logger *zap.Logger
}

func NewCompute(logger *zap.Logger) (*Compute, error) {
	if logger == nil {
		return nil, errors.New("logger is invalid")
	}

	return &Compute{
		logger: logger,
	}, nil
}

func (c *Compute) Parse(queryInput string) (Query, error) {
	args := strings.Fields(queryInput)

	if len(args) == 0 {
		c.logger.Error("invalid query length")
		return Query{}, errors.New("invalid query length")
	}

	command := strings.ToUpper(args[0])
	commandId, ok := commandToId(command)
	if !ok {
		c.logger.Error("invalid command id")
		return Query{}, errors.New("invalid command id")
	}

	number, ok := getArgumentsNumber(commandId)
	if !ok {
		c.logger.Error("invalid arguments")
		return Query{}, errors.New("invalid arguments")
	}

	query := NewQuery(commandId, args[1:])

	if number != len(query.GetArguments()) {
		c.logger.Error("invalid arguments number")
		return Query{}, errors.New("invalid arguments number")
	}

	return query, nil
}
