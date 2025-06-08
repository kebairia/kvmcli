package vms

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/digitalocean/go-libvirt"
	db "github.com/kebairia/kvmcli/internal/database"
	log "github.com/kebairia/kvmcli/internal/logger"
)

const allDomains = -1

const (
	vmCols = `id, name, namespace, cpu, ram, mac_address, network_id, image, disk_size, disk_path, created_at, labels`
)

// VirtualMachineInfo holds everything we need to print one row.
type VirtualMachineInfo struct {
	Name     string
	State    string
	CPU      int
	RAM      int     // in MB
	DiskSize float64 // in GB
	Network  string
	OS       string
	Age      string
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
		info.DiskSize,
		info.Network,
		info.OS,
		info.Age,
	)
}

// NewVirtualMachineInfo constructs a VirtualMachineInfo by querying libvirt and the DB.
func NewVirtualMachineInfo(
	ctx context.Context,
	database *sql.DB,
	conn *libvirt.Libvirt,
	rec db.VirtualMachineRecord,
) (*VirtualMachineInfo, error) {
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
	network, err := db.GetNetworkNameByID(ctx, database, rec.NetworkID)
	if err != nil {
		log.Errorf("cannot get network name for %q: %v", rec.Name, err)
	}

	return &VirtualMachineInfo{
		Name:     rec.Name,
		State:    state,
		CPU:      rec.CPU,
		RAM:      rec.RAM,
		DiskSize: disk,
		Network:  network,
		OS:       rec.Image,
		Age:      formatAge(rec.CreatedAt),
	}, nil
}

func GetVirtualMachines(
	ctx context.Context,
	database *sql.DB,
	conn *libvirt.Libvirt,
) ([]VirtualMachineInfo, error) {
	// records, err := db.GetRecords[*db.VirtualMachineRecord](ctx, database, "", db.VMsTable, vmCols)
	records, err := db.GetVMRecords(ctx, database, "")
	if err != nil {
		return nil, fmt.Errorf("failed to get VM records  %w", err)
	}

	vms := make([]VirtualMachineInfo, 0, len(records))

	for _, rec := range records {
		vmInfo, err := NewVirtualMachineInfo(ctx, database, conn, rec)
		if err != nil {
			log.Errorf("could not build VM info for %q: %v", rec.Name, err)
			continue
		}
		vms = append(vms, *vmInfo)
	}

	return vms, nil
}
