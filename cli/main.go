package main

import (
	"log"

	"github.com/gophergala/go_ne/core"
)

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

	runner, _ := NewRemoteRunner()

	core.RunAll(runner, config)
}
