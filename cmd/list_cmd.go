package cmd

import (
	"fmt"

	"github.com/kebairia/kvmcli/internal/logger"
	op "github.com/kebairia/kvmcli/internal/operations"
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
		if ConfigFile == "" {
			ConfigFile = "./configs/servers.yaml"
			// logger.Log.Fatalf("Configuration file is required (-f flag)")
		}
		op.GetAllVM(ConfigFile)
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
		if ConfigFile == "" {
			logger.Log.Fatalf("Configuration file is required (-f flag)")
		}
		fmt.Println("You networks are here")
		// op.ListSnapshost()
	},
}

func init() {
	GetVMCmd.Flags().
		StringVarP(&ConfigFile, "file", "f", "./configs/servers.yaml", "Path to the configuration file")
	GetNetworkCmd.Flags().
		StringVarP(&ConfigFile, "file", "f", "", "Configuration file for the VM(s)")
	GetSnapshotsCmd.Flags().
		StringVarP(&ConfigFile, "file", "f", "", "Configuration file for the VM(s)")
	GetCmd.AddCommand(GetVMCmd, GetSnapshotsCmd, GetNetworkCmd)
}
