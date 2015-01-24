package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"strings"

	"golang.org/x/crypto/ssh"

	"github.com/gophergala/go_ne/core"
)

var username = flag.String("username", "username", "username for remote server")
var password = flag.String("password", "password", "password for remote server")
var host = flag.String("host", "localhost", "host for remote server")
var port = flag.String("port", "22", "ssh port")

type Remote struct {
	Session *ssh.Session
}

func NewRemoteRunner() (*Remote, error) {
	flag.Parse()

	return &Remote{
		Session: createSession(*username, *password, *host, *port),
	}, nil
}

func (r *Remote) Run(task core.Task) error {
	log.Printf("Running task: %v %v\n", task.Name(), strings.Join(task.Args(), " "))

	cmd := fmt.Sprintf("%v %v", task.Name(), strings.Join(task.Args(), " "))
	var stdOut bytes.Buffer
	var stdErr bytes.Buffer
	r.Session.Stdout = &stdOut
	r.Session.Stderr = &stdErr
	if err := r.Session.Run(cmd); err != nil {
		fmt.Println(stdErr.String())
		return err
		// panic("Failed to run: " + err.Error())
	}

	fmt.Println(stdOut.String())
	return nil
}

func createSession(username, password, host, port string) *ssh.Session {
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
	}

	remoteServer := fmt.Sprintf("%v:%v", host, port)

	log.Printf("Connecting to %v\n", remoteServer)
	client, err := ssh.Dial("tcp", remoteServer, config)
	if err != nil {
		panic("Failed to dial: " + err.Error())
	}

	session, err := client.NewSession()
	if err != nil {
		panic("Failed to create session: " + err.Error())
	}

	return session
}

func (r *Remote) Close() {
	r.Session.Close()
}
