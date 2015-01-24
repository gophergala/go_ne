package main

import (
	"log"

	"github.com/gophergala/go_ne/core"
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

	core.RunAll(runner, config)
}
