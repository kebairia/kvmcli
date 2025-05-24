package operations

import (
	"github.com/kebairia/kvmcli/internal/loader"
	"github.com/kebairia/kvmcli/internal/logger"
)

func DeleteFromManifest(manifestPath string) error {
	// ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	// defer cancel()
	operator, err := NewOperator()
	if err != nil {
		return err
	}
	defer operator.Close()
	// resources := manifest.Load(manifestPath)
	resources, err := loader.LoadManifest(manifestPath)
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
