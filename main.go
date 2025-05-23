package main

import (
	"github.com/kebairia/kvmcli/cmd"
	databasesql "github.com/kebairia/kvmcli/internal/database-sql"
	"github.com/kebairia/kvmcli/internal/logger"
)

func main() {
	var err error
	databasesql.DB, err = databasesql.InitDB()
	if err != nil {
		logger.Log.Fatalf("DB initialization failed: %v", err)
	}
	// Execute the CLI commands defined in the cmd package.
	cmd.Execute()
}
