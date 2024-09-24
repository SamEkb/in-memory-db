package configuration

import (
	"os"

	"gopkg.in/yaml.v3"
)

const config = "config.yaml"

type Configuration struct {
	Network *NetworkConfig `yaml:"network"`
	Engine  *EngineConfig  `yaml:"engine"`
	Logging *LoggingConfig `yaml:"logging"`
}

type NetworkConfig struct {
	Address        string `yaml:"address"`
	MaxConnection  int    `yaml:"max_connections"`
	MaxMessageSize string `yaml:"max_message_size"`
	IdleTimeout    string `yaml:"idle_timeout"`
}

type EngineConfig struct {
	Type string `yaml:"type"`
}

type LoggingConfig struct {
	Level  string `yaml:"level"`
	Output string `yaml:"output"`
}

func NewConfiguration() (*Configuration, error) {
	yamlData, err := os.ReadFile(config)
	if err != nil {
		return &Configuration{}, err
	}

	var config *Configuration
	err = yaml.Unmarshal(yamlData, &config)
	if err != nil {
		return &Configuration{}, err
	}

	return config, nil
}
