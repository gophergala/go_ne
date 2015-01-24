package main

import (
	"github.com/gophergala/go_ne/core"
	"log"
)

func main() {
	config, err := core.NewConfig(); if err != nil {
		log.Fatal(err)
	}

	err = config.Load("config/test-tasks.yaml"); if err != nil {
		log.Fatal(err)
	}
	
	for _, step := range config.Tasks[`sayhello`].Steps {
		task, _ := core.NewTask(step.Command, step.Args)
		runner, _ := NewLocalRunner()

		runner.Run(task)
	}	

	for _, step := range config.Tasks[`deploy`].Steps {
		task, _ := core.NewTask(step.Command, step.Args)
		runner, _ := NewLocalRunner()

		runner.Run(task)
	}	
}
