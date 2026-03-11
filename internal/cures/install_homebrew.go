package cures

import (
	"context"
	"fmt"
	"os"
	"os/exec"
)

// InstallHomebrew installs Homebrew on macOS using the official install script
func InstallHomebrew(ctx context.Context) error {
	// Check if already installed
	if exec.CommandContext(ctx, "brew", "--version").Run() == nil {
		fmt.Println("  ✓ Homebrew is already installed")
		return nil
	}

	fmt.Println("  Installing Homebrew...")
	fmt.Println("  ⏱  This may take several minutes...")
	fmt.Println()

	cmd := exec.CommandContext(ctx, "/bin/bash", "-c",
		`curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh | bash`)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install Homebrew: %w", err)
	}

	fmt.Println()
	fmt.Println("  ✓ Homebrew installed successfully")
	fmt.Println()
	fmt.Println("  ⚠ You may need to add Homebrew to your PATH. Follow any instructions printed above.")

	return nil
}
