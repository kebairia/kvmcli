package network

// IDEA: the ip <=>  mac address mapping done here in Virtual Network declaration
// What I need is whenever I create a new virtual machine with a static ip, I need to update
// my virtual network declaration to add the ip <=> mac address mapping

type Network struct {
	Name      string `hcl:"name,label"`
	Namespace string `hcl:"namespace"`
	// CIDR      string `hcl:"cidr"`
	NetAddress string `hcl:"netaddress"`
	NetMask    string `hcl:"netmask"`
	Bridge     string `hcl:"bridge"`
	Mode       string `hcl:"mode"`

	// DHCP      *DHCP             `hcl:"dhcp,block"`
	// DHCP      map[string]string `hcl:"dhcp,block"`
	// DHCP      *DHCP             `hcl:"dhcp,block"`
	Autostart bool              `hcl:"autostart"`
	Labels    map[string]string `hcl:"labels,attr"`
}

// DHCP describes the dhcp block inside a network.
// type DHCP struct {
// 	// Range string `hcl:"range"`
// 	Start string `hcl:"start"`
// 	End   string `hcl:"end"`
// }

// Struct definition
// type VirtualNetworkConfig struct {
// 	// Conn to hold the libvirt connection
// 	ApiVersion string      `yaml:"apiVersion"`
// 	Kind       string      `yaml:"kind"`
// 	Metadata   NetMetadata `yaml:"metadata"`
// 	Spec       NetSpec     `yaml:"spec"`
// }
// type NetMetadata struct {
// 	Name      string            `yaml:"name"`
// 	Namespace string            `yaml:"namespace"`
// 	Labels    map[string]string `yaml:"labels"`
// }
// type NetSpec struct {
// 	DHCP       map[string]string `yaml:"dhcp"`
// 	Bridge     string            `yaml:"bridge"`
// 	Mode       string            `yaml:"mode"`
// 	Network    Network           `yaml:"network"`
// 	Autostart  bool              `yaml:"autostart"`
// 	MacAddress string            `yaml:"macAddress"`
// }
//
// type Network struct {
// 	Address string `yaml:"address"`
// 	Netmask string `yaml:"netmask"`
// }

// func (net *VirtualNetwork) SetConnection(ctx context.Context, db *sql.DB, conn *libvirt.Libvirt) {
// 	net.Conn = conn
// 	net.DB = db
// 	net.Context = ctx
// }

// type Network struct {
// 	Name      string `hcl:"name,label"`
// 	Namespace string `hcl:"namespace,label"`
// 	CIDR      string `hcl:"cidr"`
// 	Bridge    string `hcl:"bridge"`
// 	Mode      string `hcl:"mode"`
//
// 	DHCP      *DHCP             `hcl:"dhcp,block"`
// 	Autostart bool              `hcl:"autostart"`
// 	Labels    map[string]string `hcl:"labels,attr"`
// }
//
// // DHCP describes the dhcp block inside a network.
// type DHCP struct {
// 	Range string `hcl:"range"`
// }
