package vms

import (
	"github.com/hashicorp/hcl/v2"
)

// VM describes a virtual machine definition.
type Config struct {
	Name      string            `hcl:"name,label"`
	Namespace string            `hcl:"namespace"`
	Image     string            `hcl:"image"`
	CPU       int               `hcl:"cpu"`
	Memory    int               `hcl:"memory"`
	Disk      string            `hcl:"disk,optional"`
	NetExpr   hcl.Expression    `hcl:"network,attr"` // raw HCL expression, e.g. network.homelab
	NetName   string            // resolved network name; filled by ResolveReferences
	StoreExpr hcl.Expression    `hcl:"store,attr"`
	Store     string            // resolved store name; filled by ResolveReferences
	MAC       string            `hcl:"mac,optional"`
	IP        string            `hcl:"ip,optional"`
	Labels    map[string]string `hcl:"labels,optional"`
}
