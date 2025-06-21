package network

// IDEA: the ip <=>  mac address mapping done here in Virtual Network declaration
// What I need is whenever I create a new virtual machine with a static ip, I need to update
// my virtual network declaration to add the ip <=> mac address mapping

// Struct definition
type VirtualNetworkConfig struct {
	// Conn to hold the libvirt connection
	ApiVersion string      `yaml:"apiVersion"`
	Kind       string      `yaml:"kind"`
	Metadata   NetMetadata `yaml:"metadata"`
	Spec       NetSpec     `yaml:"spec"`
}
type NetMetadata struct {
	Name      string            `yaml:"name"`
	Namespace string            `yaml:"namespace"`
	Labels    map[string]string `yaml:"labels"`
}
type NetSpec struct {
	DHCP       map[string]string `yaml:"dhcp"`
	Bridge     string            `yaml:"bridge"`
	Mode       string            `yaml:"mode"`
	Network    Network           `yaml:"network"`
	Autostart  bool              `yaml:"autostart"`
	MacAddress string            `yaml:"macAddress"`
}

type Network struct {
	Address string `yaml:"address"`
	Netmask string `yaml:"netmask"`
}

// func (net *VirtualNetwork) SetConnection(ctx context.Context, db *sql.DB, conn *libvirt.Libvirt) {
// 	net.Conn = conn
// 	net.DB = db
// 	net.Context = ctx
// }
