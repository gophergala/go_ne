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

func (l *Local) Run(task core.Task) error {
	log.Printf("Running task: %v %v\n", task.Name(), strings.Join(task.Args(), " "))

	cmd := exec.Command(task.Name(), task.Args()...)
	var stdOut bytes.Buffer
	var stdErr bytes.Buffer
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	err := cmd.Run()
	if err != nil {
		fmt.Println(stdErr.String())
		return err
	}

	fmt.Println(stdOut.String())
	return nil
}

func (l *Local) Close() {

}
