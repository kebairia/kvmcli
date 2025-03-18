package cmd

import (
	"github.com/kebairia/kvmcli/internal/logger"
	"github.com/spf13/cobra"
)

// Global variables to store flag values
var (
	ConfigFile  string // To store the path of the config file
	ClusterFile string // To store the path of the cluster file
	Provision   bool   // To determine if provisioning should start
	Verbose     bool   // To determine if output should be in verbose
)

// Create the root command
var rootCmd = &cobra.Command{
	Use:   "kvmcli",
	Short: "kvmcli is a CLI for managing KVM virtual machines",
	Long:  "A CLI similar to kubectl for creating, deleting, and managing KVM VMs.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.Log.Errorf("%s", err)
	}
}

func init() {
	rootCmd.AddCommand(CreateCmd)
	rootCmd.AddCommand(DeleteCmd)
	rootCmd.AddCommand(GetCmd)
	rootCmd.AddCommand(InitVMCmd)
}
