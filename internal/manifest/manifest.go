package manifest

import (
	"errors"
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

type kindOnly struct {
	Kind string `yaml:"kind"`
}

func Load(manifestPath string) ([]resources.Resource, error) {
	f, err := os.Open(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("open manifest: %w", err)
	}
	defer f.Close()
	decoder := yaml.NewDecoder(f)

	var (
		stores   []resources.Resource
		networks []resources.Resource
		vmsList  []resources.Resource
	)
	for {
		var doc yaml.Node
		if err := decoder.Decode(&doc); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, fmt.Errorf("decode document: %w", err)
		}
		// 2) Decode just Kind
		var meta kindOnly
		if err := doc.Decode(&meta); err != nil {
			return nil, fmt.Errorf("decode kind: %w", err)
		}

		switch meta.Kind {
		case KindStore:
			var st store.Store
			if err := doc.Decode(&st); err != nil {
				return nil, fmt.Errorf("unmarshal Store: %w", err)
			}
			stores = append(stores, &st)
		case KindNetwork:
			var net network.VirtualNetwork
			if err := doc.Decode(&net); err != nil {
				return nil, fmt.Errorf("unmarshal Network: %w", err)
			}
			networks = append(networks, &net)
		case KindVirtualMachine:
			var vm vms.VirtualMachine
			if err := doc.Decode(&vm); err != nil {
				return nil, fmt.Errorf("unmarshal VirtualMachine: %w", err)
			}
			vmsList = append(vmsList, &vm)
		default:
			return nil, fmt.Errorf("unknown kind %q", meta.Kind)

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
