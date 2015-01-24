package plugin

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"

	"github.com/gophergala/go_ne/plugins/shared"
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

	log.Printf("Started plugin on `%v`\n", address)

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
