package main

import (
	"fmt"

	"github.com/kebairia/kvmcli/internal/utils"
)

func main() {
	domain := utils.NewDomain(
		"code-test",
		5120,
		2,
		"/home/zakaria/dox/homelab/artifacts/rocky-linux-9.qcow2",
		"52:54:00:e9:a9:29",
	)
	xmlOutput, _ := domain.GenerateXML()
	fmt.Println(string(xmlOutput))
}
