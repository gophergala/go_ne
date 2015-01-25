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

func Register(r shared.Responder) {
	server.Register(r)
}

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
