package operations

import (
	"fmt"
	"net/url"
	"os"
	"text/tabwriter"

	"github.com/digitalocean/go-libvirt"
	"github.com/kebairia/kvmcli/internal/config"
	"github.com/kebairia/kvmcli/internal/logger"
)

type VMInfo struct {
	Name    string
	CPU     int
	Memory  int
	Disk    float64
	Network string
	Status  string
}

func GetAllVM(configPath string) {
	uri, err := url.Parse(string(libvirt.QEMUSystem))
	if err != nil {
		fmt.Println(err)
	}
	l, err := libvirt.ConnectToURI(uri)
	if err != nil {
		logger.Log.Fatalf("failed to connect: %v", err)
	}

	config, err := config.LoadConfig(configPath)
	if err != nil {
		logger.Log.Fatalf("failed to connect: %v", err)
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "NAME\tCPU\tMEMORY\tDISK\tNETWORK\tSTATUS")
	for vmName := range config.VMs {
		info := GetVMInfo(vmName, l)

		// Use tabwriter to format the output similar to kubectl
		// Print header similar to Kubernetes (you can add more columns if needed)
		// fmt.Fprintln(w, "----\t------\t----")
		fmt.Fprintf(w, "%s\t%d\t%d GiB\t%.2f GiB\t%s\t%s\n",
			info.Name,
			info.CPU,
			info.Memory,
			info.Disk,
			info.Network,
			info.Status,
		)

	}
	w.Flush()
}

func GetVMInfo(vmName string, conn *libvirt.Libvirt) VMInfo {
	vm := VMInfo{}

	domain, err := conn.DomainLookupByName(vmName)
	if err != nil {
		logger.Log.Fatalf("Domain lookup failed for %s: %v", vmName, err)
	}
	// Retrieve basic domain info (state, memory, and CPU count)
	state, _, mem, cpu, _, err := conn.DomainGetInfo(domain)
	if err != nil {
		logger.Log.Fatalf("Failed to get info for domain %s: %v", vmName, err)
	}

	// Retrieve disk block info for device "vda"
	_, _, diskPhysSize, err := conn.DomainGetBlockInfo(domain, "vda", 0)
	if err != nil {
		logger.Log.Fatalf("Failed to get block info for domain %s: %v", vmName, err)
	}

	// Determine status based on the domain state (assumes state 1 means running)
	status := "Stopped"
	if state == 1 {
		status = "Running"
	}
	// Populate the VMInfo struct
	vm.Name = vmName
	vm.CPU = int(cpu)
	vm.Memory = int(
		mem / 1024 / 1024,
	) // Note: mem is returned in kilobytes; adjust if needed.
	vm.Disk = float64(diskPhysSize / 1024 / 1024 / 1024) // Convert bytes to MB.
	vm.Network = "homelab"                               // This is hard-coded; update as needed.
	vm.Status = status

	return vm
}
