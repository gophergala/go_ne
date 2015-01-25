package main

import (
	"github.com/gophergala/go_ne/plugins/core"
	"github.com/gophergala/go_ne/plugins/shared"
)

type Command struct {
}

func (t *Command) Execute(args shared.Args, reply *shared.Response) error {
	*reply = shared.NewResponse(shared.NewCommand("whoami"))

	return nil
}

func NewCommand() *Command {
	return new(Command)
}

func main() {
	plugin.Register(NewCommand())
	plugin.Serve()
}
