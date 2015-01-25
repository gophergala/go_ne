package plugin

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"

	"github.com/gophergala/go_ne/plugins/shared"
	"github.com/mgutz/ansi"
)

var host = flag.String("host", "localhost", "host for plugin server")
var port = flag.String("port", "1234", "port for plugin server")
var server = rpc.NewServer()

// Register registers the given responder to be used as a plugin.
// The plugin will be able to receive calls via RPC.
// Please make sure your struct is called `Command` as the remote call is `Command.Execute`
func Register(r shared.Responder) {
	server.Register(r)
}

// Serve Starts the server and listens to it.
func Serve() {
	flag.Parse()

	address := getAddress()

	server.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)
	l, e := net.Listen("tcp", address)
	if e != nil {
		log.Fatal("listen error:", e)
	}

	fmt.Println(ansi.Color(fmt.Sprintf("Started plugin on `%v`", address), "black+h"))

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go server.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}

func getAddress() string {
	return fmt.Sprintf("%v:%v", *host, *port)
}
