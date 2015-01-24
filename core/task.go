package core

type Task interface {
	Name() string
	Args() []string
}
