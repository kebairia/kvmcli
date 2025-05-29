package operations

import (
	"context"
	"time"

	"github.com/kebairia/kvmcli/internal/logger"
	"github.com/kebairia/kvmcli/internal/manifest"
)

func DeleteFromManifest(manifestPath string) error {
	// log := logger.Logger
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
		if err := operator.Delete(resource); err != nil {
			logger.Log.Errorf("failed to delete resource: %v\n", err)
			continue
		}
	}

	return nil
}
