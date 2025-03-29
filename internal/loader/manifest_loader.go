package loader

import (
	"fmt"
	"io"
	"os"

	"github.com/kebairia/kvmcli/internal/network"
	"github.com/kebairia/kvmcli/internal/resources"
	"github.com/kebairia/kvmcli/internal/store"
	"github.com/kebairia/kvmcli/internal/vms"
	"gopkg.in/yaml.v3"
)

func LoadManifest(manifestPath string) ([]resources.Resource, error) {
	var resourceList []resources.Resource
	file, err := os.Open(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	for {
		var rawMap map[string]any
		if err := decoder.Decode(&rawMap); err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("failed to decode YAML: %w", err)
		}
		kind, ok := rawMap["kind"].(string)
		if !ok {
			return nil, fmt.Errorf("missing or invalide 'kind' field")
		}
		rawBytes, err := yaml.Marshal(rawMap)
		if err != nil {
			return nil, fmt.Errorf("failed to re-marshal resource: %w", err)
		}
		switch kind {
		// Virtual Machines
		case "VirtualMachine":
			var vm vms.VirtualMachine
			if err := yaml.Unmarshal(rawBytes, &vm); err != nil {
				return nil, fmt.Errorf("failed to unmarshal VirtualMachine: %w", err)
			}
			resourceList = append(resourceList, &vm)
		// Virtual Networks
		case "Network":
			var net network.VirtualNetwork
			if err := yaml.Unmarshal(rawBytes, &net); err != nil {
				return nil, fmt.Errorf("failed to unmarshal VirtualNetwork: %w", err)
			}
			resourceList = append(resourceList, &net)
			// Stores
		case "Store":
			var store store.Store
			if err := yaml.Unmarshal(rawBytes, &store); err != nil {
				return nil, fmt.Errorf("failed to unmarshal Store: %w", err)
			}
			resourceList = append(resourceList, &store)

			// Default action
		default:
			return nil, fmt.Errorf("unknown kind: %s", kind)
		}
	}
	return resourceList, nil
}
