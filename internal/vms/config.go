package vms

// VirtualMachine represents a VM specification (loaded from YAML) and its runtime dependencies.
type VirtualMachineConfig struct {
	// manifest fields (populated by YAML unmarshal)
	APIVersion string   `yaml:"apiVersion"`
	Kind       string   `yaml:"kind"`
	Metadata   Metadata `yaml:"metadata"`
	Spec       Spec     `yaml:"spec"`
}

// Metadata contains identifying information for the VM.
type Metadata struct {
	Name      string            `yaml:"name"`
	Namespace string            `yaml:"namespace"`
	Labels    map[string]string `yaml:"labels"`
	Store     string            `yaml:"store"`
}

// Spec holds the VM’s desired configuration.
type Spec struct {
	Image     string  `yaml:"image"`
	CPU       int     `yaml:"cpu"`
	Memory    int     `yaml:"memory"`
	Disk      Disk    `yaml:"disk"`
	Network   Network `yaml:"network"`
	Autostart bool    `yaml:"autostart"`
}

// Disk describes the VM’s disk configuration.
type Disk struct {
	Size string `yaml:"size"`
	Path string `yaml:"path"`
}

// Network describes the VM’s network configuration.
type Network struct {
	Name       string `yaml:"name"`
	IP         string `yaml:"ip"`
	MacAddress string `yaml:"mac"`
}
