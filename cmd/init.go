package cmd

import (
	"github.com/spf13/cobra"
)

var InitVMCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a YAML template file with one virtual machine definition.",
	Run: func(cmd *cobra.Command, args []string) {
	},
}
