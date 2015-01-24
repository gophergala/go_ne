package core

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/gophergala/go_ne/plugins/core"
)

var pluginPrefix = "plugin"
var loadedPlugins = make(map[string]*Plugin)
var startPort = 8000

type Plugin struct {
	information *PluginInformation
	client      *rpc.Client
}

type PluginInformation struct {
	Host string
	Port string
	Cmd  *exec.Cmd
}

func (p *PluginInformation) Address() string {
	return fmt.Sprintf("%v:%v", p.Host, p.Port)
}

func StartPlugin(name string) *Plugin {
	command := fmt.Sprintf("%v-%v", pluginPrefix, name)
	host := "localhost"
	port := nextAvailblePort()

	log.Printf("Starting plugin `%v` on port %v\n", name, port)

	cmd := exec.Command(command,
		fmt.Sprintf("-host=%v", host),
		fmt.Sprintf("-port=%v", port),
	)
	cmd.Stdout = os.Stdout

	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	info := &PluginInformation{
		Host: host,
		Port: port,
		Cmd:  cmd,
	}

	var conn net.Conn
	for i := 1; i <= 5; i++ {
		log.Printf("Attempt %v to connect to plugin...", i)

		conn, err = net.Dial("tcp", info.Address())
		if err != nil {
			log.Print("FAILED")
			time.Sleep(100 * time.Millisecond)
			continue

			if i == 5 {
				return nil
			}
		}

		log.Print("OK")

		break
	}

	client := jsonrpc.NewClient(conn)

	plugin := &Plugin{
		information: info,
		client:      client,
	}

	loadedPlugins[name] = plugin

	return plugin
}

func GetPlugin(name string) (*Plugin, error) {
	var val *Plugin
	var ok bool

	if val, ok = loadedPlugins[name]; !ok {
		val = StartPlugin(name)
	}
	return val, nil
}

func (p *Plugin) GetCommand(args plugin.Args) (*Command, error) {
	// Pass in environment
	var reply plugin.Response
	err := p.client.Call("Command.Execute", args, &reply)
	if err != nil {
		return nil, err
	}

	return &Command{
		name: reply.Name,
		args: reply.Args,
	}, nil
}

func nextAvailblePort() string {
	startPort++
	return strconv.Itoa(startPort)
}

// BUG(Tobscher) Send signal to gracefully shutdown the plugin
func StopAllPlugins() {
	for k, v := range loadedPlugins {
		log.Printf("Stopping plugin: %v\n", k)
		v.information.Cmd.Process.Kill()
	}
}
