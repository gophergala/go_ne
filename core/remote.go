package core

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"

	"github.com/mgutz/ansi"
)

// Remote describes a runner which runs task
// on a remote system via SSH.
type Remote struct {
	Client *ssh.Client
}

// NewRemoteRunner creates a new runner which runs
// tasks on a remote system.
//
// An SSH connection will be establishe.
func NewRemoteRunner(options ConfigServer) (*Remote, error) {
	flag.Parse()

	client, err := createClient(options.Username, options.Password, options.Host, options.Port, options.KeyPath)
	if err != nil {
		return nil, err
	}

	return &Remote{
		Client: client,
	}, nil
}

// Run runs the given task on the remote system
func (r *Remote) Run(task Task) error {
	session, err := r.Client.NewSession()
	if err != nil {
		return errors.New("Failed to create session: " + err.Error())
	}
	defer session.Close()

	fmt.Println(ansi.Color(fmt.Sprintf("Executing `%v %v`", task.Name(), strings.Join(task.Args(), " ")), "green"))

	cmd := fmt.Sprintf("%v %v", task.Name(), strings.Join(task.Args(), " "))

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	if err := session.Start(cmd); err != nil {
		return err
	}

	return session.Wait()
}

// Close closes the SSH connection to the remote system
func (r *Remote) Close() {
	r.Client.Close()
}

func createClient(username, password, host, port, key string) (*ssh.Client, error) {
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

	if len(port) == 0 {
		port = "22"
	}

	remoteServer := fmt.Sprintf("%v:%v", host, port)

	fmt.Println(ansi.Color(fmt.Sprintf("Connecting to %v@%v", username, remoteServer), "green"))
	client, err := ssh.Dial("tcp", remoteServer, config)
	if err != nil {
		return nil, err
	}

	return client, nil
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
