package main

import "github.com/gophergala/go_ne/core"

func main() {
	task, _ := core.NewTask("ls", []string{
		"-la",
	})
	runner, _ := NewLocalRunner()

	runner.Run(task)
}
