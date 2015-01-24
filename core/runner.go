package core

type Runner interface {
	Run(Task) error
	Close()
}
