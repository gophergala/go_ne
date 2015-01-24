package core

import (
	"fmt"
	"net/rpc"

	"github.com/gophergala/go_ne/plugins/core"
)

var loadedPlugins = make(map[string]*Plugin)

type Plugin struct {
	Command
}

func NewPlugin() (Task, error) {
	client, err := rpc.DialHTTP("tcp", "127.0.01:1234")
	if err != nil {
		return nil, err
	}

	// Pass in environment
	args := plugin.Args{7, 8}
	var reply string
	err = client.Call("Command.Execute", args, &reply)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Command: %d - %d = %d\n", args.A, args.B, reply)

	plugin := Plugin{
		Command: Command{
			name: reply,
			args: []string{},
		},
	}

	return &plugin, nil
}
