package main

import (
	// "fmt"

	"github.com/kebairia/kvmcli/cmd"
	// db "github.com/kebairia/kvmcli/internal/database"
	// log "github.com/kebairia/kvmcli/internal/logger"
)

func main() {
	// var err error
	// db.DB, err = db.InitDB()
	// if err != nil {
	// 	// log.Fatalf()
	// 	log.Errorf("DB initialization failed: %v", err)
	// }
	// fmt.Printf("database @: %p\n", db.DB)
	// Execute the CLI commands defined in the cmd package.
	cmd.Execute()
}
