package compute

const (
	GetCommandID = 1
	SetCommandID = 2
	DelCommandID = 3
)

var commandIDs = map[int]string{
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

func commandToID(command string) (int, bool) {
	commandID, ok := commandNames[command]
	return commandID, ok
}

func IDToCommand(commandID int) (string, bool) {
	commandName, ok := commandIDs[commandID]
	return commandName, ok
}

func getArgumentsNumber(commandID int) (int, bool) {
	number, ok := argumentsNumber[commandID]
	return number, ok
}
