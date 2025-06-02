package cmd

import (
	log "github.com/kebairia/kvmcli/internal/logger"
	"github.com/spf13/cobra"
)

// Global flag variables.
var (
	Namespace    string // Namespace
	ManifestPath string // Path of the manifest file.
	ConfigFile   string // Path of the configuration file.
	ClusterFile  string // Path of the cluster file.
	Provision    bool   // Flag to start provisioning.
	DeleteAll    bool   // Flag to delete all VMs.
	Verbose      bool   // Flag for verbose output.
)

// rootCmd is the base command for kvmcli.
var rootCmd = &cobra.Command{
	Use:   "kvmcli",
	Short: "A CLI for managing KVM virtual machines",
	Long:  "A CLI similar to kubectl for creating, deleting, and managing KVM VMs.",
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Errorf("Error executing command: %v", err)
	}
}

func init() {
	rootCmd.AddCommand(CreateCmd)
	rootCmd.AddCommand(DeleteCmd)
	rootCmd.AddCommand(GetCmd)
	rootCmd.AddCommand(InitVMCmd)
}
