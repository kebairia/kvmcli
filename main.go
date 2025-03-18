package main

import (
	"github.com/kebairia/kvmcli/cmd"
	_ "github.com/kebairia/kvmcli/internal/database"
)

func main() {
	// Execute the CLI commands defined in the cmd package.
	cmd.Execute()
	// database.ConnectToMongo()
}
