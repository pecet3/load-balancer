package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

type ConfigServer struct {
	URL           string `yaml:"URL"`
	StatusURL     string `yaml:"statusURL"`
	IsWsCandidate bool   `yaml:"isWsCandidate"`
}

type Config struct {
	Port           int             `yaml:"port"`
	StatusInterval int             `yaml:"statusInterval"`
	Servers        []*ConfigServer `yaml:"servers"`
	LogBuffSize    int             `yaml:"loggerBufferSize"`
}

func GetConfig() (*Config, error) {
	data, err := os.ReadFile("cfg/config.yaml")
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
	if cfg.LogBuffSize <= 0 {
		cfg.LogBuffSize = 100
	}
	return cfg, nil
}
