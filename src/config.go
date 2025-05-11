package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

type ConfigServer struct {
	URL       string `yaml:"url"`
	StatusURL string `yaml:"statusURL"`
}

type Config struct {
	Port           int             `yaml:"port"`
	StatusInterval int             `yaml:"statusInterval"`
	Servers        []*ConfigServer `yaml:"servers"`
}

func GetConfig() (*Config, error) {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		data, err = os.ReadFile("config.yml")
		if err != nil {
			return nil, err
		}
	}
	cfg := &Config{}
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
