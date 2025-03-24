package vms

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/digitalocean/go-libvirt"
	"github.com/kebairia/kvmcli/internal/logger"
	"github.com/kebairia/kvmcli/internal/utils"
)

const imagesPath = "/home/zakaria/dox/homelab/images/"

type VirtualMachine struct {
	// Conn to hold the libvirt connection
	Conn       *libvirt.Libvirt
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

func (vm *VirtualMachine) Create() error {
	// Check connection
	if vm.Conn == nil {
		return fmt.Errorf("libvirt connection is nil")
	}
	// Create overlay image
	if err := CreateOverlay("rocky.qcow2", vm.Spec.Disk.Path); err != nil {
		return fmt.Errorf("Failed to create overlay for VM %q: %v", vm.Metadata.Name, err)
	}
	// Creating domain out of infos
	domain := utils.NewDomain(

		vm.Metadata.Name,
		vm.Spec.Memory,
		vm.Spec.CPU,
		vm.Spec.Disk.Path,
		vm.Spec.Network.MacAddress,
	)
	xmlConfig, err := domain.GenerateXML()
	if err != nil {
		return fmt.Errorf("failed to generate XML for VM %s: %v", vm.Metadata.Name, err)
	}

	vmInstance, err := vm.Conn.DomainDefineXML(string(xmlConfig))
	if err != nil {
		return fmt.Errorf("failed to define domain for VM %s: %v", vmInstance.Name, err)
	}
	if err := vm.Conn.DomainCreate(vmInstance); err != nil {
		return fmt.Errorf("failed to start VM %s: %w", vmInstance.Name, err)
	}

	fmt.Printf("vm/%s created\n", vm.Metadata.Name)
	return nil
}

func (vm *VirtualMachine) Delete() error {
	// Check connection
	if vm.Conn == nil {
		return fmt.Errorf("libvirt connection is nil")
	}

	vmName := vm.Metadata.Name
	domain, err := vm.Conn.DomainLookupByName(vmName)
	if err != nil {
		return fmt.Errorf("Failed to find VM %s: %w", vmName, err)
	}

	// Attempt to destroy the domain.
	if err := vm.Conn.DomainDestroy(domain); err != nil {
		return fmt.Errorf(
			"failed to delete VM %q (it might not be running): %w",
			vmName,
			err,
		)
	}

	// Undefine the domain
	if err := vm.Conn.DomainUndefine(domain); err != nil {
		return fmt.Errorf("failed to undefine VM %q: %w", vmName, err)
	}

	// Remove the disk associated with the VM.
	diskPath := filepath.Join(imagesPath, vmName+".qcow2")
	if err := os.Remove(diskPath); err != nil {
		return fmt.Errorf("failed to delete disk for VM %q: %w", vmName, err)
	}

	logger.Log.Infof("%s/%s deleted", "vm", vmName)

	return nil
}

func (vm *VirtualMachine) SetConnection(conn *libvirt.Libvirt) {
	vm.Conn = conn
}
