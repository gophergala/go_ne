package core

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
//	"log"
)


type ConfigStep struct {
	Command   string
	Args      []string
}


type ConfigTask struct {
	Steps   []ConfigStep
}


type Config struct {
	Tasks   map[string]ConfigTask
}


func NewConfig() (*Config, error) {
	c := Config{}
	return &c, nil
}


func (c* Config) Load(filepath string) (error) {
	rawYaml, err := ioutil.ReadFile(filepath); if err != nil {
		return err
	}
	
	yaml.Unmarshal(rawYaml, &c); if err != nil {
		return err
	}
	
	return nil
}
