package main

import "github.com/gophergala/go_ne/plugins/core"

type Command string

func (t *Command) Execute(args *plugin.Args, reply *string) error {
	*reply = "env"
	return nil
}

func NewCommand() *Command {
	return new(Command)
}

func main() {
	plugin.Register(NewCommand())
	plugin.Serve()
}
