package plugin

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
)

// BUG(Tobscher) These arguments are passed from
// the main process.
type Args struct {
	A, B int
}

type Commander interface {
}

func Register(c Commander) {
	rpc.Register(c)
}

func Serve() {
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", getAddress())
	if e != nil {
		log.Fatal("listen error:", e)
	}
	http.Serve(l, nil)
}

// BUG(Tobscher) This should support a dynamic port
func getAddress() string {
	return "127.0.0.1:1234"
}
