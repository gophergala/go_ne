package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gophergala/go_ne/core"
	"github.com/mgutz/ansi"
)

func main() {
	config, err := core.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = config.Load("config/test-tasks.yaml")
	if err != nil {
		log.Fatal(err)
	}

	runner, err := NewRemoteRunner()
	if err != nil {
		fail(err)
	}

	err = core.RunAll(runner, config)
	if err != nil {
		fail(err)
	} else {
		fmt.Println(ansi.Color("Tasks completed successfully", "green"))
	}
}

func fail(err error) {
	fmt.Println(ansi.Color(fmt.Sprintf("Tasks failed: %v", err), "red"))
	os.Exit(1)
}
