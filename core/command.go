package core

type Command struct {
	name string
	args []string
}

func (c *Command) Name() string {
	return c.name
}

func (c *Command) Args() []string {
	return c.args
}

func NewCommand(name string, args []string) (*Command, error) {
	command := Command{
		name: name,
		args: args,
	}

	return &command, nil
}
