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
	"github.com/kebairia/kvmcli/internal/network"
	"github.com/kebairia/kvmcli/internal/resources"
	"github.com/kebairia/kvmcli/internal/store"
	"github.com/kebairia/kvmcli/internal/vms"
	"github.com/zclconf/go-cty/cty"
)

// Config represents a complete kvmcli configuration file.
type Config struct {
	Networks []network.Network   `hcl:"network,block"`
	VMs      []vms.VM            `hcl:"vm,block"`
	Stores   []store.StoreConfig `hcl:"store,block"`
	Clusters []Cluster           `hcl:"cluster,block"`
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
	if err := cfg.ResolveReferences(); err != nil {
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
		netRes, err := network.NewVirtualNetwork(
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

func (cfg *Config) ResolveReferences() error {
	if cfg == nil {
		return fmt.Errorf("config is nil")
	}

	networksByName, err := cfg.indexNetworks()
	if err != nil {
		return err
	}

	evalCtx := cfg.evalContextForNetworks()

	for i := range cfg.VMs {
		if err := resolveVMNetwork(&cfg.VMs[i], networksByName, evalCtx); err != nil {
			return err
		}
	}
	// for _, vm := range cfg.VMs {
	// 	if err := resolveVMNetwork(vm, networksByName, evalCtx); err != nil {
	// 		return err
	// 	}
	// }

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
	vm *vms.VM,
	networks map[string]struct{},
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
	if _, ok := networks[netName]; !ok {
		return fmt.Errorf("vm %q: unknown network %q", vm.Name, netName)
	}

	vm.NetName = netName
	return nil
}

func (cfg *Config) evalContextForNetworks() *hcl.EvalContext {
	attrs := make(map[string]cty.Value, len(cfg.Networks))

	for _, n := range cfg.Networks {
		attrs[n.Name] = cty.StringVal(n.Name)
	}

	return &hcl.EvalContext{
		Variables: map[string]cty.Value{
			"network": cty.ObjectVal(attrs),
		},
	}
}
