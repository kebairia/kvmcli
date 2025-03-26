package operations

import (
	"fmt"

	"github.com/kebairia/kvmcli/internal"
	"github.com/kebairia/kvmcli/internal/loader"
	"github.com/kebairia/kvmcli/internal/logger"
	"github.com/kebairia/kvmcli/internal/resources"
)

// NOTE: I need to use DeleteMany for mongodb when I delete from a manifest

func DeleteFromManifest(manifestPath string) error {
	conn, err := internal.InitConnection()
	if err != nil {
		return fmt.Errorf("failed to establish libvirt connection: %w", err)
	}

	// Ensure that the libvirt connection is closed when the function exits.
	defer conn.Disconnect()

	resList, err := loader.LoadManifest(manifestPath)
	if err != nil {
		logger.Log.Errorf("failed to Load config file: %v", err)
	}

	for _, res := range resList {
		if cs, ok := res.(resources.ClientSetter); ok {
			cs.SetConnection(conn)
		}
		if err := res.Delete(); err != nil {

			logger.Log.Errorf("failed to delete resource: %v", err)
			// Optionally, you could return the error here instead of continuing.
			continue

		}
	}

	return nil
}
