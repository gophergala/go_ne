package main

import "github.com/gophergala/go_ne/plugins/core"

type Command struct {
}

func (t *Command) Execute(args plugin.Args, reply *plugin.Response) error {
	*reply = plugin.NewResponse("ping", args.Options)

	return nil
}

func NewEnvCommand() *Command {
	return new(Command)
}

func main() {
	plugin.Register(NewEnvCommand())
	plugin.Serve()
}
