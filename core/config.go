package core

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type ConfigStep struct {
	Plugin  *string
	Command *string
	Options map[string]interface{}
	Args    []string
}

type ConfigTask struct {
	Steps []ConfigStep
}

type Config struct {
	Tasks map[string]ConfigTask
}

func NewConfig() (*Config, error) {
	c := Config{}
	return &c, nil
}

func (c *Config) Load(filepath string) error {
	rawYaml, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}

	yaml.Unmarshal(rawYaml, &c)
	if err != nil {
		return err
	}

	return nil
}
