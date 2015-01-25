package core

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)


type ConfigServer struct {
	Ip           string
	SshUser      string
	SshKey       string
	SshKeyPath   string
	Vars         map[string]string
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
	Type          string
	Period        uint
	Endpoint      string
	Secret        string
	ServerGroup   string
	Task          string
}
	
type Config struct {
	Vars            map[string]string
	ServerGroups    map[string][]ConfigServer
	Tasks           map[string]ConfigTask
	Triggers        map[string]ConfigEvent
	Interfaces      map[string]map[string]string
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
