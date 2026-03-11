package cures

import (
	"context"
	"fmt"
	"os"
	"os/exec"
)

// InstallDockerCompose installs docker-compose via Homebrew
func InstallDockerCompose(ctx context.Context) error {
	fmt.Println("  Checking for docker-compose...")
	fmt.Println()

	// Check if docker compose v2 plugin is available (comes with Docker Desktop)
	cmd := exec.CommandContext(ctx, "docker", "compose", "version")
	if err := cmd.Run(); err == nil {
		fmt.Println("  ✓ docker-compose v2 plugin is already available")
		fmt.Println("  (Installed with Docker Desktop)")
		return nil
	}

	// Check if standalone docker-compose v1 is installed
	cmd = exec.CommandContext(ctx, "docker-compose", "--version")
	if err := cmd.Run(); err == nil {
		fmt.Println("  ✓ docker-compose v1 is already installed")
		return nil
	}

	fmt.Println("  Checking Homebrew installation...")

	// Check if brew is installed
	if exec.CommandContext(ctx, "brew", "--version").Run() != nil {
		return fmt.Errorf("Homebrew is not installed. Install from https://brew.sh")
	}

	fmt.Println("  ✓ Homebrew is installed")
	fmt.Println()

	fmt.Println("  Installing docker-compose via Homebrew...")
	fmt.Println()

	// Install docker-compose
	cmd = exec.CommandContext(ctx, "brew", "install", "docker-compose")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install docker-compose: %w", err)
	}

	fmt.Println()
	fmt.Println("  ✓ docker-compose installed successfully")

	return nil
}
