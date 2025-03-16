package operations

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/kebairia/kvmcli/internal/logger"
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
		logger.Log.Fatalf("error executing qemu-img: %s\n\touptut: %s", err, string(output))
	}
	logger.Log.Debugf("Overlay image created successfully \n")
	return nil
}
