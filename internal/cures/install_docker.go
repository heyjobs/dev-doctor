package cures

import (
	"context"
	"fmt"
	"os"
	"os/exec"
)

// InstallDocker installs Docker Desktop via Homebrew
func InstallDocker(ctx context.Context) error {
	fmt.Println("  Checking Homebrew installation...")

	// Check if brew is installed
	if exec.CommandContext(ctx, "brew", "--version").Run() != nil {
		return fmt.Errorf("Homebrew is not installed. Install from https://brew.sh")
	}

	fmt.Println("  ✓ Homebrew is installed")
	fmt.Println()

	// Check if Docker is already installed
	cmd := exec.CommandContext(ctx, "brew", "list", "--cask", "docker")
	if err := cmd.Run(); err == nil {
		fmt.Println("  ✓ Docker is already installed via Homebrew")
		return nil
	}

	fmt.Println("  Installing Docker Desktop via Homebrew...")
	fmt.Println("  ⏱  This may take several minutes...")
	fmt.Println()

	// Install Docker Desktop as a cask
	cmd = exec.CommandContext(ctx, "brew", "install", "--cask", "docker")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install Docker: %w", err)
	}

	fmt.Println()
	fmt.Println("  ✓ Docker Desktop installed successfully")
	fmt.Println()
	fmt.Println("  ⚠ Please launch Docker Desktop from Applications to complete setup")

	return nil
}
