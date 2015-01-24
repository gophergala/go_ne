package plugin

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

var host = flag.String("host", "localhost", "host for plugin server")
var port = flag.String("port", "1234", "port for plugin server")
var server = rpc.NewServer()

// BUG(Tobscher) These arguments are passed from
// the main process.
type Args struct {
	// A, B int
	Environment []string
	Options     []string
}

type Response struct {
	Name string
	Args []string
}

type Responder interface {
	Execute(args Args, reply *Response) error
}

func NewResponse(name string, args []string) Response {
	return Response{
		Name: name,
		Args: args,
	}
}

func Register(r Responder) {
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
