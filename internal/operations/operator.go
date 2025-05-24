package operations

import (
	"context"

	"github.com/digitalocean/go-libvirt"
	"github.com/kebairia/kvmcli/internal"
	"github.com/kebairia/kvmcli/internal/resources"
)

// Operator manages the lifecycle of resource lifecycle.
// and other !operations!
// It holds the execution context, configuration, libvirt connection, and logger.
type Operator struct {
	ctx context.Context
	// configuration file
	// config config.Config
	// libvirt connection
	conn *libvirt.Libvirt
	// logger
	// log logger.Logger
}

func NewOperator() (*Operator, error) {
	ctx := context.Background()
	conn, err := internal.InitConnection()
	if err != nil {
		return nil, err
	}
	// setting logger
	// log := logger.Global()

	return &Operator{
		ctx:  ctx,
		conn: conn,
		// log:  log,
	}, nil
}

func (o *Operator) Close() {
	if o.conn != nil {
		o.conn.Disconnect()
	}
}

func (o *Operator) Create(resource resources.Resource) error {
	// If the resource requires a libvirt connection, inject it.
	if cs, ok := resource.(resources.ClientSetter); ok {
		cs.SetConnection(o.conn)
	}
	// create resource
	err := resource.Create()
	if err != nil {
		return err
	}
	return nil
}

func (o *Operator) Delete(resource resources.Resource) error {
	// If the resource requires a libvirt connection, inject it.
	if cs, ok := resource.(resources.ClientSetter); ok {
		cs.SetConnection(o.conn)
	}
	// delete resource
	err := resource.Delete()
	if err != nil {
		return err
	}
	return nil
}
