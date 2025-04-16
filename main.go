package main

import (
	"github.com/kebairia/kvmcli/cmd"
	_ "github.com/kebairia/kvmcli/internal/database"
	databasesql "github.com/kebairia/kvmcli/internal/database-sql"
)

func main() {
	databasesql.InitDB()
	// Execute the CLI commands defined in the cmd package.
	cmd.Execute()
}
