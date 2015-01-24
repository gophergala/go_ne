package core

import (
	"fmt"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/gophergala/go_ne/plugins/shared"
	"github.com/mgutz/ansi"
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

func StartPlugin(name string) (*Plugin, error) {
	command := fmt.Sprintf("%v-%v", pluginPrefix, name)
	host := "localhost"
	port := nextAvailblePort()

	fmt.Println(ansi.Color(fmt.Sprintf("-- Starting plugin `%v` on port %v", name, port), "black+h"))

	cmd := exec.Command(command,
		fmt.Sprintf("-host=%v", host),
		fmt.Sprintf("-port=%v", port),
	)
	cmd.Stdout = os.Stdout

	err := cmd.Start()
	if err != nil {
		return nil, err
	}

	info := &PluginInformation{
		Host: host,
		Port: port,
		Cmd:  cmd,
	}

	var conn net.Conn
	for i := 1; i <= 5; i++ {
		fmt.Print(ansi.Color(fmt.Sprintf("-- Attempt %v to connect to plugin...", i), "black+h"))

		conn, err = net.Dial("tcp", info.Address())
		if err != nil {
			fmt.Println(ansi.Color("FAILED", "black+h"))
			time.Sleep(100 * time.Millisecond)
			continue

			if i == 5 {
				return nil, err
			}
		}

		fmt.Println(ansi.Color("OK", "black+h"))

		break
	}

	client := jsonrpc.NewClient(conn)

	plugin := &Plugin{
		information: info,
		client:      client,
	}

	loadedPlugins[name] = plugin

	return plugin, nil
}

func GetPlugin(name string) (*Plugin, error) {
	var val *Plugin
	var ok bool
	var err error

	val, ok = loadedPlugins[name]
	if !ok {
		val, err = StartPlugin(name)
		if err != nil {
			return nil, err
		}
	}
	return val, nil
}

func (p *Plugin) GetCommand(args shared.Args) (*Command, error) {
	// Pass in environment
	var reply shared.Response
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
		fmt.Println(ansi.Color(fmt.Sprintf("-- Stopping plugin: %v", k), "black+h"))
		v.information.Cmd.Process.Kill()
	}
}
