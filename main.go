package main

import (
	"github.com/kebairia/kvmcli/cmd"
	db "github.com/kebairia/kvmcli/internal/database"
	"github.com/kebairia/kvmcli/internal/logger"
)

func main() {
	var err error
	db.DB, err = db.InitDB()
	if err != nil {
		logger.Log.Fatalf("DB initialization failed: %v", err)
	}
	// Execute the CLI commands defined in the cmd package.
	cmd.Execute()
}
