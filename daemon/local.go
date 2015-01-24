package main

import (
	"io"
	"os/exec"

	"github.com/gophergala/go_ne/core"
)


const READ_BUFFER_SIZE = 1024


type Local struct {
}


func NewLocalRunner() (*Local, error) {
	return &Local{}, nil
}


func (l *Local) Run(task *core.Task, cStdOut chan<- []byte) (error) {
	cmd := exec.Command(task.Command, task.Args...)
	stdOut, err := cmd.StdoutPipe(); if err != nil {
		return err
	}
	defer stdOut.Close()

	err = cmd.Start(); if err != nil {
		return err
	}
		
	buffer := make([]byte, READ_BUFFER_SIZE)
	for {
		bytes, err := stdOut.Read(buffer); if err != nil && err != io.EOF {
			return err		
		}
				
		cStdOut <- buffer[:bytes]

		if(bytes == 0) {
			break
		}
	}
	
	cmd.Wait()
	return nil
}
