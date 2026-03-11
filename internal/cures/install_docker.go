package cures

import (
	"context"
	"fmt"
	"os"
	"os/exec"
)

// InstallDocker installs Docker Desktop via Homebrew
func InstallDocker(ctx context.Context) error {
	fmt.Println("  Checking for Docker Desktop...")

	// Check if Docker Desktop app already exists
	if _, err := os.Stat("/Applications/Docker.app"); err == nil {
		fmt.Println("  ✓ Docker Desktop is already installed in Applications")
		fmt.Println()
		fmt.Println("  ℹ If docker commands are not working, try:")
		fmt.Println("  1. Launch Docker Desktop from Applications")
		fmt.Println("  2. Restart your terminal")
		return nil
	}

	fmt.Println("  Checking Homebrew installation...")

	// Check if brew is installed
	if exec.CommandContext(ctx, "brew", "--version").Run() != nil {
		fmt.Println()
		fmt.Println("  ⚠ Homebrew is not installed")
		fmt.Println()
		fmt.Println("  Please install Docker Desktop manually:")
		fmt.Println("  https://docs.docker.com/desktop/install/mac-install/")
		return nil
	}

	fmt.Println("  ✓ Homebrew is installed")
	fmt.Println()

	// Check if Docker is already installed via brew
	cmd := exec.CommandContext(ctx, "brew", "list", "--cask", "docker")
	if err := cmd.Run(); err == nil {
		fmt.Println("  ✓ Docker cask is already installed via Homebrew")
		fmt.Println("  ℹ Docker Desktop app should be in /Applications/Docker.app")
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
		fmt.Println()
		fmt.Println("  ✖ Failed to install Docker via Homebrew")
		fmt.Println()
		fmt.Println("  Please install Docker Desktop manually:")
		fmt.Println("  https://docs.docker.com/desktop/install/mac-install/")
		return nil
	}

	fmt.Println()
	fmt.Println("  ✓ Docker Desktop installed successfully")
	fmt.Println()
	fmt.Println("  ⚠ Please launch Docker Desktop from Applications to complete setup")

	return nil
}
