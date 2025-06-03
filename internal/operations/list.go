package operations

import (
	"context"
	"fmt"
	"time"

	"github.com/kebairia/kvmcli/internal/vms"
)

// -------------------------------------------
// IDEA: resources, err := operator.GetResources()
//
//		if err != nil {
//			return err
//	}
//
// w := resource.Header()
//
//	for _, resource := range resources {
//		resource.PrintRow(w)
//	}
//
// -------------------------------------------

// rec := db.VirtualMachineRecord{}
// err = rec.GetRecordByNamespace(operator.ctx, operator.db, vm.Name, namespace)
// if err != nil {
// 	return err
// }
// info, err = vms.NewVirtualMachineInfo(operator.ctx, operator.conn, rec)
// -------------------------------------------

func ListByNamespace(namespace string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	operator, err := NewOperator(ctx)
	if err != nil {
		return err
	}
	defer operator.Close()
	// resources , err := operator.GetResourcesByNamespace(namespace)
	virtualMachines, err := vms.GetVirtualMachineByNamespace(
		operator.ctx,
		operator.db,
		operator.conn,
		namespace,
	)

	info := &vms.VirtualMachineInfo{}
	w := info.Header()
	// -------------------------------------------

	for _, vm := range virtualMachines {
		vm.PrintInfo(w)
	}
	w.Flush()
	return nil
}

func ListAll() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	operator, err := NewOperator(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize operator: %w", err)
	}
	defer operator.Close()

	virtualMachines, err := vms.GetVirtualMachines(operator.ctx, operator.db, operator.conn)
	if err != nil {
		return fmt.Errorf("failed to retrieve VMs: %w", err)
	}

	info := &vms.VirtualMachineInfo{}
	w := info.Header()

	for _, virtualMachine := range virtualMachines {
		// fmt.Printf("database ==> %p\t", &operator.db)
		// fmt.Printf("connection ==> %p\n", &operator.conn)
		virtualMachine.PrintInfo(w)
	}
	w.Flush()
	return nil
}
