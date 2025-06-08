package vms

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/digitalocean/go-libvirt"
	"gopkg.in/yaml.v3"
)

// ErrNilLibvirtConn is returned when the Libvirt connection is missing.
var ErrNilLibvirtConn = errors.New("libvirt connection is nil")

// ErrNilDBConn is returned when the database connection is missing.
var ErrNilDBConn = errors.New("database connection is nil")

// VirtualMachine represents a VM specification (loaded from YAML) and its runtime dependencies.
type VirtualMachine struct {
	// runtime-only dependencies (never part of YAML)
	Conn    *libvirt.Libvirt `yaml:"-"`
	DB      *sql.DB          `yaml:"-"`
	Context context.Context  `yaml:"-"`

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

// Resources defines CPU and memory for the VM.
type Resources struct{}

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

// VirtualMachineOption is a functional option for configuring a VirtualMachine.
type VirtualMachineOption func(*VirtualMachine)

// WithLibvirtConnection injects a non-nil Libvirt connection.
func WithLibvirtConnection(conn *libvirt.Libvirt) VirtualMachineOption {
	return func(vm *VirtualMachine) {
		vm.Conn = conn
	}
}

// WithDatabaseConnection injects a non-nil *sql.DB.
func WithDatabaseConnection(db *sql.DB) VirtualMachineOption {
	return func(vm *VirtualMachine) {
		vm.DB = db
	}
}

// WithContext injects a context.Context; if nil is passed, a background context is set.
func WithContext(ctx context.Context) VirtualMachineOption {
	return func(vm *VirtualMachine) {
		if ctx == nil {
			vm.Context = context.Background()
		} else {
			vm.Context = ctx
		}
	}
}

func (vm *VirtualMachine) SetConnection(ctx context.Context, db *sql.DB, conn *libvirt.Libvirt) {
	vm.Conn = conn
	vm.DB = db
	vm.Context = ctx
}

// NewVirtualMachine parses the YAML manifest into a VirtualMachine struct,
// applies any options, and validates required dependencies.
func NewVirtualMachine(manifest []byte, opts ...VirtualMachineOption) (*VirtualMachine, error) {
	var vm VirtualMachine

	if err := yaml.Unmarshal(manifest, &vm); err != nil {
		return nil, fmt.Errorf("failed to parse VM manifest: %w", err)
	}

	// Apply all options (some of them might set vm.DB, vm.Conn, vm.Context, etc.)
	for _, opt := range opts {
		opt(&vm)
	}

	// Ensure required dependencies are present
	if vm.Conn == nil {
		return nil, ErrNilLibvirtConn
	}
	if vm.DB == nil {
		return nil, ErrNilDBConn
	}
	if vm.Context == nil {
		vm.Context = context.Background()
	}
	return &vm, nil
}
