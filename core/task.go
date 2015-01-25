package core

import (
	"os"

	"github.com/gophergala/go_ne/plugins/shared"

	"errors"
)

type Task interface {
	Name() string
	Args() []string
}

func RunAll(runner Runner, config *Config) error {
	defer StopAllPlugins()
	defer runner.Close()

	for _, t := range config.Tasks {
		for _, s := range t.Steps {
			err := RunStep(runner, &s)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func RunTask(runner Runner, config *Config, taskName string) error {
	defer StopAllPlugins()
	defer runner.Close()

	task, ok := config.Tasks[taskName]
	if !ok {
		return errors.New("No task exists with that name")
	}

	for _, s := range task.Steps {
		err := RunStep(runner, &s)
		if err != nil {
			return err
		}
	}

	return nil
}

func RunStep(runner Runner, s *ConfigStep) error {
	var commands []*Command
	var err error

	if s.Plugin != nil {
		// Load plugin
		p, err := GetPlugin(*s.Plugin)
		if err != nil {
			return err
		}

		pluginArgs := shared.Args{
			Environment: os.Environ(),
			Args:        s.Args,
			Options:     s.Options,
		}

		commands, err = p.GetCommands(pluginArgs)
		if err != nil {
			return err
		}
	} else {
		// Run arbitrary command
		command, err := NewCommand(*s.Command, s.Args)
		if err != nil {
			return err
		}

		commands = append(commands, command)
	}

	for _, c := range commands {
		err = runner.Run(c)
		if err != nil {
			return err
		}
	}

	return nil
}
