package vms

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/digitalocean/go-libvirt"
	db "github.com/kebairia/kvmcli/internal/database"
	log "github.com/kebairia/kvmcli/internal/logger"
)

const allDomains = -1

// VirtualMachineInfo holds everything we need to print one row.
type VirtualMachineInfo struct {
	Name      string
	State     string
	CPU       int
	RAM       int     // in MB
	Disk      float64 // in GB
	Network   string
	Image     string
	CreatedAt time.Time
}

func (info *VirtualMachineInfo) Header() *tabwriter.Writer {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "NAME\tSTATE\tCPU\tMEMORY\tDISK\tNETWORK\tOS\tAGE")
	return w
}

func (info *VirtualMachineInfo) PrintInfo(w *tabwriter.Writer) {
	fmt.Fprintf(w, "%s\t%s\t%d\t%d MB\t%.2f GB\t%s\t%s\t%s\n",
		info.Name,
		info.State,
		info.CPU,
		info.RAM, // Convert to MB
		info.Disk,
		info.Network,
		info.Image,
		formatAge(info.CreatedAt),
	)
}

// NOTE:: Later, it must be (ListResourcesByNamespace), something like that
// Since namespaces are subset in (all), I can add it as a condition.
// if no namespace is required, I list all resources
// of course resources are specified as argument or something like that

// NewVirtualMachineInfo constructs a VirtualMachineInfo by querying libvirt and the DB.
func NewVirtualMachineInfo(
	ctx context.Context,
	conn *libvirt.Libvirt,
	rec db.VirtualMachineRecord,
) (*VirtualMachineInfo, error) {
	// log := logger
	// Domain lookup
	dom, err := conn.DomainLookupByName(rec.Name)
	if err != nil {
		return nil, fmt.Errorf("lookup domain %q: %w", rec.Name, err)
	}

	// State
	state, err := getState(conn, dom)
	if err != nil {
		log.Errorf("cannot get state for %q: %v", rec.Name, err)
		state = "unknown"
	}

	// Disk size
	disk, err := getDiskSize(conn, dom)
	if err != nil {
		log.Errorf("cannot get disk size for %q: %v", rec.Name, err)
	}

	// Network name
	network, err := db.GetNetworkNameByID(ctx, db.DB, rec.NetworkID)
	if err != nil {
		log.Errorf("cannot get network name for %q: %v", rec.Name, err)
	}

	return &VirtualMachineInfo{
		Name:      rec.Name,
		State:     state,
		CPU:       rec.CPU,
		RAM:       rec.RAM,
		Disk:      disk,
		Network:   network,
		Image:     rec.Image,
		CreatedAt: rec.CreatedAt,
	}, nil
}

func GetVirtualMachines(conn *libvirt.Libvirt) ([]VirtualMachineInfo, error) {

	records, err := db.GetRecords(db.Ctx, db.DB)
	if err != nil {
		return nil, fmt.Errorf("failed to get VM records  %w", err)
	}

	vms := make([]VirtualMachineInfo, 0, len(records))

	for _, rec := range records {
		vmInfo, err := NewVirtualMachineInfo(db.Ctx, conn, rec)
		if err != nil {
			log.Errorf("could not build VM info for %q: %v", rec.Name, err)
			continue
		}
		vms = append(vms, *vmInfo)
	}

	return vms, nil
}

func GetVirtualMachineByNamespace(
	conn *libvirt.Libvirt,
	namespace string,
) ([]VirtualMachineInfo, error) {

	records, err := db.GetRecordsByNamespace(db.Ctx, db.DB, namespace, db.VMsTable)
	if err != nil {
		return nil, fmt.Errorf("failed to get VM records for namespace %q: %w", namespace, err)
	}

	vms := make([]VirtualMachineInfo, 0, len(records))

	for _, rec := range records {
		vmInfo, err := NewVirtualMachineInfo(db.Ctx, conn, rec)
		if err != nil {
			log.Errorf("could not build VM info for %q: %v", rec.Name, err)
			continue
		}
		vms = append(vms, *vmInfo)
	}

	return vms, nil
}

// FIX: This must be on database package, not here.
// func (vm *VirtualMachine) Get(
// 	ctx context.Context,
// 	conn *libvirt.Libvirt,
// ) (*VirtualMachineInfo, error) {
// 	// NOTE: <-----
// 	// defer conn.Disconnect()
//
// 	// Get the domain by name.
// 	domain, err := conn.DomainLookupByName(vm.Metadata.Name)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to lookup domain %s: %w", vm.Metadata.Name, err)
// 	}
//
// 	// Retrieve the VM record from the database.
// 	rec := db.VirtualMachineRecord{}
// 	err = rec.GetRecord(db.Ctx, db.DB, vm.Metadata.Name)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get VM record: %w", err)
// 	}
//
// 	// Load the VM information.
// 	// info := &VirtualMachineInfo{}
// 	// Load the VM info using the connection, domain, and record.
//
// 	state, err := getState(conn, domain)
// 	if err != nil {
// 		log.Error("state(%q): %v", rec.Name, err)
// 		state = "Unknown"
// 	}
// 	// disk
// 	disk, err := getDiskSize(conn, domain)
// 	if err != nil {
// 		log.Error("disk(%q): %v", rec.Name, err)
// 	}
// 	// network name lookup
// 	network, err := db.GetNetworkNameByID(db.Ctx, db.DB, rec.NetworkID)
// 	if err != nil {
// 		log.Error("network(%q): %v", rec.Name, err)
// 	}
//
// 	// err = info.Load(conn, domain, rec)
// 	// if err != nil {
// 	// 	return nil, fmt.Errorf("failed to load VM info: %w", err)
// 	// }
//
// 	return &VirtualMachineInfo{
// 		Name:      rec.Name,
// 		State:     state,
// 		CPU:       rec.CPU,
// 		RAM:       rec.RAM / 1024, // Convert to MB
// 		Disk:      disk,
// 		Network:   network,
// 		Image:     rec.Image,
// 		CreatedAt: rec.CreatedAt,
// 	}, nil
// }
