package main

import (
	"io"
	"os/exec"

	"github.com/gophergala/go_ne/core"
	"log"
	"strings"
)


const READ_BUFFER_SIZE = 1024


type Local struct {
	chStdOut   chan []byte
	chStdErr   chan []byte
}


func NewLocalRunner() (*Local, error) {
	lr := Local{
		chStdOut: make(chan []byte),
		chStdErr: make(chan []byte),
	}

	return &lr, nil
}


func (l *Local) Run(task core.Task) (error) {
	log.Printf("Running task: %v %v\n", task.Name(), strings.Join(task.Args(), " "))

	cmd := exec.Command(task.Name(), task.Args()...)
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
	listeningOut := true
	listeningErr := true
	for {
		if(listeningOut) {
			outBytes, err := stdOut.Read(bufferStdOut); if err != nil {
				if(err == io.EOF) {
					listeningOut = false
				} else {
					return err
				}
			}
			
			if(outBytes > 0) {
				l.chStdOut <- bufferStdOut[:outBytes]
			}
		}
		
		if(listeningErr) {
			errBytes, err := stdOut.Read(bufferStdErr); if err != nil {
				if(err == io.EOF) {
					listeningErr = false
				} else {
					return err
				}
			}

			if(errBytes > 0) {
				l.chStdErr <- bufferStdErr[:errBytes]
			}
		}
		
		if(!listeningOut && !listeningErr) {	// both complete
			break
		}
	}
	
	cmd.Wait()
	return nil
}

func (l *Local) Close() {

}
