package compute

const (
	GetCommandId = 1
	SetCommandId = 2
	DelCommandId = 3
)

var commandIds = map[int]string{
	GetCommandId: "GET",
	SetCommandId: "SET",
	DelCommandId: "DEL",
}

var commandNames = map[string]int{
	"GET": GetCommandId,
	"SET": SetCommandId,
	"DEL": DelCommandId,
}

var argumentsNumber = map[int]int{
	GetCommandId: 1,
	SetCommandId: 2,
	DelCommandId: 1,
}

func commandToId(command string) (int, bool) {
	commandId, ok := commandNames[command]
	return commandId, ok
}

func idToCommand(commandId int) (string, bool) {
	commandName, ok := commandIds[commandId]
	return commandName, ok
}

func getArgumentsNumber(commandId int) (int, bool) {
	number, ok := argumentsNumber[commandId]
	return number, ok
}
