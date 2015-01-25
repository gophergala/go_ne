package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gophergala/go_ne/core"
	"github.com/mgutz/ansi"
)

var (
	taskName   = flag.String("task", "", "which task to run")
	configFile = flag.String("config", ".kiss.yml", "path to config file")
)

func main() {
	flag.Parse()

	config, err := core.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = config.Load(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	runner, err := NewRemoteRunner()
	if err != nil {
		fail(err)
	}

	if len(*taskName) == 0 {
		fail(errors.New("Please select the task to execute by passing the `-task=name` flag"))
	}

	fmt.Println(ansi.Color(fmt.Sprintf("Executing `%v`", *taskName), "green"))
	err = core.RunTask(runner, config, *taskName)
	if err != nil {
		fail(err)
	} else {
		fmt.Println(ansi.Color("Tasks completed successfully", "green"))
	}
}

func fail(err error) {
	fmt.Println(ansi.Color(fmt.Sprintf("Task failed: %v", err), "red"))
	os.Exit(1)
}
