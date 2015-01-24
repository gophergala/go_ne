package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"

	"github.com/gophergala/go_ne/core"
)

type Local struct {
}

func NewLocalRunner() (*Local, error) {
	return &Local{}, nil
}

func (l *Local) Run(task core.Task) {
	cmd := exec.Command(task.Name(), task.Args()...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run(); if err != nil {
		log.Fatal(err)
	}
	fmt.Println(out.String())
}
