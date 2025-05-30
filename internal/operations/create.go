package operations

import (
	"context"
	"time"

	"github.com/kebairia/kvmcli/internal/logger"
	"github.com/kebairia/kvmcli/internal/manifest"
	"github.com/kebairia/kvmcli/internal/resources"
)

// TODO: Create a context with a timeout for the operations.
//		   This function is responsible for creating resources defined in a manifest file.
//		   1. Create a new operator
//	     2. Load the manifest file, and extract the resources.
//	     3. Loop through the resources and create them one by one.
//			 !. the operator has a connection to the libvirt daemon.
//			 !. Create/Delete/Update the resources as needed.

// IDEA: use go routines to create the resources in parallel.

// NOTICE: using go routines has an issue, because sometimes I need to create network resources
//       before creating the VMs.

func CreateFromManifest(manifestPath string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	operator, err := NewOperator(ctx)
	if err != nil {
		return err
	}
	defer operator.Close()
	resources, err := manifest.Load(manifestPath)
	if err != nil {
		logger.Log.Errorf("failed to load configuration: %v", err)
		return err
	}
	for _, resource := range resources {
		if err := operator.Create(resource); err != nil {
			logger.Log.Errorf("failed to create resource: %v\n", err)
			continue
		}
	}

	return nil
}

// Create provisions the given Resource.
func (o *Operator) Create(r resources.Resource) error {
	o.SetConnection(r)
	return r.Create() // assumes your interface takes (ctx, db)
}
