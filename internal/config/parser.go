package config

import (
	"fmt"
	"log"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
)

// ==================================================
// 1. BASIC: Reference existing networks
// ==================================================

// Config represents the entire HCL file
type Config struct {
	Networks     []*Network     `hcl:"network,block"`
	DataNetworks []*DataNetwork `hcl:"data,block"`
	VMs          []*VM          `hcl:"vm,block"`
	Clusters     []*Cluster     `hcl:"cluster,block"`
	Remain       hcl.Body       `hcl:",remain"`
}

// DataNetwork represents a reference to existing network
type DataNetwork struct {
	Type   string   `hcl:"type,label"` // "network"
	Name   string   `hcl:"name,label"` // network name
	Remain hcl.Body `hcl:",remain"`
}

// Network represents a new network definition
type Network struct {
	Name   string   `hcl:"name,label"`
	CIDR   string   `hcl:"cidr"`
	Mode   string   `hcl:"mode,optional"`
	Bridge string   `hcl:"bridge,optional"`
	Remain hcl.Body `hcl:",remain"`
}

// VM represents a virtual machine
type VM struct {
	Name       string         `hcl:"name,label"`
	Image      string         `hcl:"image"`
	CPU        int            `hcl:"cpu"`
	Memory     string         `hcl:"memory"`
	Disk       string         `hcl:"disk"`
	NetRef     hcl.Expression `hcl:"net"` // This can be a reference!
	MAC        string         `hcl:"mac,optional"`
	IP         string         `hcl:"ip,optional"`
	ClusterTag string         `hcl:"cluster,optional"`
	Role       string         `hcl:"role,optional"`
	Remain     hcl.Body       `hcl:",remain"`
}

// Cluster represents a logical grouping
type Cluster struct {
	Name   string            `hcl:"name,label"`
	VMs    []hcl.Expression  `hcl:"vms"` // References to VMs
	Labels map[string]string `hcl:"labels,optional"`
	Remain hcl.Body          `hcl:",remain"`
}

// Load parses an HCL configuration file and returns a Config struct
func Load(path string) (*Config, error) {
	parser := hclparse.NewParser()

	// Parse the HCL file
	file, diags := parser.ParseHCLFile(path)
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse HCL file: %w", diags)
	}

	var config Config

	// First pass: Decode basic structure without evaluation context
	// This collects all blocks (networks, VMs, clusters, etc.)
	confDiags := gohcl.DecodeBody(file.Body, nil, &config)
	if confDiags.HasErrors() {
		// Log warnings but continue - some errors are expected (unresolved references)
		log.Printf("First pass diagnostics: %s", confDiags)
	}

	// Second pass: Build evaluation context with all resources
	// ctx := buildEvalContext(&config)

	// Re-decode with evaluation context to resolve references
	confDiags = gohcl.DecodeBody(file.Body, nil, &config)
	if confDiags.HasErrors() {
		return nil, fmt.Errorf("failed to decode HCL config: %w", confDiags)
	}

	// Validate the configuration
	// if err := config.Validate(); err != nil {
	// 	return nil, fmt.Errorf("invalid configuration: %w", err)
	// }

	return &config, nil
}
