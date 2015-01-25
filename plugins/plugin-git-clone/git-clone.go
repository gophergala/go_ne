package main

import (
	"strings"

	"github.com/gophergala/go_ne/plugins/core"
	"github.com/gophergala/go_ne/plugins/shared"
)

type Command struct {
}

func (t *Command) Execute(args shared.Args, reply *shared.Response) error {
	directory := args.Options["directory"].(string)

	a := []string{
		"git",
		"clone",
		args.Options["repo"].(string),
		directory,
	}

	cmd := shared.NewCommand(strings.Join(a, " "))
	cmd.Unless(shared.DirectoryExists(directory))

	*reply = shared.NewResponse(cmd)

	return nil
}

func NewCommand() *Command {
	return new(Command)
}

func main() {
	plugin.Register(NewCommand())
	plugin.Serve()
}
