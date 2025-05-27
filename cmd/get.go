package cmd

import (
	"fmt"

	"github.com/kebairia/kvmcli/internal/logger"
	"github.com/kebairia/kvmcli/internal/network"
	"github.com/kebairia/kvmcli/internal/operations"
	"github.com/spf13/cobra"
)

// Create the "get" parent command.
var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "Retrieve information about VMs, snapshots, networks, etc.",
}

// 'get vm' subcommand: shows virtual machines.
var GetVMCmd = &cobra.Command{
	Use:   "vm",
	Short: "Display information about virtual machines",
	Run: func(cmd *cobra.Command, args []string) {
		if Namespace != "" {
			operations.ListByNamespace(Namespace)
			// operations.ListAll()
			return
		}
		operations.ListAll()
	},
}

// 'get snapshots' subcommand: shows snapshots.
var GetSnapshotsCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Display snapshots for virtual machines",
	Run: func(cmd *cobra.Command, args []string) {
		if ConfigFile == "" {
			logger.Log.Fatalf("Configuration file is required (-f flag)")
		}
		fmt.Println("You snapshosts are here")
		// op.ListSnapshost()
	},
}

// 'get snapshots' subcommand: shows snapshots.
var GetNetworkCmd = &cobra.Command{
	Use:   "network",
	Short: "Display network details",
	Run: func(cmd *cobra.Command, args []string) {
		if Namespace != "" {
			network.ListNetworksByNamespace(Namespace)
			return
		}
		network.ListAllNetworks()
	},
}

func init() {
	// Flags for virtual machines
	GetVMCmd.Flags().
		StringVarP(&Namespace, "namespace", "n", "", "Namespace")
		// Flags for Networks
	GetNetworkCmd.Flags().
		StringVarP(&Namespace, "namespace", "n", "", "Namespace")
		// Flags for Snapshots
	GetSnapshotsCmd.Flags().
		StringVarP(&Namespace, "namespace", "n", "", "Namespace")
	GetCmd.AddCommand(GetVMCmd, GetSnapshotsCmd, GetNetworkCmd)
}
