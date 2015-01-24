package shared

type Args struct {
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
