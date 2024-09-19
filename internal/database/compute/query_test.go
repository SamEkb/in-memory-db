package compute

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewQuery(t *testing.T) {
	query := NewQuery(1, []string{"a"})

	assert.NotNil(t, query.Arguments, "Arguments shouldn't be nil")
	assert.NotNil(t, query.CommandID, "CommandID shouldn't be nil")
}

func TestQuery_GetCommandId(t *testing.T) {
	query := &Query{
		CommandID: 1,
		Arguments: []string{"a"},
	}

	id := query.GetCommandID()
	assert.Equal(t, 1, id, "Expected commandID should be 1")

	query.CommandID = 0
	id = query.GetCommandID()
	assert.Equal(t, 0, id, "Expected commandID should be 0 after change")
}

func TestQuery_GetArguments(t *testing.T) {
	query := &Query{
		CommandID: 1,
		Arguments: []string{"a"},
	}

	args := query.GetArguments()
	assert.NotNil(t, args, "Arguments shouldn't be nil")
	assert.Len(t, args, 1, "Arguments length should be equal")
}
