package plugin

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

var host = flag.String("host", "localhost", "host for plugin server")
var port = flag.String("port", "1234", "port for plugin server")

// BUG(Tobscher) These arguments are passed from
// the main process.
type Args struct {
	A, B int
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
	rpc.Register(r)
}

func Serve() {
	flag.Parse()

	address := getAddress()

	rpc.HandleHTTP()
	l, e := net.Listen("tcp", address)
	if e != nil {
		log.Fatal("listen error:", e)
	}

	log.Printf("Started plugin on `%v`\n", address)

	http.Serve(l, nil)
}

func getAddress() string {
	return fmt.Sprintf("%v:%v", *host, *port)
}
