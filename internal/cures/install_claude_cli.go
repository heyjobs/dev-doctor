package cures

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// InstallClaudeCLI installs Claude CLI using Homebrew
func InstallClaudeCLI(ctx context.Context) error {
	fmt.Println("  Checking Homebrew installation...")

	if exec.CommandContext(ctx, "brew", "--version").Run() != nil {
		return fmt.Errorf("Homebrew is not installed. Install from https://brew.sh")
	}

	fmt.Println("  ✓ Homebrew is installed")
	fmt.Println()

	switch runtime.GOOS {
	case "darwin":
		fmt.Println("  Installing Claude CLI via Homebrew...")
		cmd := exec.CommandContext(ctx, "brew", "install", "--cask", "claude-code")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to install Claude CLI: %w", err)
		}

		fmt.Println()
		fmt.Println("  ✓ Claude CLI installed successfully")

	default:
		return fmt.Errorf("unsupported operating system: %s. Install manually: npm install -g @anthropic-ai/claude-code", runtime.GOOS)
	}

	return nil
}
