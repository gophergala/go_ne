package core

import (
	"log"
	"os"

	"github.com/gophergala/go_ne/plugins/core"
)

type Task interface {
	Name() string
	Args() []string
}

func RunAll(runner Runner, config *Config) {
	var err error
	var task Task

	for _, t := range config.Tasks {
		for _, s := range t.Steps {
			if s.Plugin != nil {
				// Load plugin
				p, err := GetPlugin(*s.Plugin)
				if err != nil {
					log.Println(err)
					continue
				}

				pluginArgs := plugin.Args{
					Environment: os.Environ(),
					Options:     s.Args,
				}

				task, err = p.GetCommand(pluginArgs)
				if err != nil {
					log.Println(err)
					continue
				}
			} else {
				// Run arbitrary command
				task, err = NewCommand(*s.Command, s.Args)
				if err != nil {
					log.Println(err)
					continue
				}
			}

			runner.Run(task)
		}
	}

	StopAllPlugins()
}
