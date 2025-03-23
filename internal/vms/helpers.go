package vms

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/kebairia/kvmcli/internal/logger"
)

const artifactsPath = "/home/zakaria/dox/homelab/artifacts/rocky"

// ------------------------------------
// Function for creating the overlay
// ------------------------------------
func CreateOverlay(namespace string, baseImage, destImage string) error {
	baseImage = filepath.Join(artifactsPath, baseImage)
	// destImage = filepath.Join(namespace, "-", destImage)
	cmd := exec.Command(
		"qemu-img",
		"create",
		"-o",
		fmt.Sprintf("backing_file=%s,backing_fmt=qcow2", baseImage),
		"-f",
		"qcow2",
		destImage,
	)
	if _, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%w", err)
	}
	logger.Log.Debugf("Overlay image created successfully \n")
	return nil
}
