package operations

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/digitalocean/go-libvirt"
	"github.com/kebairia/kvmcli/internal"
	"github.com/kebairia/kvmcli/internal/database"
	"github.com/kebairia/kvmcli/internal/logger"
)

const (
	deviceName = "vda"

	// Domain state constants.
	domainStateRunning = 1
	domainStatePaused  = 3
	domainStateStopped = 5
)

// ListAllVMsInNamespace retrieves VMs for a given namespace from MongoDB,
// looks up their libvirt domain, and prints their details.
func ListAllVMsInNamespace(namespace string) {
	conn, err := internal.InitConnection()
	if err != nil {
		logger.Log.Fatalf("failed to connect to libvirt: %v", err)
	}
	// Ensure the connection is closed after processing.
	defer conn.Disconnect()

	// Retrieve VMs for the specific namespace from MongoDB.
	vms, err := database.GetObjectsByNamespace[database.VMRecord](namespace, database.VMsCollection)
	if err != nil {
		logger.Log.Errorf("failed to retrieve VMs for namespace %s: %v", namespace, err)
		return
	}
	if len(vms) == 0 {
		logger.Log.Infof("no VMs found in namespace %s", namespace)
		return
	}

	// Create a tabwriter for neat column formatting.
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "NAME\tSTATE\tCPU\tMEMORY\tDISK\tNETWORK\tOS\tAGE")

	for _, vm := range vms {
		domain, err := conn.DomainLookupByName(vm.Name)
		if err != nil {
			logger.Log.Errorf("failed to lookup domain for VM %s: %v", vm.Name, err)
			continue
		}

		// Retrieve the VM state.
		state, err := getState(conn, domain, vm.Name)
		if err != nil {
			logger.Log.Errorf("failed to get state for VM %s: %v", vm.Name, err)
			state = "Unknown"
		}

		// Retrieve the VM disk size.
		diskSizeGB, err := getDiskSize(conn, domain, vm.Name)
		if err != nil {
			logger.Log.Errorf("failed to get disk size for VM %s: %v", vm.Name, err)
			diskSizeGB = 0
		}

		// Print the VM details.
		fmt.Fprintf(w, "%s\t%s\t%d\t%d GB\t%.2f GB\t%s\t%s\t%s\n",
			vm.Name,
			state,
			vm.CPU,
			vm.RAM/1024,
			diskSizeGB,
			vm.Network,
			vm.Image,
			formatAge(vm.CreatedAt),
		)
	}
	w.Flush()
}

// ListAllVMs retrieves all VM domains via libvirt,
// then gets corresponding details from MongoDB and prints them.
func ListAllVMs() {
	conn, err := internal.InitConnection()
	if err != nil {
		logger.Log.Fatalf("failed to connect to libvirt: %v", err)
	}
	defer conn.Disconnect()

	flags := libvirt.ConnectListDomainsActive | libvirt.ConnectListDomainsInactive
	domains, _, err := conn.ConnectListAllDomains(1, flags)
	if err != nil {
		logger.Log.Fatalf("failed to retrieve domains: %v", err)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "NAME\tSTATE\tCPU\tMEMORY\tDISK\tNETWORK\tOS\tAGE")

	for _, domain := range domains {
		vm, err := database.GetRecord[database.VMRecord](domain.Name, database.VMsCollection)
		if err != nil {
			logger.Log.Errorf("failed to get details for VM %s: %v", domain.Name, err)
			continue
		}

		state, err := getState(conn, domain, vm.Name)
		if err != nil {
			logger.Log.Errorf("failed to get state for VM %s: %v", vm.Name, err)
			state = "Unknown"
		}

		diskSizeGB, err := getDiskSize(conn, domain, vm.Name)
		if err != nil {
			logger.Log.Errorf("failed to get disk size for VM %s: %v", vm.Name, err)
			diskSizeGB = 0
		}

		fmt.Fprintf(w, "%s\t%s\t%d\t%d GB\t%.2f GB\t%s\t%s\t%s\n",
			vm.Name,
			state,
			vm.CPU,
			vm.RAM/1024,
			diskSizeGB,
			vm.Network,
			vm.Image,
			formatAge(vm.CreatedAt),
		)
	}
	w.Flush()
}

// getState returns a string representation of the VM state based on its domain info.
func getState(conn *libvirt.Libvirt, domain libvirt.Domain, name string) (string, error) {
	state, _, _, _, _, err := conn.DomainGetInfo(domain)
	if err != nil {
		return "", fmt.Errorf("failed to get info for domain %s: %w", name, err)
	}

	switch int(state) {
	case domainStateRunning:
		return "Running", nil
	case domainStatePaused:
		return "Paused", nil
	case domainStateStopped:
		return "Stopped", nil
	default:
		return "Unknown", nil
	}
}

// getDiskSize returns the disk size (in gigabytes) for the specified VM domain.
func getDiskSize(conn *libvirt.Libvirt, domain libvirt.Domain, name string) (float64, error) {
	_, _, diskPhysSize, err := conn.DomainGetBlockInfo(domain, deviceName, 0)
	if err != nil {
		return 0, fmt.Errorf("failed to get block info for domain %s: %w", name, err)
	}

	return float64(diskPhysSize) / (1024 * 1024 * 1024), nil
}

// formatAge returns a human-friendly string for the time elapsed since t.
func formatAge(t time.Time) string {
	duration := time.Since(t)
	if duration < 0 {
		duration = -duration
	}

	if days := int(duration.Hours() / 24); days >= 1 {
		return fmt.Sprintf("%dd", days)
	}
	if hours := int(duration.Hours()); hours >= 1 {
		return fmt.Sprintf("%dh", hours)
	}
	if minutes := int(duration.Minutes()); minutes >= 1 {
		return fmt.Sprintf("%dm", minutes)
	}
	return fmt.Sprintf("%ds", int(duration.Seconds()))
}
