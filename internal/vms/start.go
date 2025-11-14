package vms

import "fmt"

func (vm *VirtualMachine) Start() error {
	vm.domain.Start(vm.ctx, vm.Config.Metadata.Name)
	fmt.Printf("vm/%s started\n", vm.Config.Metadata.Name)
	return nil
}
