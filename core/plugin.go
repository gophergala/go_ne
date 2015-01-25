package core

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"

	"github.com/gophergala/go_ne/plugins/shared"
	"github.com/mgutz/ansi"
)

// PluginCache stores loaded plugins
type PluginCache struct {
	sync.Mutex
	cache map[string]*Plugin
}

var pluginPrefix = "plugin"
var loadedPlugins = PluginCache{
	cache: make(map[string]*Plugin),
}
var maxAttempts = 5
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

	// Log to logfile
	cmd := exec.Command(command,
		fmt.Sprintf("-host=%v", host),
		fmt.Sprintf("-port=%v", port),
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

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
	for i := 1; i <= maxAttempts; i++ {
		fmt.Print(ansi.Color(fmt.Sprintf("-- Attempt %v to connect to plugin...", i), "black+h"))

		conn, err = net.Dial("tcp", info.Address())
		if err != nil {
			fmt.Println(ansi.Color("FAILED", "black+h"))
			time.Sleep(100 * time.Millisecond)

			if i == maxAttempts {
				cmd.Process.Kill()
				return nil, errors.New("Could not connect to plugin.")
			}

			continue
		}

		fmt.Println(ansi.Color("OK", "black+h"))

		break
	}

	client := jsonrpc.NewClient(conn)

	plugin := &Plugin{
		information: info,
		client:      client,
	}

	loadedPlugins.Lock()
	loadedPlugins.cache[name] = plugin
	loadedPlugins.Unlock()

	return plugin, nil
}

func GetPlugin(name string) (*Plugin, error) {
	var val *Plugin
	var ok bool
	var err error

	val, ok = loadedPlugins.cache[name]
	if !ok {
		val, err = StartPlugin(name)
		if err != nil {
			return nil, err
		}
	}
	return val, nil
}

func (p *Plugin) GetCommands(args shared.Args) ([]*Command, error) {
	var reply shared.Response
	var commands []*Command

	err := p.client.Call("Command.Execute", args, &reply)
	if err != nil {
		return nil, err
	}

	for _, value := range reply.Commands {
		command := &Command{
			name: value.Name,
			args: value.Args,
		}

		commands = append(commands, command)
	}

	return commands, nil
}

func nextAvailblePort() string {
	startPort++
	return strconv.Itoa(startPort)
}

// BUG(Tobscher) Send signal to gracefully shutdown the plugin
// BUG(Tobscher) Use lock
func StopAllPlugins() {
	loadedPlugins.Lock()
	defer loadedPlugins.Unlock()

	for k, v := range loadedPlugins.cache {
		fmt.Println(ansi.Color(fmt.Sprintf("-- Stopping plugin: %v", k), "black+h"))
		if err := v.information.Cmd.Process.Kill(); err != nil {
			log.Println(err)
		}
	}

	loadedPlugins.cache = make(map[string]*Plugin)
	startPort = 8000
}
