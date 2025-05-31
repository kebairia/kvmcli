package operations

import "github.com/kebairia/kvmcli/internal/resources"

func (o *Operator) Insert(r resources.Record) error {
	return r.Insert(o.ctx, o.db)
}
