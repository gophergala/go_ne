package main

import (
	"log"
	"os"

	"github.com/gophergala/go_ne/core"
	"github.com/gophergala/go_ne/plugins/core"
)

// BUG(Tobscher) use command line arguments to perform correct task
func main() {
	log.SetPrefix("[go-ne] ")

	config, err := core.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = config.Load("config/test-tasks.yaml")
	if err != nil {
		log.Fatal(err)
	}

	runner, _ := NewLocalRunner()

	var task core.Task
	for _, t := range config.Tasks {
		for _, s := range t.Steps {
			if s.Plugin != nil {
				// Load plugin
				p, err := core.GetPlugin(*s.Plugin)
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
				task, err = core.NewCommand(*s.Command, s.Args)
				if err != nil {
					log.Println(err)
					continue
				}
			}

			runner.Run(task)
		}
	}

	core.StopAllPlugins()
}
