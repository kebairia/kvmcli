package resources

import "github.com/digitalocean/go-libvirt"

// Resource defines operations for managing a resource.
type Resource interface {
	Create() error
	Delete() error
}

// ClientSetter is implemented by types that require a libvirt connection.
type ClientSetter interface {
	SetConnection(*libvirt.Libvirt)
}
