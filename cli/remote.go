package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"

	"github.com/gophergala/go_ne/core"
	"github.com/mgutz/ansi"
)

var username = flag.String("username", "", "username for remote server")
var password = flag.String("password", "", "password for remote server")
var key = flag.String("key", "", "path to private key")
var host = flag.String("host", "", "host for remote server")
var port = flag.String("port", "22", "ssh port")

type Remote struct {
	Client *ssh.Client
}

func NewRemoteRunner() (*Remote, error) {
	flag.Parse()

	client := createClient(*username, *password, *host, *port, *key)

	return &Remote{
		Client: client,
	}, nil
}

func (r *Remote) Run(task core.Task) error {
	session, err := r.Client.NewSession()
	if err != nil {
		panic("Failed to create session: " + err.Error())
	}
	defer session.Close()

	fmt.Println(ansi.Color(fmt.Sprintf("executing `%v %v`", task.Name(), strings.Join(task.Args(), " ")), "green"))

	cmd := fmt.Sprintf("%v %v", task.Name(), strings.Join(task.Args(), " "))

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	if err := session.Start(cmd); err != nil {
		log.Println(err)
		return err
	}

	session.Wait()

	return nil
}

func createClient(username, password, host, port, key string) *ssh.Client {
	authMethods := []ssh.AuthMethod{}

	if len(password) > 0 {
		authMethods = append(authMethods, ssh.Password(password))
	}

	if len(key) > 0 {
		priv, err := loadKey(key)
		if err != nil {
			log.Println(err)
		} else {
			signers, err := ssh.NewSignerFromKey(priv)
			if err != nil {
				log.Println(err)
			} else {
				authMethods = append(authMethods, ssh.PublicKeys(signers))
			}
		}
	}

	config := &ssh.ClientConfig{
		User: username,
		Auth: authMethods,
	}

	remoteServer := fmt.Sprintf("%v:%v", host, port)

	fmt.Println(ansi.Color(fmt.Sprintf("Connecting to %v@%v", username, remoteServer), "blue"))
	client, err := ssh.Dial("tcp", remoteServer, config)
	if err != nil {
		panic("Failed to dial: " + err.Error())
	}

	return client
}

func loadKey(file string) (interface{}, error) {
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	key, err := ssh.ParseRawPrivateKey(buf)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func (r *Remote) Close() {
	r.Client.Close()
}
