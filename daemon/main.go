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
	for _, step := range config.Tasks[`sayhello`].Steps {
		task, _ := core.NewCommand(step.Command, step.Args)
		runner.Run(task)
	}

	for _, step := range config.Tasks[`deploy`].Steps {
		task, _ := core.NewCommand(step.Command, step.Args)
		runner.Run(task)
	}

	core.StopAllPlugins()
}
