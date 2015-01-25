package core

import (
	"errors"
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type ConfigServer struct {
	Host       string
	Username   string
	Password   string
	Port       string
	KeyPath    string `yaml:"key_path"`
	RunLocally bool   `yaml:"run_locally"`
	Vars       map[string]string
}

type ConfigStep struct {
	Plugin  *string
	Command *string
	Options map[string]interface{}
	Args    []string
}

type ConfigTask struct {
	Steps []ConfigStep
}

type ConfigEvent struct {
	Type        string
	Period      uint
	Endpoint    string
	Secret      string
	ServerGroup string
	Task        string
}

type ConfigWebUsers struct {
	Username string
	Password string
}

type ConfigWebInterface struct {
	Settings map[string]string
	Users    []ConfigWebUsers
}

type ConfigInterfaces struct {
	Web ConfigWebInterface
}

type Config struct {
	Vars         map[string]string
	ServerGroups map[string][]ConfigServer
	Tasks        map[string]ConfigTask
	Triggers     map[string]ConfigEvent
	Interfaces   ConfigInterfaces
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

func (c *Config) GetServerGroup(name string) ([]ConfigServer, error) {
	value, ok := c.ServerGroups[name]
	if !ok {
		return []ConfigServer{}, errors.New(fmt.Sprintf("Server group `%v` not found.", name))
	}

	return value, nil
}
