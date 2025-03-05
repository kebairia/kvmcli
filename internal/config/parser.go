package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type VMConfig struct {
	Version string        `yaml:"version"`
	VMs     map[string]VM `yaml:"vms"`
}

type VM struct {
	Name      string  `yaml:"name"`
	CPU       int     `yaml:"cpu"`
	Memory    string  `yaml:"memory"`
	Image     string  `yaml:"image"`
	Disk      Disk    `yaml:"disk"`
	Network   Network `yaml:"network"`
	OS        Os      `yaml:"os"`
	Autostart bool    `yaml:"autostart"`
}

type Disk struct {
	Size string `yaml:"size"`
	Path string `yaml:"path"`
}

type Network struct {
	Type   string `yaml:"type"`
	Source string `yaml:"source"`
	MAC    string `yaml:"mac_address"`
}

type Os struct {
	Type string `yaml:"type"`
	ISO  string `yaml:"iso"` // Fixed
}

func LoadConfig(path string) (*VMConfig, error) {
	// Read the YAML file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	var config VMConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	return &config, nil
}
