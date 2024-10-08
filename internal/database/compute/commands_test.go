package compute

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommandToID(t *testing.T) {
	setCommandID, ok := commandToID("SET")
	assert.Equal(t, SetCommandID, setCommandID)
	assert.True(t, ok, "Set commandID should exist")

	getCommandID, ok := commandToID("GET")
	assert.Equal(t, GetCommandID, getCommandID)
	assert.True(t, ok, "Get commandID should exist")

	delCommandID, ok := commandToID("DEL")
	assert.Equal(t, DelCommandID, delCommandID)
	assert.True(t, ok, "Del commandID should exist")
}

func TestIdToCommand(t *testing.T) {
	setCommandName, ok := IDToCommand(SetCommandID)
	assert.Equal(t, "SET", setCommandName)
	assert.True(t, ok, "Set command name should exist")

	getCommandName, ok := IDToCommand(GetCommandID)
	assert.Equal(t, "GET", getCommandName)
	assert.True(t, ok, "Get command name should exist")

	delCommandName, ok := IDToCommand(DelCommandID)
	assert.Equal(t, "DEL", delCommandName)
	assert.True(t, ok, "Del command name should exist")
}
