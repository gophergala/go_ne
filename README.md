# go_ne

## Description

go_ne is an automation tool which allows you to execute arbitrary tasks. It is intended
to make deployment easier.

go_ne can be used in two different ways:
* Remote execution of scripts via SSH
* Execution of scripts via the web interface

## Plugins

### How it works

We make use of a communication concept called [RPC](http://en.wikipedia.org/wiki/Remote_procedure_call). RPC allows
us to communicate with plugins in an elegant way.

### Write your own plugin

You can easily write your own plugin by using our plugin framework. Here is an example:

```go
package main

import (
	"github.com/gophergala/go_ne/plugins/core"
	"github.com/gophergala/go_ne/plugins/shared"
)

type Command struct {
}

func (t *Command) Execute(args shared.Args, reply *shared.Response) error {
	*reply = shared.NewResponse("env")

	return nil
}

func NewEnvCommand() *Command {
	return new(Command)
}

func main() {
	plugin.Register(NewEnvCommand())
	plugin.Serve()
}
```

The example above defines a plugin which runs the `env` command on your server.

## License

MIT License. Copyright 2015 James Rutherford & Tobias Haar.
