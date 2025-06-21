package manifest

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/digitalocean/go-libvirt"
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
type Manifest interface{}

func Load(
	manifestPath string,
	ctx context.Context,
	db *sql.DB,
	conn *libvirt.Libvirt,
) ([]resources.Resource, error) {
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
			var cfg store.StoreConfig
			if err := doc.Decode(&cfg); err != nil {
				return nil, fmt.Errorf("unmarshal Store: %w", err)
			}
			stRes, err := store.NewStore(cfg,
				store.WithContext(ctx),
				store.WithDatabaseConnection(db),
			)
			if err != nil {
				return nil, err
			}
			stores = append(stores, stRes)
		case KindNetwork:
			// var net network.VirtualNetwork
			// if err := doc.Decode(&net); err != nil {
			// 	return nil, fmt.Errorf("unmarshal Network: %w", err)
			// }
			var cfg network.VirtualNetworkConfig
			if err := doc.Decode(&cfg); err != nil {
				return nil, fmt.Errorf("unmarshal NetworkConfig: %w", err)
			}
			netRes, err := network.NewVirtualNetwork(cfg,
				network.WithContext(ctx),
				network.WithDatabaseConnection(db),
				network.WithLibvirtConnection(conn),
			)
			if err != nil {
				return nil, err
			}
			networks = append(networks, netRes)

			// networks = append(networks, &net)
		case KindVirtualMachine:
			var cfg vms.VirtualMachineConfig
			if err := doc.Decode(&cfg); err != nil {
				return nil, fmt.Errorf("unmarshal VMConfig: %w", err)
			}
			vmRes, err := vms.NewVirtualMachine(
				cfg,
				vms.WithContext(ctx),
				vms.WithDatabaseConnection(db),
				vms.WithLibvirtConnection(conn),
			)
			if err != nil {
				return nil, err
			}
			vmsList = append(vmsList, vmRes)
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

// func InitResources(configs []ResourceConfig) ([]resources.Resource, error) {
//     var (
//         stores   []resources.Resource
//         nets     []resources.Resource
//         vmsList  []resources.Resource
//     )
//
//     for _, rc := range configs {
//         switch rc.Kind {
//         case KindStore:
//             st, err := store.NewStore(*rc.Store, db)
//             if err != nil { return nil, err }
//             stores = append(stores, st)
//
//         case KindNetwork:
//             netRes, err := network.NewVirtualNetwork(*rc.Network, db, conn)
//             if err != nil { return nil, err }
//             nets = append(nets, netRes)
//
//         case KindVirtualMachine:
//             vmRes, err := vms.NewVirtualMachine(
//                 *rc.VM,
//                 vms.WithContext(ctx),
//                 vms.WithDatabaseConnection(db),
//                 vms.WithLibvirtConnection(conn),
//             )
//             if err != nil { return nil, err }
//             vmsList = append(vmsList, vmRes)
//         }
//     }
//
//     // “Flatten” them in the exact order you want
//     out := make([]resources.Resource, 0,
//         len(stores)+len(nets)+len(vmsList),
//     )
//     out = append(out, stores...)
//     out = append(out, nets...)
//     out = append(out, vmsList...)
//     return out, nil
// }
