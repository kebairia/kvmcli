package cmd

import (
	"github.com/spf13/cobra"
)

// Root command
var RootCmd = &cobra.Command{
	Use:   "kvmcli",
	Short: "A CLI to manage KVM virtual machines using YAML files",
}
