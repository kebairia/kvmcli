package internal

import (
	"log"
	"net/url"

	"github.com/digitalocean/go-libvirt"
)

// ConnectLibvirt create and returns a libvirt connection
func InitConnection() (*libvirt.Libvirt, error) {
	uri, _ := url.Parse(string(libvirt.QEMUSystem))
	l, err := libvirt.ConnectToURI(uri)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	return l, nil
}
