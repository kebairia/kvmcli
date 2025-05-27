package vms

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/digitalocean/go-libvirt"
)

// Struct definition
type VirtualMachine struct {
	// Conn to hold the libvirt connection
	Conn       *libvirt.Libvirt `yaml:"-"`
	ApiVersion string           `yaml:"apiVersion"`
	Kind       string           `yaml:"kind"`
	Metadata   Metadata         `yaml:"metadata"`
	Spec       Spec             `yaml:"spec"`
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

func (vm *VirtualMachine) SetConnection(conn *libvirt.Libvirt) {
	vm.Conn = conn
}

func (vm *VirtualMachine) Header() *tabwriter.Writer {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "NAME\tSTATE\tCPU\tMEMORY\tDISK\tNETWORK\tOS\tAGE")

	return w
}
