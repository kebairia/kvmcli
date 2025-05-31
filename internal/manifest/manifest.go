package manifest

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

const (
	KindStore          = "Store"
	KindNetwork        = "Network"
	KindVirtualMachine = "VirtualMachine"
)

// unmarshalStore decodes raw YAML into a *store.Store
func unmarshalStore(raw []byte) (resources.Resource, error) {
	var s store.Store
	if err := yaml.Unmarshal(raw, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

// unmarshalNetwork decodes raw YAML into a *network.VirtualNetwork
func unmarshalNetwork(raw []byte) (resources.Resource, error) {
	var n network.VirtualNetwork
	if err := yaml.Unmarshal(raw, &n); err != nil {
		return nil, err
	}
	return &n, nil
}

// unmarshalVirtualMachine decodes raw YAML into a *vms.VirtualMachine
func unmarshalVirtualMachine(raw []byte) (resources.Resource, error) {
	var m vms.VirtualMachine
	if err := yaml.Unmarshal(raw, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

// Load reads all YAML documents from the file at manifestPath, decodes each
// into its appropriate Resource struct, then returns a slice sorted in this
// order: [Stores | Networks | VirtualMachines]. If any document has an
// unknown kind, or fails to decode, Load returns an error immediately.
func Load(manifestPath string) ([]resources.Resource, error) {
	file, err := os.Open(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open manifest %q: %w", manifestPath, err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)

	var (
		stores   []resources.Resource
		networks []resources.Resource
		vmsList  []resources.Resource
	)

	for {
		// 1) Decode one document into a generic map
		var rawDoc map[string]any
		if err := decoder.Decode(&rawDoc); err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("failed to decode YAML: %w", err)
		}

		// 2) Extract the \"kind\" field
		kindVal, ok := rawDoc["kind"]
		if !ok {
			return nil, fmt.Errorf("document missing \"kind\" field")
		}
		kind, ok := kindVal.(string)
		if !ok {
			return nil, fmt.Errorf("invalid \"kind\" value: %v", kindVal)
		}

		// 3) Re-marshal just this document back into bytes
		rawBytes, err := yaml.Marshal(rawDoc)
		if err != nil {
			return nil, fmt.Errorf("failed to re-marshal %q resource: %w", kind, err)
		}

		// 4) Unmarshal into the correct struct based on kind
		switch kind {
		case KindStore:
			store, err := unmarshalStore(rawBytes)
			if err != nil {
				return nil, fmt.Errorf("store unmarshal error: %w", err)
			}
			stores = append(stores, store)

		case KindNetwork:
			network, err := unmarshalNetwork(rawBytes)
			if err != nil {
				return nil, fmt.Errorf("network unmarshal error: %w", err)
			}

			networks = append(networks, network)

		case KindVirtualMachine:
			vm, err := unmarshalVirtualMachine(rawBytes)
			if err != nil {
				return nil, fmt.Errorf("virtualmachine unmarshal error: %w", err)
			}
			vmsList = append(vmsList, vm)

		default:
			return nil, fmt.Errorf("unknown resource kind: %q", kind)
		}
	}

	// 5) Concatenate in the order: stores → networks → vms
	sorted := make([]resources.Resource, 0,
		len(stores)+len(networks)+len(vmsList),
	)
	sorted = append(sorted, stores...)
	sorted = append(sorted, networks...)
	sorted = append(sorted, vmsList...)

	return sorted, nil
}
