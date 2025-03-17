package utils

import "github.com/kebairia/kvmcli/internal/config"

func GenerateYAMLTemplate() {
	sampleConfig := config.VMConfig{
		Version: "1.0",
		VMs: map[string]config.VM{
			"vm01": {
				CPU:    1,
				Memory: 1,
			},
		},
	}
	_ = sampleConfig
}

func GenerateTomlTemplate() {
}
