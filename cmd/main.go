package main

import (
	"flag"

	op "github.com/kebairia/kvmcli/internal/operations"
)

func main() {
	// Declare the config flag with a default value ("servers.yaml") and a short description.
	configPath := flag.String("config", "servers.yaml", "Path to the YAML configuration file")
	// Declare a boolean flag (e.g., --verbose) if you need more logging or debugging info.
	// verbose := flag.Bool("verbose", false, "Enable verbose logging")

	// Parse all the flags. This reads os.Args and sets the values of the flag variables.
	flag.Parse()

	op.ProvisionVMs(*configPath)
}
