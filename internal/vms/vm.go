package vms

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/digitalocean/go-libvirt"
	"github.com/kebairia/kvmcli/internal/database"
	db "github.com/kebairia/kvmcli/internal/database"
	"github.com/kebairia/kvmcli/internal/logger"
)

const imagesPath = "/home/zakaria/dox/homelab/images/"

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

func (vm *VirtualMachine) Create() error {
	// Check connection
	if vm.Conn == nil {
		return fmt.Errorf("libvirt connection is nil")
	}
	// Initiliaze a new vm record
	record := NewVMRecord(vm)

	// Insert the vm record
	_, err := db.Insert(record)
	if err != nil {
		return fmt.Errorf("failed to create database record for %q: %w", vm.Metadata.Name, err)
	}

	// Create overlay image
	if err := CreateOverlay("rocky.qcow2", vm.Spec.Disk.Path); err != nil {
		return fmt.Errorf("Failed to create overlay for VM %q: %v", vm.Metadata.Name, err)
	}

	// Prepare the domain and generate its XML configuration.
	xmlConfig, err := vm.prepareDomain()
	if err != nil {
		logger.Log.Errorf("%v", err)
	}
	// Define the domain and start the VM.
	if err := vm.defineAndStartDomain(xmlConfig); err != nil {
		logger.Log.Errorf("%v", err)
	}

	fmt.Printf("vm/%s created\n", vm.Metadata.Name)
	return nil
}

// Delete function
func (vm *VirtualMachine) Delete() error {
	var err error
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
	err = database.Delete(vm.Metadata.Name)
	if err != nil {
		logger.Log.Errorf("failed to delete record for VM %s: %v", vm.Metadata.Name, err)
	}
	logger.Log.Infof("%s/%s deleted", "vm", vmName)

	return nil
}

func (vm *VirtualMachine) SetConnection(conn *libvirt.Libvirt) {
	vm.Conn = conn
}
