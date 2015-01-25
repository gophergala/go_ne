package main

import (
	"strings"

	"github.com/gophergala/go_ne/plugins/core"
	"github.com/gophergala/go_ne/plugins/shared"
)

type Command struct {
}

func (t *Command) Execute(args shared.Args, reply *shared.Response) error {
	var commands []shared.Command

	update := shared.ExtractBool(args.Options["update"])
	if update {
		commands = append(commands, updateCommand())
	}

	commands = append(commands, installCommand(args))

	*reply = shared.NewResponse(commands...)

	return nil
}

func NewCommand() *Command {
	return new(Command)
}

func updateCommand() shared.Command {
	cmd := []string{
		"sudo",
		"apt-get",
		"update",
		"-y",
	}

	return shared.NewCommand(strings.Join(cmd, " "))
}

func installCommand(args shared.Args) shared.Command {
	packages := shared.ExtractOptions(args.Options["packages"])

	cmd := []string{
		"sudo",
		"apt-get",
		"install",
		"-y",
		strings.Join(packages, " "),
	}

	return shared.NewCommand(strings.Join(cmd, " "))
}

func main() {
	plugin.Register(NewCommand())
	plugin.Serve()
}
