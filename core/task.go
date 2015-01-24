package core

type Task struct {
	Command string
	Args    []string
}

func NewTask(command string, args []string) (*Task, error) {
	return &Task{
		Command: command,
		Args:    args,
	}, nil
}
