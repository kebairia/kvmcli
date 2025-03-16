package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Global variables to store flag values
var (
	ConfigFile  string // To store the path of the config file
	ClusterFile string // To store the path of the cluster file
	Provision   bool   // To determine if provisioning should start
	Verbose     bool   // To determine if output should be in verbose
)

func Execute() {
	// Create the root command
	rootCmd := &cobra.Command{
		Use:   "kvmcli",                                 // How the command is called from the terminal
		Short: "kvmcli is a tool to manage VM clusters", // A short description
		Long: `
    __                        ___ 
   / /___   ______ ___  _____/ (_)
  / //_/ | / / __ '__ \/ ___/ / / 
 / ,<  | |/ / / / / / / /__/ / /  
/_/|_| |___/_/ /_/ /_/\___/_/_/

kvmcli is a command line tool that helps you manage and provision 
virtual machine clusters using configuration files. `,
		// The Run function executes when the root command is called without any subcommands.
		Run: func(cmd *cobra.Command, args []string) {
			// Here you could call further functions to load configuration, read cluster file, and provision VMs.
		},
	}

	// Define flags for the root command using PersistentFlags.
	// Persistent flags are available to the command and all its subcommands.
	// If you wanted flags only for this command, you could use Flags() instead.

	// 1. --config or -c: specifies the path to the config file.
	rootCmd.PersistentFlags().
		StringVarP(&ConfigFile, "config", "c", "config.toml", "Path to the configuration file")

	// 2. --file or -f: specifies the file that defines your cluster info.
	rootCmd.PersistentFlags().
		StringVarP(&ClusterFile, "file", "f", "servers.yaml", "Path to the cluster file")

	// 3. --provision or -p: a boolean flag to start provisioning VMs.
	rootCmd.PersistentFlags().
		BoolVarP(&Provision, "provision", "p", false, "Start provisioning VMs")

	// 4. --verbose or -v: a boolean flag to start provisioning VMs.
	rootCmd.PersistentFlags().
		BoolVarP(&Verbose, "verbose", "v", false, "Output in verbose mode")

	// Execute the command line application.
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
