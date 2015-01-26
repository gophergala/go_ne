# GoKiss (by team go_ne)

## Inspiration

James runs a few hobby projects on a VM.  
Tobias helps manage large-scale infrastructure for a multinational company.   

We wanted to meet somewhere in the middle to create an easily set up, easy to use, highly configurable tool to help manage our infrastructures - small and large scale.  

## Overview

Kiss is a plugin-based automation tool which allows you to execute arbitrary tasks on a remote machine (or even locally). It is intended to make deployments easier.

You can interact with Kiss in three different ways:
* Via a Command Line Interface  
* Via a web interface  
* Using webhooks or periodic triggers  

It is able to:
* Remotely execute scripts via SSH
* Locally execute scripts
(tested on Ubuntu Linux)

The execution of your tasks is driven by a YAML file. Each task consists of a number of steps. Tasks can be executed against a list of server groups (e.g. run task `mysqldump` against all servers in group `db`).

Having your deployment scripts in a local place gives you and your team a number of benefits:
* Team members don't have to remember how certain tasks have to be executed
* You can easily add additional tasks by modifying a single configuration file
* The configuration can be in version control (tasks don't get lost)

The plugin-based architecture allows you to control more complex deployment tasks.

This project has been developed during [Gopher Gala 2015](http://gophergala.com/).

## The CLI

![Preview](https://i.imgflip.com/gt47e.gif "kiss preview")

For the CLI, the YAML file can be placed in your project's root folder (`.kiss.yml`).

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
    - host: www.example.org
      username: "root"
      key_path: "/path/to/private_key"
  db:
    - host: localhost
      username: "vagrant"
      password: "vagrant"
      port: 2222
    - host: db1.example.org
      username: "db"
      key_path: "/path/to/private_key"
    - host: db2.example.org
      username: "db"
      key_path: "/path/to/private_key"

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

(If the server is to operate locally, use `run_locally: true` instead of ssh username/password/key_path).  

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

## The Web Interface

![The Web Interface](https://raw.githubusercontent.com/gophergala/go_ne/master/githubdocs/gokiss-web.jpg)

[Note: the web interface is currently less complete than the CLI and has some messaging issues]

To start the daemon/web interface, build and run it:

```
$ cd go_ne/daemon/
```

```
$ go build -o daemon.exe
```

```
$ daemon.exe
```

The configuration file is the same format as for the CLI, but since it's running as a web facing daemon, we can set up triggers (currently GitHub Webhooks and periodic jobs are supported).


## Configuration

An example configuration looks like this:

```yaml
servergroups:
  web:
    - host: localhost
      username: "vagrant"
      password: "vagrant"
      port: 2222
    - host: www.example.org
      username: "root"
      key_path: "/path/to/private_key"
  db:
    - host: localhost
      username: "vagrant"
      password: "vagrant"
      port: 2222
    - host: db1.example.org
      username: "db"
      key_path: "/path/to/private_key"
    - host: db2.example.org
      username: "db"
      key_path: "/path/to/private_key"

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

### Additional Daemon Configuration 

The daemon also supports a 'triggers' block.

The following responds to the GitHub webhook on `/gokiss/triggers/gitpush` by running the `deploy` task on the `web` servergroup.

It also runs a periodic job `every day` (86400 seconds), running the `mysqldump` task on the `db` servergroup.

```yaml
triggers:
  githubrepopush:
    type: github-webhook
    endpoint: gitpush
    secret: mygithubsecrethere
    servergroup: web
    task: deploy
  monitor:
    type: periodic
    period: 86400
    servergroup: db
    task: mysqldump
```

The web setup is configurable. This serves from `domain:port/gokiss/` and configures one user, `gokiss` with password `default`.

```yaml
interfaces:
  web:
    settings:
      folder: /gokiss
      port: 20000
      sessionsecret: e53cr3t5h3re
    users:
      - username: gokiss
        password: default
```

### Options

#### -config

You can also specify the location of the config file for the daemon (default is /config/test-tasks.yaml):

```
$ daemon.exe -config=.kiss.yml
```

## Plugins

### How plugins work

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
* Daemon has problems running remote tasks [partly patched in a branch]
* The Web interface should run under SSL
* The configuration file could take a variable block
* Results from tasks could be assigned to variables and feed into other tasks, or branch the flow (so periodic server monitoring => alert emails could then be possible)

## Contributing

We hope you find this project useful. Please help us build upon it!

1. Fork it ( https://github.com/gophergala/go_ne/fork )
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create a new Pull Request


![Goofy Gopher](https://raw.githubusercontent.com/gophergala/go_ne/master/githubdocs/gopher.jpg)

## License

MIT License. Copyright 2015 [James Rutherford](https://twitter.com/jtruk "Twitter - JTRUK") & [Tobias Haar](https://twitter.com/tobscher "Twitter - Tobscher").
