package main

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Upstream string `yaml:"upstream"`
	Routes   []struct {
		Path       string   `yaml:"path"`
		ClearCache []string `yaml:"clear_cache"`
	} `yaml:"routes"`
}

func NewConfig() *Config {
	configFile, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	config := &Config{}
	if err := yaml.Unmarshal(configFile, config); err != nil {
		log.Fatal(err)
	}

	return config
}
