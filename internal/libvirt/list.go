package libvirt

import (
	"fmt"
	"log"
	"net/url"

	"github.com/digitalocean/go-libvirt"
)

func ListAll() {
	uri, err := url.Parse(string(libvirt.QEMUSystem))
	if err != nil {
		fmt.Println(err)
	}
	l, err := libvirt.ConnectToURI(uri)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	// First we print the version
	v, err := l.ConnectGetLibVersion()
	if err != nil {
		log.Fatalf("failed to retrieve libvirt version: %v", err)
	}
	fmt.Printf("Version: %d\n", v)

	// Now we retreive all the instance (running and down)
	flags := libvirt.ConnectListDomainsActive | libvirt.ConnectListDomainsInactive
	domains, _, err := l.ConnectListAllDomains(1, flags)
	if err != nil {
		log.Fatalf("failed to retrieve domains: %v", err)
	}

	fmt.Println("ID\tName\t\tUUID")
	fmt.Printf("--------------------------------------------------------\n")
	for _, d := range domains {
		fmt.Printf("%d\t%s\t\t%x\n", d.ID, d.Name, d.UUID)
	}

	if err = l.Disconnect(); err != nil {
		log.Fatalf("failed to disconnect: %v", err)
	}
}
