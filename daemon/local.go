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


func (l *Local) Run(task *core.Task, chStdOut chan<- []byte, chStdErr chan<- []byte) (error) {
	cmd := exec.Command(task.Command, task.Args...)
	stdOut, err := cmd.StdoutPipe(); if err != nil {
		return err
	}
	defer stdOut.Close()

	stdErr, err := cmd.StderrPipe(); if err != nil {
		return err
	}
	defer stdErr.Close()

	err = cmd.Start(); if err != nil {
		return err
	}
	
	// COULDDO: Might be slicker to spin up two async processes that communicate back
	bufferStdOut := make([]byte, READ_BUFFER_SIZE)
	bufferStdErr := make([]byte, READ_BUFFER_SIZE)
	for {
		outBytes, err := stdOut.Read(bufferStdOut); if err != nil && err != io.EOF {
			return err
		}
				
		errBytes, err := stdOut.Read(bufferStdErr); if err != nil && err != io.EOF {
			return err
		}

		chStdOut <- bufferStdOut[:outBytes]
		chStdErr <- bufferStdErr[:errBytes]

		if(outBytes == 0 && errBytes == 0) {
			break
		}
	}
	
	cmd.Wait()
	return nil
}
