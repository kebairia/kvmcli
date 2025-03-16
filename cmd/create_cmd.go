package cmd

import (
	"github.com/kebairia/kvmcli/internal/logger"
	op "github.com/kebairia/kvmcli/internal/operations"
	"github.com/spf13/cobra"
)

var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create VM(s) from a configuration file",
	Run: func(cmd *cobra.Command, args []string) {
		if ConfigFile == "" {
			logger.Log.Fatalf("Configuration file is required (-f flag)")
		}
		op.ProvisionVMs(ConfigFile)
	},
}

func init() {
	CreateCmd.Flags().StringVarP(&ConfigFile, "file", "f", "", "Configuration file for the VM(s)")
}
