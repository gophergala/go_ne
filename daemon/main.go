package main

import (
	"log"

	"github.com/gophergala/go_ne/core"
)

// BUG(Tobscher) use command line arguments to perform correct task
func main() {
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
				task, _ = core.NewPlugin(*s.Plugin)
			} else {
				task, _ = core.NewCommand(*s.Command, s.Args)
			}
			runner.Run(task)
		}
	}

	core.StopAllPlugins()
}
