package config

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/digitalocean/go-libvirt"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/kebairia/kvmcli/internal/database"
	"github.com/kebairia/kvmcli/internal/network"
	"github.com/kebairia/kvmcli/internal/resources"
	"github.com/kebairia/kvmcli/internal/store"
	"github.com/kebairia/kvmcli/internal/vms"
	"github.com/zclconf/go-cty/cty"
)

// Config represents a complete kvmcli configuration file.
type Config struct {
	Networks []network.Config `hcl:"network,block"`
	VMs      []vms.Config     `hcl:"vm,block"`
	Stores   []store.Config   `hcl:"store,block"`
	Clusters []Cluster        `hcl:"cluster,block"`
	Data     []DataResource   `hcl:"data,block"`
}

type DataResource struct {
	Type string `hcl:"type,label"`
	Name string `hcl:"name,label"`
}

// Cluster describes a logical grouping of VMs.
type Cluster struct {
	Name      string            `hcl:"name,label"`
	VMExprs   hcl.Expression    `hcl:"vms,attr"` // List of VM references
	VMNames   []string          // Resolved VM names
	Labels    map[string]string `hcl:"labels,optional"`
	Lifecycle *Lifecycle        `hcl:"lifecycle,block"`
}

type Lifecycle struct {
	StartOrder []string `hcl:"start_order,optional"`
	StopOrder  []string `hcl:"stop_order,optional"`
}

// Load parses and decodes the configuration file at the given path.
func Load(
	path string,
	ctx context.Context,
	db *sql.DB,
	conn *libvirt.Libvirt,
) ([]resources.Resource, error) {
	src, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config %q: %w", path, err)
	}

	parser := hclparse.NewParser()
	file, diags := parser.ParseHCL(src, path)
	if diags.HasErrors() {
		return nil, fmt.Errorf("parse hcl %q: %w", path, diags)
	}

	var cfg Config
	if diags := gohcl.DecodeBody(file.Body, nil, &cfg); diags.HasErrors() {
		return nil, fmt.Errorf("decode hcl %q: %w", path, diags)
	}
	if err := cfg.ResolveReferences(ctx, db); err != nil {
		return nil, err
	}

	var networks []resources.Resource
	var stores []resources.Resource
	var vmsList []resources.Resource
	// for _, n := range cfg.Networks {
	// 	networks = append(networks, n)
	// }

	for _, v := range cfg.VMs {
		vmRes, err := vms.NewVirtualMachine(
			v,
			vms.WithContext(ctx),
			vms.WithDatabaseConnection(db),
			vms.WithLibvirtConnection(conn),
		)
		if err != nil {
			return nil, err
		}

		vmsList = append(vmsList, vmRes)
	}
	for _, n := range cfg.Networks {
		netRes, err := network.NewNetwork(
			n,
			network.WithContext(ctx),
			network.WithDatabaseConnection(db),
			network.WithLibvirtConnection(conn),
		)
		if err != nil {
			return nil, err
		}

		networks = append(networks, netRes)
	}

	for _, s := range cfg.Stores {
		stRes, err := store.NewStore(
			s,
			store.WithContext(ctx),
			store.WithDatabaseConnection(db),
		)
		if err != nil {
			return nil, err
		}

		stores = append(stores, stRes)
	}

	sorted := make([]resources.Resource, 0,
		len(stores)+len(vmsList),
	)
	sorted = append(sorted, stores...)
	sorted = append(sorted, networks...)
	sorted = append(sorted, vmsList...)

	return sorted, nil
}

func (cfg *Config) ResolveReferences(ctx context.Context, db *sql.DB) error {
	if cfg == nil {
		return fmt.Errorf("config is nil")
	}

	// 1. Build evaluation context for existing config blocks
	networksByName, err := cfg.indexNetworks()
	if err != nil {
		return err
	}
	storesByName, err := cfg.indexStores()
	if err != nil {
		return err
	}

	// 2. Process data blocks (data "store" "..." {})
	// We'll verify they exist in the DB, then add them to our evaluation maps.
	dataStores := make(map[string]struct{})
	dataNetworks := make(map[string]struct{})

	for _, d := range cfg.Data {
		switch d.Type {
		case "store":
			// Check DB
			if _, err := database.GetStoreIDByName(ctx, db, d.Name); err != nil {
				return fmt.Errorf("data.store.%s: %w", d.Name, err)
			}
			dataStores[d.Name] = struct{}{}
		case "network":
			if _, err := database.GetNetworkIDByName(ctx, db, d.Name); err != nil {
				return fmt.Errorf("data.network.%s: %w", d.Name, err)
			}
			dataNetworks[d.Name] = struct{}{}
		default:
			return fmt.Errorf("unknown data type %q (supported: store, network)", d.Type)
		}
	}

	// 3. Construct the shared EvalContext
	evalCtx := cfg.evalContext(networksByName, storesByName, dataNetworks, dataStores)

	// 4. Resolve VM references
	for i := range cfg.VMs {
		vm := &cfg.VMs[i]

		// Resolve Network
		if err := resolveVMNetwork(vm, networksByName, dataNetworks, evalCtx); err != nil {
			return err
		}

		// Resolve Store
		if err := resolveVMStore(vm, storesByName, dataStores, evalCtx); err != nil {
			return err
		}
	}

	return nil
}

