package main

import (
	"log"
	"strings"

	"github.com/gophergala/go_ne/core"
)

func main() {
	tasks := []core.Task{}

	task, err := core.NewCommand("ls", []string{
		"-la",
	})
	if err != nil {
		log.Printf("Could not load task: %v\n", err)
	} else {
		tasks = append(tasks, task)
	}

	env, err := core.NewPlugin("env")
	if err != nil {
		log.Printf("Could not load plugin: %v\n", err)
	} else {
		tasks = append(tasks, env)
	}

	ping, err := core.NewPlugin("ping")
	if err != nil {
		log.Printf("Could not load plugin: %v\n", err)
	} else {
		tasks = append(tasks, ping)
	}

	runner, _ := NewLocalRunner()

	for _, v := range tasks {
		log.Printf("Running task `%v` with arguments `%v`\n", v.Name(), strings.Join(v.Args(), " "))
		runner.Run(v)
	}

	core.StopAllPlugins()
}
