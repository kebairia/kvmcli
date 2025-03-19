package config

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type VirtualMachine struct {
	ApiVersion string   `yaml:"apiVersion"`
	Kind       string   `yaml:"kind"`
	Metadata   Metadata `yaml:"metadata"`
	Spec       Spec     `yaml:"spec"`
}
type Metadata struct {
	Name      string            `yaml:"name"`
	Namespace string            `yaml:"namespace"`
	Labels    map[string]string `yaml:"labels"`
}
type Spec struct {
	CPU       int     `yaml:"cpu"`
	Memory    int     `yaml:"memory"`
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

func LoadConfig(configPath string) ([]VirtualMachine, error) {
	// return error if you failed to read/open file
	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	defer file.Close()

	var vms []VirtualMachine
	decoder := yaml.NewDecoder(file)
	for {
		var vm VirtualMachine
		if err := decoder.Decode(&vm); err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("failed to decode YAML: %w", err)
		}
		vms = append(vms, vm)
	}

	return vms, nil
}
