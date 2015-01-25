package main

import (
	"fmt"
	"strings"

	"github.com/gophergala/go_ne/plugins/core"
	"github.com/gophergala/go_ne/plugins/shared"
)

type Command struct {
}

func (t *Command) Execute(args shared.Args, reply *shared.Response) error {
	cmd := []string{
		"ping",
	}

	if count, ok := args.Options["count"]; ok {
		c := count.(float64)
		cmd = append(cmd, fmt.Sprintf("-c %v", c))
	}

	cmd = append(cmd, shared.ExtractString(args.Options["url"]))
	*reply = shared.NewResponse(shared.NewCommand(strings.Join(cmd, " ")))

	return nil
}

func NewEnvCommand() *Command {
	return new(Command)
}

func main() {
	plugin.Register(NewEnvCommand())
	plugin.Serve()
}
