package operations

import (
	"context"
	"fmt"
	"time"

	"github.com/kebairia/kvmcli/internal/network"
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
		fmt.Printf("database ==> %p\t", &operator.db)
		fmt.Printf("connection ==> %p\n", &operator.conn)
		virtualMachine.PrintInfo(w)
	}
	w.Flush()
	return nil
}

func ListAllNetworks() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	operator, err := NewOperator(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize operator: %w", err)
	}
	defer operator.Close()

	virtualNetworks, err := network.GetVirtualNetworks(operator.ctx, operator.db, operator.conn)
	if err != nil {
		return fmt.Errorf("failed to retrieve networks: %w", err)
	}

	info := &network.VirtualNetworkInfo{}
	w := info.Header()

	for _, virtualNetwork := range virtualNetworks {
		fmt.Printf("database ==> %p\t", &operator.db)
		fmt.Printf("connection ==> %p\n", &operator.conn)
		virtualNetwork.PrintInfo(w)
	}
	w.Flush()
	return nil
}

// func List(namespace, table string) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
// 	defer cancel()
// 	operator, err := NewOperator(ctx)
// 	if err != nil {
// 		return fmt.Errorf("failed to initialize operator: %w", err)
// 	}
// 	defer operator.Close()
// 	if namespace == "" {
// 		return nil
// 	}
//
// 	resourcesInfo, err := GetResourcesInfo(operator.ctx, operator.db, table)
// 	if err != nil {
// 		return fmt.Errorf("failed to retreive resources information: %w", err)
// 	}
// 	w := *&tabwriter.Writer{}
// 	for _, info := range resourcesInfo {
// 		info.PrintInfo(w)
// 	}
// 	return nil
// }
