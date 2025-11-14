package operations

import "github.com/kebairia/kvmcli/internal/resources"

func (o *Operator) Start(r resources.Resource) error {
	return r.Start()
}
