package main

/* EXPERIMENTAL - This will likely do bad things with some inputs until we properly escape them!
*/

import (
	"github.com/gophergala/go_ne/plugins/core"
	"github.com/gophergala/go_ne/plugins/shared"
	"strings"
	"errors"
)

type Command struct {
}

func (t *Command) Execute(args shared.Args, reply *shared.Response) error {
	optionUser, ok := args.Options["user"]; if !ok {
		return errors.New("Mandatory argument omitted: user")
	}

	optionEndpoint, ok := args.Options["endpoint"]; if !ok {
		return errors.New("Mandatory argument omitted: endpoints")
	}

	optionFrom, ok := args.Options["from"]; if !ok {
		return errors.New("Mandatory argument omitted: from")
	}

	optionTo, ok := args.Options["to"]; if !ok {
		return errors.New("Mandatory argument omitted: to")
	}

	optionSubject, ok := args.Options["subject"]; if !ok {
		return errors.New("Mandatory argument omitted: subject")
	}

	optionText, ok := args.Options["text"]; if !ok {
		return errors.New("Mandatory argument omitted: text")
	}

	cmd := []string{
		"curl",
		"--user '" + optionUser.(string) + "'",
		optionEndpoint.(string),
		"-F from='" + optionFrom.(string) + "'",
		"-F to='" + optionTo.(string) + "'",
		"-F subject='" + optionSubject.(string) + "'",
		"-F text='" + optionText.(string) + "'",
	}

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
