# go_ne (codename kiss)

## Description

kiss is an automation tool which allows you to execute arbitrary tasks. It is intended
to make deployment easier.

kiss can be used in two different ways:
* Remote execution of scripts via SSH
* Execution of scripts via the web interface

This project has been developed during [Gopher Gala 2015](http://gophergala.com/).

## Deploy via the web interface

Describe how to deploy via the web interface...

## Deploy to a remote system

![Preview](https://i.imgflip.com/gsz38.gif "kiss preview")

Follow the example on [Tobscher/go-example-app](https://github.com/Tobscher/go-example-app#test-deployment-via-kiss) or do it manually:

Download the binary from GitHub:

```
$ go get github.com/gophergala/go_ne/kiss
```

You can start your deployment by running the following commands:

```
$ kiss -group=web -task=deploy
```

The commands above assume you have placed a `.kiss.yml` file in the root of your project. Here is an example
configuration file:

```yaml
servergroups:
  web:
    - host: localhost
      username: "vagrant"
      password: "vagrant"
      port: 2222

tasks:
  setup:
    steps:
      - plugin: apt-get
        options:
          update: true
          packages:
            - "git"
            - "golang"
            - "python-setuptools"
      - command: go version
      - command: easy_install supervisor
      - command: rm -rf example-app
      - plugin: git-clone
        options:
          repo: "https://github.com/your-fork/go-example-app.git"
          directory: "example-app"
      - command: cp example-app/supervisord.conf /etc/supervisord.conf
      - command: supervisord || echo "Looks like supervisord is already running"
  deploy:
    steps:
      - plugin: whoami
      - plugin: env
      - command: supervisorctl stop example-app
      - command: cd example-app && git pull
      - command: cd example-app && go test -v
      - command: cd example-app && go build -v
      - command: supervisorctl start example-app
      - command: curl http://your.server.org:8080/
  start:
    steps:
      - command: supervisorctl start example-app
  stop:
    steps:
      - command: supervisorctl stop example-app
```

TIP: You can use our test application to test the steps above: https://github.com/Tobscher/go-example-app

### Options

#### -group

```
$ kiss -group=web
```

Defines the group for which the task should run. This flag is mandatory.

#### -task

```
$ kiss -task=deploy
```

Defines the task that should run. This flag is mandatory.

#### -config

```
$ kiss -config=.kiss-staging.yml
```

Defines the config file which includes the task definition. Default .kiss.yml

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
	*reply = shared.NewResponse(shared.NewCommand("env"))

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

Please refer to the [plugins directory](https://github.com/gophergala/go_ne/tree/master/plugins) for more examples.

## Limitations

* Plugins will get a port assigned starting from 8000
* Plugins need to be prefixed wit plugin-, e.g. plugin-apt-get
* Some tasks require sudo

## Contributing

1. Fork it ( https://github.com/gophergala/go_ne/fork )
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create a new Pull Request

## License

MIT License. Copyright 2015 James Rutherford & Tobias Haar.