func (cfg *Config) indexStores() (map[string]struct{}, error) {
	stores := make(map[string]struct{}, len(cfg.Stores))
	for _, s := range cfg.Stores {
		if s.Name == "" {
			return nil, fmt.Errorf("store with empty name")
		}
		if _, exists := stores[s.Name]; exists {
			return nil, fmt.Errorf("duplicate store %q", s.Name)
		}
		stores[s.Name] = struct{}{}
	}
	return stores, nil
}

// evalContext builds a single *hcl.EvalContext containing variables:
//
//	network.<name>
//	store.<name>
//	data.network.<name>
//	data.store.<name>
func (cfg *Config) evalContext(
	networks, stores, dataNetworks, dataStores map[string]struct{},
) *hcl.EvalContext {
	// Objects for 'network' and 'store'
	netMap := make(map[string]cty.Value)
	for n := range networks {
		netMap[n] = cty.StringVal(n)
	}
	storeMap := make(map[string]cty.Value)
	for s := range stores {
		storeMap[s] = cty.StringVal(s)
	}

	// Objects for 'data.network' and 'data.store'
	// In HCL, data variables are usually top-level `data` object containing types.
	dNetMap := make(map[string]cty.Value)
	for n := range dataNetworks {
		dNetMap[n] = cty.StringVal(n)
	}
	dStoreMap := make(map[string]cty.Value)
	for s := range dataStores {
		dStoreMap[s] = cty.StringVal(s)
	}

	dataObj := cty.ObjectVal(map[string]cty.Value{
		"network": cty.ObjectVal(dNetMap),
		"store":   cty.ObjectVal(dStoreMap),
	})

	return &hcl.EvalContext{
		Variables: map[string]cty.Value{
			"network": cty.ObjectVal(netMap),
			"store":   cty.ObjectVal(storeMap),
			"data":    dataObj,
		},
	}
}

// resolveVMStore resolves the `store = ...` expression on a VM.
func resolveVMStore(
	vm *vms.Config,
	configStores map[string]struct{},
	dataStores map[string]struct{},
	ctx *hcl.EvalContext,
) error {
	if vm.StoreExpr == nil {
		return fmt.Errorf("vm %q: missing required attribute \"store\"", vm.Name)
	}

	val, diags := vm.StoreExpr.Value(ctx)
	if diags.HasErrors() {
		return fmt.Errorf("vm %q: invalid store expression: %w", vm.Name, diags)
	}
	if val.Type() != cty.String {
		return fmt.Errorf(
			"vm %q: store must be a string (reference), got %s",
			vm.Name,
			val.Type().FriendlyName(),
		)
	}

	storeName := val.AsString()

	// Check if present in locally defined stores OR data stores
	_, local := configStores[storeName]
	_, data := dataStores[storeName]

	if !local && !data {
		return fmt.Errorf("vm %q: unknown store %q", vm.Name, storeName)
	}

	vm.Store = storeName
	return nil
}

// NOTE: this is basically check the validity of network keywords in the config
func (cfg *Config) indexNetworks() (map[string]struct{}, error) {
	networks := make(map[string]struct{}, len(cfg.Networks))

	for _, n := range cfg.Networks {
		if n.Name == "" {
			return nil, fmt.Errorf("network with empty name")
		}
		if _, exists := networks[n.Name]; exists {
			return nil, fmt.Errorf("duplicate network %q", n.Name)
		}
		networks[n.Name] = struct{}{}
	}

	return networks, nil
}

// NOTE: this is resolve the network name from the network experession
func resolveVMNetwork(
	vm *vms.Config,
	networks map[string]struct{},
	dataNetworks map[string]struct{},
	evalCtx *hcl.EvalContext,
) error {
	if vm.NetExpr == nil {
		return fmt.Errorf("vm %q: missing required attribute \"net\"", vm.Name)
	}

	val, diags := vm.NetExpr.Value(evalCtx)
	if diags.HasErrors() {
		return fmt.Errorf("vm %q: invalid net expression: %w", vm.Name, diags)
	}

	if val.Type() != cty.String {
		return fmt.Errorf(
			"vm %q: net must be a string, got %s",
			vm.Name,
			val.Type().FriendlyName(),
		)
	}

	netName := val.AsString()
	// Check local or data
	_, local := networks[netName]
	_, data := dataNetworks[netName]
	if !local && !data {
		return fmt.Errorf("vm %q: unknown network %q", vm.Name, netName)
	}

	vm.NetName = netName
	return nil
}
