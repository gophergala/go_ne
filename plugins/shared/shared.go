package shared

type Args struct {
	Environment []string
	Args        []string
	Options     map[string]interface{}
}

type Command struct {
	Name string
	Args []string
}

type Response struct {
	Commands []Command
}

type Responder interface {
	Execute(args Args, reply *Response) error
}

func NewCommand(name string, args ...string) Command {
	return Command{
		Name: name,
		Args: args,
	}
}

func NewResponse(commands ...Command) Response {
	return Response{Commands: commands}
}
