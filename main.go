package main

import (
	"log"

	"github.com/kebairia/kvmcli/cmd"
	db "github.com/kebairia/kvmcli/internal/database"
)

func main() {
	var err error
	db.DB, err = db.InitDB()
	if err != nil {
		log.Fatalf("DB initialization failed: %v", err)
	}
	// Execute the CLI commands defined in the cmd package.
	cmd.Execute()
}
