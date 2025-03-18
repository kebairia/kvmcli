package cmd

import (
	"github.com/kebairia/kvmcli/internal/logger"
	"github.com/kebairia/kvmcli/internal/operations/vms"
	"github.com/spf13/cobra"
)

var DeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete VM(s) from a configuration file",
	Run: func(cmd *cobra.Command, args []string) {
		if ConfigFile == "" {
			logger.Log.Fatalf("Configuration file is required (-f flag)")
		}
		// Call your delete operation with the provided file.
		// op.DestroyFromFile(ConfigFile)
		vms.DestroyFromFile(ConfigFile)
	},
}

func init() {
	DeleteCmd.Flags().
		StringVarP(&ConfigFile, "file", "f", "", "Configuration file for the VM(s) to delete")
}
