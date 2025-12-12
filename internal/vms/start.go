package vms

import (
	"fmt"
)

func (vm *VirtualMachine) Start() error {
	vm.domain.Start(vm.ctx, vm.Spec.Name)
	fmt.Printf("vm/%s started\n", vm.Spec.Name)
	// cfg.NetworksByName, cfg.VMsByName, cfg.ClustersByName are now available
	return nil
}
