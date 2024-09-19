package compute

const (
	GetCommandID = 1
	SetCommandID = 2
	DelCommandID = 3
)

var commandIds = map[int]string{
	GetCommandID: "GET",
	SetCommandID: "SET",
	DelCommandID: "DEL",
}

var commandNames = map[string]int{
	"GET": GetCommandID,
	"SET": SetCommandID,
	"DEL": DelCommandID,
}

var argumentsNumber = map[int]int{
	GetCommandID: 1,
	SetCommandID: 2,
	DelCommandID: 1,
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
