package cures

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// InstallAWSVault installs aws-vault using Homebrew
func InstallAWSVault(ctx context.Context) error {
	fmt.Println("  Checking Homebrew installation...")

	// Check if brew is installed
	if exec.CommandContext(ctx, "brew", "--version").Run() != nil {
		return fmt.Errorf("Homebrew is not installed. Install from https://brew.sh")
	}

	fmt.Println("  ✓ Homebrew is installed")
	fmt.Println()

	// Install aws-vault based on OS
	switch runtime.GOOS {
	case "darwin": // macOS
		fmt.Println("  Installing aws-vault via Homebrew...")
		cmd := exec.CommandContext(ctx, "brew", "install", "aws-vault")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to install aws-vault: %w", err)
		}

		fmt.Println()
		fmt.Println("  ✓ aws-vault installed successfully")

	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	return nil
}
