package operations

import (
	"log"
	"net"

	"github.com/digitalocean/go-libvirt"
)

// CustomDialer implements socket.Dialer
type CustomDialer struct {
	network, address string
}

// CustomDialer implements socket.Dialer
func (d *CustomDialer) Dial() (net.Conn, error) {
	return net.Dial(d.network, d.address)
}

// ConnectLibvirt create and returns a libvirt connection
func ConnectLibvirt(network, address string) (*libvirt.Libvirt, error) {
	dialer := &CustomDialer{
		network: network,
		address: address,
	}
	l := libvirt.NewWithDialer(dialer)
	if err := l.Connect(); err != nil {
		log.Fatalf("failed to establish libvirt connection: %v", err)
	}

	return l, nil
}
