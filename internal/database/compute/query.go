package compute

type Query struct {
	CommandID int
	Arguments []string
}

func NewQuery(commandID int, arguments []string) Query {
	return Query{
		CommandID: commandID,
		Arguments: arguments,
	}
}

func (query *Query) GetCommandID() int {
	return query.CommandID
}

func (query *Query) GetArguments() []string {
	return query.Arguments
}
