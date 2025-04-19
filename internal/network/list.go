package network

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/digitalocean/go-libvirt"
	"github.com/kebairia/kvmcli/internal"
	db "github.com/kebairia/kvmcli/internal/database-sql"
	"github.com/kebairia/kvmcli/internal/logger"
)

const (
	// Network state constants as defined by libvirt.
	networkStateActive   = 1
	networkStateInactive = 0
)

// ListNetworksByNamespace retrieves networks for a given namespace from MongoDB,
// looks up their corresponding libvirt instances, and prints the details in a tabular format.
func ListNetworksByNamespace(namespace string) {
	conn, err := internal.InitConnection()
	if err != nil {
		logger.Log.Fatalf("failed to connect to libvirt: %v", err)
	}
	defer conn.Disconnect()

	networks, err := db.GetNetworkObjectsByNamespace(
		db.Ctx,
		db.DB,
		namespace,
		db.NetworksTable,
	)
	if err != nil {
		logger.Log.Errorf("failed to retrieve VMs for namespace %s: %v", namespace, err)
		return
	}
	if len(networks) == 0 {
		logger.Log.Infof("no networks found in namespace %q", namespace)
		return
	}

	// Setup tabwriter for clean columnar output.
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "NAME\tSTATE\tBRIDGE\tSUBNET\tGATEWAY\tDHCP RANGE\tAGE")

	// Process each network record.
	for _, nwRecord := range networks {

		// Lookup the network instance via libvirt.
		netInstance, err := conn.NetworkLookupByName(nwRecord.Name)
		if err != nil {
			logger.Log.Errorf("failed to lookup network %q: %v", nwRecord.Name, err)
			continue
		}

		// Get the current state of the network.
		state, err := getState(conn, netInstance, nwRecord.Name)
		if err != nil {
			logger.Log.Errorf("failed to get state for network %q: %v", nwRecord.Name, err)
			state = "Unknown"
		}

		// Format the DHCP range.
		dhcpRange := nwRecord.DHCP["start"] + " → " + nwRecord.DHCP["end"]

		// Print network information.
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			nwRecord.Name,
			state,
			nwRecord.Bridge,
			nwRecord.Netmask,
			nwRecord.NetAddress,
			dhcpRange,
			formatAge(nwRecord.CreatedAt),
		)
	}
	w.Flush()
}

// ListAllNetworks retrieves all networks (active and inactive) from libvirt,
// gets additional details from the database, and prints them in a tabular format.
func ListAllNetworks() {
	// Initialize libvirt connection.
	conn, err := internal.InitConnection()
	if err != nil {
		logger.Log.Fatalf("failed to initialize libvirt connection: %v", err)
	}
	defer conn.Disconnect()

	// Define flags to list both active and inactive networks.
	flags := libvirt.ConnectListNetworksActive | libvirt.ConnectListNetworksInactive

	// Retrieve networks.
	networks, _, err := conn.ConnectListAllNetworks(1, flags)
	if err != nil {
		logger.Log.Fatalf("failed to retrieve networks: %v", err)
	}

	// Setup tabwriter for clean columnar output.
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "NAME\tSTATE\tBRIDGE\tSUBNET\tGATEWAY\tDHCP RANGE\tAGE")

	// Process each network.
	for _, network := range networks {
		// Retrieve network details from the database.
		networkDetails, err := db.GetNetworkRecord(
			db.Ctx,
			db.DB,
			network.Name,
			db.NetworksTable,
		)
		if err != nil {
			logger.Log.Errorf("failed to get details for network %s: %v", network.Name, err)
			continue
		}

		// Get the current state of the network.
		state, err := getState(conn, network, network.Name)
		if err != nil {
			logger.Log.Errorf("failed to get state for network %s: %v", network.Name, err)
		}

		// Format the DHCP range.
		dhcpRange := networkDetails.DHCP["start"] + " → " + networkDetails.DHCP["end"]

		// Print network information.
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			networkDetails.Name,
			state,
			networkDetails.Bridge,
			networkDetails.Netmask,
			networkDetails.NetAddress,
			dhcpRange,
			formatAge(networkDetails.CreatedAt),
		)
	}
	w.Flush()
}

// getState retrieves the active/inactive state of the network using libvirt.
func getState(conn *libvirt.Libvirt, network libvirt.Network, name string) (string, error) {
	state, err := conn.NetworkIsActive(network)
	if err != nil {
		return "", fmt.Errorf("failed to get state for network %s: %w", name, err)
	}

	switch int(state) {
	case networkStateActive:
		return "Active", nil
	case networkStateInactive:
		return "Inactive", nil
	default:
		return "Unknown", nil
	}
}

// formatAge returns a human-friendly string for the time elapsed since t.
func formatAge(t time.Time) string {
	duration := time.Since(t)
	if duration < 0 {
		duration = -duration
	}

	switch {
	case duration.Hours() >= 24:
		return fmt.Sprintf("%dd", int(duration.Hours()/24))
	case duration.Hours() >= 1:
		return fmt.Sprintf("%dh", int(duration.Hours()))
	case duration.Minutes() >= 1:
		return fmt.Sprintf("%dm", int(duration.Minutes()))
	default:
		return fmt.Sprintf("%ds", int(duration.Seconds()))
	}
}
