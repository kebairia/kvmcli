package main

import (
	"fmt"

	"github.com/kebairia/kvmcli/internal/config"
)

func main() {
	config, err := config.LoadConfig("./servers.yaml")
	if err != nil {
		fmt.Println("Error: ", err)
	}
	for name, vm := range config.VMs {
		fmt.Println(name)
		fmt.Println(vm.Disk.Path)
		fmt.Println(vm.Network.MAC)

	}
}
