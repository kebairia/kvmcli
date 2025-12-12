package network

// IDEA: the ip <=>  mac address mapping done here in Virtual Network declaration
// What I need is whenever I create a new virtual machine with a static ip, I need to update
// my virtual network declaration to add the ip <=> mac address mapping

type Config struct {
	Name       string `hcl:"name,label"`
	Namespace  string `hcl:"namespace"`
	CIDR       string `hcl:"cidr,optional"`
	NetAddress string `hcl:"netaddress,optional"`
	NetMask    string `hcl:"netmask,optional"`
	Bridge     string `hcl:"bridge,optional"`
	Mode       string `hcl:"mode,optional"`

	DHCP      *DHCP             `hcl:"dhcp,block"`
	Autostart bool              `hcl:"autostart,optional"`
	Labels    map[string]string `hcl:"labels,optional"`
}

// DHCP describes the dhcp block inside a network.
type DHCP struct {
	// Range string `hcl:"range"`
	Start string `hcl:"start"`
	End   string `hcl:"end"`
}
