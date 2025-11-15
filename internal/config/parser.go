package config

import (
	"fmt"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/zclconf/go-cty/cty"
)

// Config represents a complete kvmcli configuration file.
type Config struct {
	Networks []Network `hcl:"network,block"`
	VMs      []VM      `hcl:"vm,block"`
}

// Network describes a virtual network definition.
type Network struct {
	Name   string `hcl:"name,label"`
	CIDR   string `hcl:"cidr"`
	Bridge string `hcl:"bridge"`
	Mode   string `hcl:"mode"`

	DHCP   *DHCP             `hcl:"dhcp,block"`
	Labels map[string]string `hcl:"labels,attr"`
}

// DHCP describes the dhcp block inside a network.
type DHCP struct {
	Range string `hcl:"range"`
}

// VM describes a virtual machine definition.
type VM struct {
	Name    string            `hcl:"name,label"`
	Image   string            `hcl:"image"`
	CPU     int               `hcl:"cpu"`
	Memory  string            `hcl:"memory"`
	Disk    string            `hcl:"disk"`
	NetExpr hcl.Expression    `hcl:"net,attr"` // raw HCL expression, e.g. network.homelab
	NetName string            // resolved network name; filled by ResolveReferences
	MAC     string            `hcl:"mac"`
	IP      string            `hcl:"ip"`
	Labels  map[string]string `hcl:"labels,attr"`
}

// Load parses and decodes the configuration file at the given path.
func Load(path string) (*Config, error) {
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

	return &cfg, nil
}

// ResolveReferences evaluates cross-resource references in the config
// (currently only VM.net) and performs basic validation. On success,
// the resolved network name is stored in VM.NetName.
func (cfg *Config) ResolveReferences() error {
	if cfg == nil {
		return fmt.Errorf("nil config")
	}

	// 1. Index networks by name and validate.
	networksByName := make(map[string]struct{}, len(cfg.Networks))
	for _, n := range cfg.Networks {
		if n.Name == "" {
			return fmt.Errorf("network with empty name defined")
		}
		if _, exists := networksByName[n.Name]; exists {
			return fmt.Errorf("duplicate network %q", n.Name)
		}
		networksByName[n.Name] = struct{}{}
	}

	// 2. Build evaluation context for HCL expressions.
	evalCtx := cfg.evalContextForNetworks()

	// 3. Resolve VM.net expressions.
	for i := range cfg.VMs {
		vm := &cfg.VMs[i]

		if vm.NetExpr == nil {
			return fmt.Errorf("vm %q: missing required attribute \"net\"", vm.Name)
		}

		val, diags := vm.NetExpr.Value(evalCtx)
		if diags.HasErrors() {
			return fmt.Errorf("vm %q: evaluating net: %w", vm.Name, diags)
		}

		if !val.Type().Equals(cty.String) {
			return fmt.Errorf("vm %q: net must evaluate to string, got %s",
				vm.Name, val.Type().FriendlyName())
		}

		netName := val.AsString()
		if _, ok := networksByName[netName]; !ok {
			return fmt.Errorf("vm %q: references unknown network %q", vm.Name, netName)
		}

		vm.NetName = netName
	}

	return nil
}

// evalContextForNetworks builds an evaluation context exposing networks for
// expressions like `net = network.homelab`.
func (cfg *Config) evalContextForNetworks() *hcl.EvalContext {
	networkAttrs := make(map[string]cty.Value, len(cfg.Networks))

	for _, n := range cfg.Networks {
		// For now the attribute value is just the network name.
		// Later you could expose more info (cidr, bridge, ...).
		networkAttrs[n.Name] = cty.StringVal(n.Name)
	}

	return &hcl.EvalContext{
		Variables: map[string]cty.Value{
			"network": cty.ObjectVal(networkAttrs),
		},
	}
}
