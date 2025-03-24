package vms

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/kebairia/kvmcli/internal/logger"
)

const artifactsPath = "/home/zakaria/dox/homelab/artifacts/rocky"

// CreateOverlay creates a qcow2 overlay image based on a backing file.
func CreateOverlay(baseImage, destImage string) error {
	// Construct the full path for the base image.
	baseImagePath := filepath.Join(artifactsPath, baseImage)

	// Prepare the qemu-img command.
	cmdArgs := []string{
		"create",
		"-o", fmt.Sprintf("backing_file=%s,backing_fmt=qcow2", baseImagePath),
		"-f", "qcow2",
		destImage,
	}

	// Create a context with a timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "qemu-img", cmdArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log.Debugf("%s", output)
		return fmt.Errorf("failed execute qemu-img command: %w", err)
	}

	logger.Log.Debug("Overlay image created successfully")
	return nil
}
