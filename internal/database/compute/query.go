package compute

type Query struct {
	CommandId int
	Arguments []string
}

func NewQuery(commandId int, arguments []string) Query {
	return Query{
		CommandId: commandId,
		Arguments: arguments,
	}
}

func (query *Query) GetCommandId() int {
	return query.CommandId
}

func (query *Query) GetArguments() []string {
	return query.Arguments
}
