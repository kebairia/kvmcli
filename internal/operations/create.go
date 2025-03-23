package operations

import (
	"github.com/kebairia/kvmcli/internal"
	"github.com/kebairia/kvmcli/internal/loader"
	"github.com/kebairia/kvmcli/internal/logger"
	"github.com/kebairia/kvmcli/internal/resources"
)

// CreateFromManifest loads a resource manifest from the given path,
// injects a libvirt connection where necessary, and creates the resources.
func CreateFromManifest(manifestPath string) error {
	conn, err := internal.InitConnection()
	if err != nil {
		logger.Log.Errorf("failed to establish libvirt connection: %v", err)
		return err
	}
	// Ensure that the libvirt connection is closed when the function exits.
	defer conn.Disconnect()

	resList, err := loader.LoadManifest(manifestPath)
	if err != nil {
		logger.Log.Errorf("failed to load configuration: %v", err)
		return err
	}

	for _, res := range resList {
		// If the resource requires a libvirt connection, inject it.
		if cs, ok := res.(resources.ClientSetter); ok {
			cs.SetConnection(conn)
		}
		// Create the resource and log any errors.
		if err := res.Create(); err != nil {
			logger.Log.Errorf("failed to create resource: %v", err)
			// Optionally, you could return the error here instead of continuing.
			continue
		}
	}

	return nil
}
