package config

import (
	"io"
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

//	func LoadConfig[T any](configPath string) (T, error) {
//		var resource T
//		data, err := os.ReadFile(configPath)
//		if err != nil {
//			logger.Log.Errorf("failed to read file: %v", err)
//		}
//
//		if err := yaml.Unmarshal(data, &resource); err != nil {
//			logger.Log.Errorf("failed to parse YAML: %v", err)
//		}
//
//		return resource, nil
//	}
func LoadConfig(configPath string) ([]VirtualMachine, error) {
	file, err := os.Open(configPath)
	if err != nil {
		logger.Log.Errorf("failed to read file: %v", err)
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
			logger.Log.Errorf("failed to decode YAML: %v", err)
		}
		vms = append(vms, vm)
	}
	// if err := yaml.Unmarshal(file, &resource); err != nil {
	// 	log.Printf("failed to parse YAML: %v", err)
	// }

	return vms, nil
}
