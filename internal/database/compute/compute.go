package compute

import (
	"errors"
	"fmt"
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
	commandID, ok := commandToID(command)
	if !ok {
		c.logger.Error("invalid command id")
		return Query{}, errors.New("invalid command id")
	}

	number, ok := getArgumentsNumber(commandID)
	if !ok {
		c.logger.Error("invalid arguments")
		return Query{}, errors.New("invalid arguments")
	}

	if len(args[1:]) == 0 {
		c.logger.Error("invalid arguments number")
		return Query{}, fmt.Errorf("invalid number of arguments for command %s, expected %d but got %d", command, number, len(args[1:]))

	}

	query := NewQuery(commandID, args[1:])

	if number != len(query.GetArguments()) {
		c.logger.Error("invalid arguments length for command")
		return Query{}, errors.New("invalid arguments length for command")
	}

	return query, nil
}
