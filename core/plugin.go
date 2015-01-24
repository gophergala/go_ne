package core

import (
	"bytes"
	"fmt"
	"log"
	"net/rpc"
	"os/exec"
	"strconv"
	"time"

	"github.com/gophergala/go_ne/plugins/core"
)

var pluginPrefix = "plugin"
var loadedPlugins = make(map[string]*PluginInformation)
var startPort = 8000

type Plugin struct {
	Command
}

type PluginInformation struct {
	Host string
	Port string
	Cmd  *exec.Cmd
}

func (p *PluginInformation) Address() string {
	return fmt.Sprintf("%v:%v", p.Host, p.Port)
}

func StartPlugin(name string) *PluginInformation {
	command := fmt.Sprintf("%v-%v", pluginPrefix, name)
	host := "localhost"
	port := nextAvailblePort()

	log.Printf("Starting plugin `%v` on port %v\n", name, port)

	cmd := exec.Command(command,
		fmt.Sprintf("-host=%v", host),
		fmt.Sprintf("-port=%v", port),
	)
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	info := &PluginInformation{
		Host: host,
		Port: port,
		Cmd:  cmd,
	}

	loadedPlugins[name] = info

	return info
}

func NewPlugin(name string) (Task, error) {
	var val *PluginInformation
	var ok bool
	var client *rpc.Client
	var err error

	if val, ok = loadedPlugins[name]; !ok {
		val = StartPlugin(name)
	}

	for i := 1; i <= 5; i++ {
		log.Printf("Attempt %v to connect to plugin...", i)

		client, err = rpc.DialHTTP("tcp", val.Address())
		if err != nil {
			log.Print("FAILED")
			time.Sleep(100 * time.Millisecond)
			continue

			if i == 5 {
				return nil, err
			}
		}

		log.Print("OK")

		break
	}

	// Pass in environment
	args := &plugin.Args{7, 8}
	var reply plugin.Response
	err = client.Call("Command.Execute", args, &reply)
	if err != nil {
		return nil, err
	}

	plugin := Plugin{
		Command: Command{
			name: reply.Name,
			args: reply.Args,
		},
	}

	return &plugin, nil
}

func nextAvailblePort() string {
	startPort++
	return strconv.Itoa(startPort)
}

// BUG(Tobscher) Send signal to gracefully shutdown the plugin
func StopAllPlugins() {
	for k, v := range loadedPlugins {
		log.Printf("Stopping plugin: %v\n", k)
		v.Cmd.Process.Kill()
	}
}
