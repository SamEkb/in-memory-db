package compute

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommandToID(t *testing.T) {
	setCommandId, ok := commandToId("SET")
	assert.Equal(t, SetCommandId, setCommandId)
	assert.True(t, ok, "Set commandId should exist")

	getCommandId, ok := commandToId("GET")
	assert.Equal(t, GetCommandId, getCommandId)
	assert.True(t, ok, "Get commandId should exist")

	delCommandId, ok := commandToId("DEL")
	assert.Equal(t, DelCommandId, delCommandId)
	assert.True(t, ok, "Del commandId should exist")
}

func TestIdToCommand(t *testing.T) {
	setCommandName, ok := idToCommand(SetCommandId)
	assert.Equal(t, "SET", setCommandName)
	assert.True(t, ok, "Set command name should exist")

	getCommandName, ok := idToCommand(GetCommandId)
	assert.Equal(t, "GET", getCommandName)
	assert.True(t, ok, "Get command name should exist")

	delCommandName, ok := idToCommand(DelCommandId)
	assert.Equal(t, "DEL", delCommandName)
	assert.True(t, ok, "Del command name should exist")
}
