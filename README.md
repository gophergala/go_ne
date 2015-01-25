# go_ne

## Description

go_ne is an automation tool which allows you to execute arbitrary tasks. It is intended
to make deployment easier.

go_ne can be used in two different ways:
* Remote execution of scripts via SSH
* Execution of scripts via the web interface

This project has been developed during [Gopher Gala 2015](http://gophergala.com/).

## Deploy via the web interface

Describe how to deploy via the web interface...

## Deploy to a remote system

Download the binary from GitHub:

```
$ go get github.com/gophergala/go_ne/cli
```

You can start your deployment by running the following command:

```
$ cli -task=deploy
```

### Options

#### -host

```
$ cli -host=www.example.org -task=deploy
```

Defines the remote host to connect to via SSH.

#### -port

```
$ cli -port=23
```

Defines the port which is used to connect via SSH. Default: 22

#### -username

```
$ cli -username=deployer -task=deploy
```

Defines the username to use to connect via SSH.

#### -password

```
$ cli -password=secret
```

Defines the password to be used to connect via SSH.

#### -key

```
$ cli -key=$HOME/.ssh/id_rsa
```

Defines the key file to be used to connect via SSH.

## Plugins

### How it works

We make use of a communication concept called [RPC](http://en.wikipedia.org/wiki/Remote_procedure_call). RPC allows
us to communicate with plugins in an elegant way.

NOTE: Your plugin has to start with the prefix `plugin-` in order to be revealed by the plugin framework.

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

## Contributing

1. Fork it ( https://github.com/gophergala/go_ne/fork )
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create a new Pull Request

## License

MIT License. Copyright 2015 James Rutherford & Tobias Haar.
