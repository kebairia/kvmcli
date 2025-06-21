package operations

import (
	"context"
	"fmt"
	"slices"
	"time"

	log "github.com/kebairia/kvmcli/internal/logger"
	"github.com/kebairia/kvmcli/internal/manifest"
	"github.com/kebairia/kvmcli/internal/resources"
)

func DeleteFromManifest(manifestPath string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	operator, err := NewOperator(ctx)
	if err != nil {
		return fmt.Errorf("failed to create operator: %w", err)
	}
	defer operator.Close()

	resources, err := manifest.Load(manifestPath)
	if err != nil {
		return fmt.Errorf("failed to load manifest %q: %w", manifestPath, err)
	}

	slices.Reverse(resources)
	for _, resource := range resources {
		if err := operator.Delete(resource); err != nil {
			log.Errorf("failed to delete resource: %v", err)
			continue
		}
	}

	return nil
}

// Delete removes the given Resource.
func (o *Operator) Delete(r resources.Resource) error {
	return r.Delete()
}
