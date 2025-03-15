package operations

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

const artifactsPath = "/home/zakaria/dox/homelab/artifacts/rocky"

func CreateOverlay(baseImage, destImage string) error {
	baseImage = filepath.Join(artifactsPath, baseImage)
	cmd := exec.Command(
		"qemu-img",
		"create",
		"-o",
		fmt.Sprintf("backing_file=%s,backing_fmt=qcow2", baseImage),
		"-f",
		"qcow2",
		destImage,
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error executing qemu-img: %v\nouptut: %s", err, string(output))
	}
	fmt.Printf("Overlay image created successfully \n")
	return nil
}
