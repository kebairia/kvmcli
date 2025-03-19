package vms

import (
	"fmt"
	"net/url"
	"os"
	"text/tabwriter"

	"github.com/digitalocean/go-libvirt"
	"github.com/kebairia/kvmcli/internal/logger"
)

const deviceName = "vda"

type VMInfo struct {
	Name    string
	CPU     int
	Memory  int
	Disk    float64
	Network string
	Status  string
	OS      string
}

func ListAllVM(configPath string) {
	uri, err := url.Parse(string(libvirt.QEMUSystem))
	if err != nil {
		logger.Log.Println(err)
	}
	l, err := libvirt.ConnectToURI(uri)
	if err != nil {
		logger.Log.Fatalf("failed to connect: %v", err)
	}
	flags := libvirt.ConnectListDomainsActive | libvirt.ConnectListDomainsInactive
	domains, _, err := l.ConnectListAllDomains(1, flags)
	if err != nil {
		logger.Log.Fatalf("can't retreive domains infos: %v", err)
	}

	// read config file, return error if you failed
	// var vms []config.VirtualMachine
	// if vms, err = config.LoadConfig(configPath); err != nil {
	// 	logger.Log.Fatalf("%s", err)
	// }
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "NAME\tSTATUS\tCPU\tMEMORY\tDISK\tNETWORK\tOS")

	manager := VMManager{conn}

	for _, domain := range domains {
		info := manager.GetInfo(domain.Name)

		// Use tabwriter to format the output similar to kubectl
		// Print header similar to Kubernetes (you can add more columns if needed)
		// fmt.Fprintln(w, "----\t------\t----")
		fmt.Fprintf(w, "%s\t%s\t%d\t%d GB\t%.2f GB\t%s\n",
			info.Name,
			info.Status,
			info.CPU,
			info.Memory,
			info.Disk,
			info.Network,
			// vm.Spec.Image,
		)

	}
	w.Flush()
}

func GetVMInfo(vmName string, conn *libvirt.Libvirt) *VMInfo {
// func GetVMInfo(vmName string, conn *libvirt.Libvirt) *VMInfo {
func (m *VMManager) GetInfo(name string) *VMInfo {
	var info VMInfo
	domain, err := m.Conn.DomainLookupByName(name)
	if err != nil {
		logger.Log.Fatalf("Domain lookup failed for %s: %v", name, err)
	}
	// Retrieve basic domain info (state, memory, and CPU count)
	state, _, mem, cpu, _, err := m.Conn.DomainGetInfo(domain)
	if err != nil {
		logger.Log.Fatalf("Failed to get info for domain %s: %v", name, err)
	}

	// Retrieve disk block info for device "vda"
	_, _, diskPhysSize, err := m.Conn.DomainGetBlockInfo(domain, deviceName, 0)
	if err != nil {
		logger.Log.Fatalf("Failed to get block info for domain %s: %v", name, err)
	}

	// Determine status based on the domain state (assumes state 1 means running)
	status := "Stopped"
	if state == 1 {
		status = "Running"
	}
	// Populate the VMInfo struct.
	info = VMInfo{
		Name:    name,
		CPU:     int(cpu),
		Memory:  int(mem) / (1024 * 1024), // Convert from kilobytes to GiB.
		Disk:    float64(diskPhysSize) / (1024 * 1024 * 1024),
		Network: "homelab", // This value is hard-coded; consider retrieving from config or domain XML.
		Status:  status,
	}

	return &info
}
