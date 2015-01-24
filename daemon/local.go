package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/gophergala/go_ne/core"
)

type Local struct {
}

func NewLocalRunner() (*Local, error) {
	return &Local{}, nil
}

func (l *Local) Run(task core.Task) {
	log.Printf("Running task: %v %v\n", task.Name(), strings.Join(task.Args(), " "))

	cmd := exec.Command(task.Name(), task.Args()...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(out.String())
}
