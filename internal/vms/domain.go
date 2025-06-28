package vms

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"path/filepath"

	"github.com/digitalocean/go-libvirt"
	"github.com/kebairia/kvmcli/internal/database"
	"github.com/kebairia/kvmcli/internal/utils"
)

const (
	domainStateRunning = 1
	domainStatePaused  = 3
	domainStateStopped = 5
)

type DomainManager interface {
	// BuildXML returns the full libvirt XML for this VM.
	BuildXML(ctx context.Context, db *sql.DB, cfg VirtualMachineConfig) (string, error)

	// Define registers the domain but doesn’t start it.
	Define(ctx context.Context, xmlConfig string) error

	// Start actually powers on the domain.
	Start(ctx context.Context, name string) error

	// Stop gracefully shuts down the domain.
	Stop(ctx context.Context, name string) error

	// Stop (kills) shuts down the domain immediately.
	Destroy(ctx context.Context, name string) error

	// Undefine removes the domain metadata.
	Undefine(ctx context.Context, name string) error

	// State returns one of “Running”/“Stopped”/etc.
	State(ctx context.Context, name string) (string, error)
}
type LibvirtDomainManager struct {
	conn *libvirt.Libvirt
}

// NewLibvirtDomainManager constructs a new manager with the given libvirt client and logger.
func NewLibvirtDomainManager(conn *libvirt.Libvirt) *LibvirtDomainManager {
	return &LibvirtDomainManager{conn: conn}
}

func (d *LibvirtDomainManager) BuildXML(
	ctx context.Context,
	db *sql.DB,
	config VirtualMachineConfig,
) (string, error) {
	var st database.StoreRecord
	var err error
	st.ID, err = database.GetStoreIDByName(ctx, db, config.Metadata.Store)
	if err != nil {
		return "", nil
	}
	img, err := st.GetImageRecord(ctx, db, config.Spec.Image)
	if err != nil {
		return "", nil
	}

	// Build the disk image path for the domain configuration.
	diskImagePath := fmt.Sprintf(
		"%s.qcow2",
		filepath.Join(st.ImagesPath, config.Metadata.Name),
	)
	domain := utils.NewDomain(
		config.Metadata.Name,
		config.Spec.Memory,
		config.Spec.CPU,
		diskImagePath,
		config.Spec.Network.Name,
		config.Spec.Network.MacAddress,
		img.OsProfile,
	)
	xmlConfig, err := domain.GenerateXML()
	if err != nil {
		return "", fmt.Errorf("failed to generate XML for VM %s: %v", config.Metadata.Name, err)
	}
	return xml.Header + string(xmlConfig), nil
}

// Define registers the domain with libvirt (but does not start it).
func (m *LibvirtDomainManager) Define(ctx context.Context, xmlConfig string) error {
	dom, err := m.conn.DomainDefineXML(xmlConfig)
	if err != nil {
		return fmt.Errorf("define domain XML: %w", err)
	}
	_ = dom
	return nil
}

// Start powers on the given domain by name.
func (m *LibvirtDomainManager) Start(ctx context.Context, name string) error {
	dom, err := m.conn.DomainLookupByName(name)
	if err != nil {
		return fmt.Errorf("lookup domain %q: %w", name, err)
	}
	if err := m.conn.DomainCreate(dom); err != nil {
		return fmt.Errorf("start domain %q: %w", name, err)
	}
	return nil
}

// Stop attempts a graceful shutdown of the domain.
func (m *LibvirtDomainManager) Stop(ctx context.Context, name string) error {
	dom, err := m.conn.DomainLookupByName(name)
	if err != nil {
		return fmt.Errorf("lookup domain %q: %w", name, err)
	}
	if err := m.conn.DomainShutdown(dom); err != nil {
		return fmt.Errorf("shutdown domain %q: %w", name, err)
	}
	return nil
}

// Destroy force-stops (kills) the domain immediately.
func (m *LibvirtDomainManager) Destroy(ctx context.Context, name string) error {
	dom, err := m.conn.DomainLookupByName(name)
	if err != nil {
		return fmt.Errorf("lookup domain %q: %w", name, err)
	}
	if err := m.conn.DomainDestroy(dom); err != nil {
		return fmt.Errorf("destroy domain %q: %w", name, err)
	}
	return nil
}

// Undefine removes the domain’s metadata from libvirt (after it’s stopped).
func (m *LibvirtDomainManager) Undefine(ctx context.Context, name string) error {
	dom, err := m.conn.DomainLookupByName(name)
	if err != nil {
		return fmt.Errorf("lookup domain %q: %w", name, err)
	}
	if err := m.conn.DomainUndefine(dom); err != nil {
		return fmt.Errorf("undefine domain %q: %w", name, err)
	}
	return nil
}

// State returns a human-readable state (“Running”, “Paused”, “Shut off”, etc.).
func (m *LibvirtDomainManager) State(ctx context.Context, name string) (string, error) {
	dom, err := m.conn.DomainLookupByName(name)
	if err != nil {
		return "", fmt.Errorf("lookup domain %q: %w", name, err)
	}
	state, _, _, _, _, err := m.conn.DomainGetInfo(dom)
	if err != nil {
		return "", fmt.Errorf("get state for %q: %w", name, err)
	}
	switch int(state) {
	case 1:
		return "Running", nil
	case 3:
		return "Paused", nil
	case 5:
		return "Shut off", nil
	default:
		return "Unknown", nil
	}
}
