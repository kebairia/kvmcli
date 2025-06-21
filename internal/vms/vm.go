package vms

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/digitalocean/go-libvirt"
)

// Error definitions
var (
	ErrNilLibvirtConn   = errors.New("libvirt connection is nil")
	ErrNilDBConn        = errors.New("database connection is nil")
	ErrNilDiskManager   = errors.New("disk manager is nil")
	ErrNilDomainManager = errors.New("domain manager is nil")
)

// VirtualMachine represents a VM with injected dependencies.
type VirtualMachine struct {
	ctx    context.Context
	conn   *libvirt.Libvirt
	db     *sql.DB
	Config VirtualMachineConfig
	disk   DiskManager
	domain DomainManager
	// network NetworkManager
}

// VirtualMachineOption is a functional option for configuring a VirtualMachine.
type VirtualMachineOption func(*VirtualMachine)

// WithLibvirtConnection injects a non-nil Libvirt connection.
func WithLibvirtConnection(conn *libvirt.Libvirt) VirtualMachineOption {
	return func(vm *VirtualMachine) {
		vm.conn = conn
	}
}

// WithDatabaseConnection injects a non-nil *sql.DB.
func WithDatabaseConnection(db *sql.DB) VirtualMachineOption {
	return func(vm *VirtualMachine) {
		vm.db = db
	}
}

// WithContext injects a context.Context; if nil is passed, a background context is set.
func WithContext(ctx context.Context) VirtualMachineOption {
	return func(vm *VirtualMachine) {
		if ctx == nil {
			vm.ctx = context.Background()
		} else {
			vm.ctx = ctx
		}
	}
}

// New option to inject a DiskManager:
func WithDiskManager(d DiskManager) VirtualMachineOption {
	return func(vm *VirtualMachine) {
		vm.disk = d
	}
}

// New option to inject a DomainManager:
func WithDomainManager(d DomainManager) VirtualMachineOption {
	return func(vm *VirtualMachine) {
		vm.domain = d
	}
}

// NewVirtualMachine constructs a VM, applies options, and validates dependencies.
func NewVirtualMachine(
	cfg VirtualMachineConfig,
	opts ...VirtualMachineOption,
) (*VirtualMachine, error) {
	vm := &VirtualMachine{
		Config: cfg,
		ctx:    context.Background(),
	}

	// apply options
	for _, opt := range opts {
		opt(vm)
	}

	// wire defaults if not provided
	// if vm.disk == nil {
	// 	// user must provide via option; no hard-coded defaults here
	// 	return nil, ErrNilDiskManager
	// }
	// if vm.domain == nil {
	// 	return nil, ErrNilDomainManager
	// }

	if vm.disk == nil {
		vm.disk = &QemuDiskManager{
			QemuImgPath:    "/usr/bin/qemu-img",
			BaseImagesPath: "/home/zakaria/dox/homelab/artifacts/rocky",
			DestImagesPath: "/home/zakaria/dox/homelab/images/",
			Timeout:        10 * time.Second,
		}
	}
	vm.domain = NewLibvirtDomainManager(vm.conn)

	// validate core dependencies
	if vm.conn == nil {
		return nil, ErrNilLibvirtConn
	}
	if vm.db == nil {
		return nil, ErrNilDBConn
	}

	return vm, nil
}
