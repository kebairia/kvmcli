package libvirt

import (
	"log"
	"net"
	"time"

	"github.com/digitalocean/go-libvirt"
)

func CreateVM() {
	// Connect to Libvirt
	conn, err := net.DialTimeout("unix", "/var/run/libvirt/libvirt-sock", 2*time.Second)
	if err != nil {
		log.Fatalf("failed to connect to libvirt: %v", err)
	}
	defer conn.Close()

	l := libvirt.New(conn)
	if err := l.Connect(); err != nil {
		log.Fatalf("failed to establish libvirt connection: %v", err)
	}
	defer l.Disconnect()

	// Define the XML for the new VM
}
