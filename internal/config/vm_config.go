package config

import (
	"os"

	"github.com/kebairia/kvmcli/internal/logger"
	"gopkg.in/yaml.v3"
)

type Resource interface {
	Validate()
}

type VirtualMachine struct {
	Kind     string   `yaml:"kind"`
	Metadata Metadata `yaml:"metadata"`
	Spec     Spec     `yaml:"spec"`
}
type Metadata struct {
	Name   string            `yaml:"name"`
	Labels map[string]string `yaml:"labels"`
}
type Spec struct {
	CPU       int     `yaml:"cpu"`
	Memory    string  `yaml:"memory"`
	Image     string  `yaml:"image"`
	Disk      Disk    `yaml:"disk"`
	Network   Network `yaml:"network"`
	Autostart bool    `yaml:"autostart"`
}

type Disk struct {
	Size string `yaml:"size"`
	Path string `yaml:"path"`
}
type Network struct {
	Name       string `yaml:"name"`
	MacAddress string `yaml:"macAddress"`
}

func LoadConfig[T any](configPath string) (T, error) {
	var resource T
	data, err := os.ReadFile(configPath)
	if err != nil {
		logger.Log.Errorf("failed to read file: %v", err)
	}

	if err := yaml.Unmarshal(data, &resource); err != nil {
		logger.Log.Errorf("failed to parse YAML: %v", err)
	}

	return resource, nil
}
