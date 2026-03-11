package cures

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// InstallBrewfile installs all packages from the infrastructure Brewfile
func InstallBrewfile(ctx context.Context) error {
	fmt.Println("  Checking Homebrew installation...")

	// Check if brew is installed
	if exec.CommandContext(ctx, "brew", "--version").Run() != nil {
		return fmt.Errorf("Homebrew is not installed. Install from https://brew.sh")
	}

	fmt.Println("  ✓ Homebrew is installed")
	fmt.Println()

	// Fetch Brewfile from infrastructure repository
	fmt.Println("  Fetching Brewfile from infrastructure repository...")
	brewfileContent, err := fetchBrewfileContent(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch Brewfile: %w", err)
	}
	fmt.Println("  ✓ Brewfile fetched successfully")
	fmt.Println()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("cannot determine home directory: %w", err)
	}

	brewfilePath := filepath.Join(homeDir, ".dev-doctor-Brewfile")
	err = os.WriteFile(brewfilePath, []byte(brewfileContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to create Brewfile: %w", err)
	}
	defer os.Remove(brewfilePath)

	fmt.Println("  Installing infrastructure tools...")
	fmt.Println("  This may take several minutes...")
	fmt.Println()

	// Run brew bundle install
	cmd := exec.CommandContext(ctx, "brew", "bundle", "install", "--file", brewfilePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install packages: %w", err)
	}

	fmt.Println()
	fmt.Println("  ✓ All infrastructure tools installed successfully")

	return nil
}

// fetchBrewfileContent fetches the full Brewfile content from the infrastructure repo
func fetchBrewfileContent(ctx context.Context) (string, error) {
	// Check if gh CLI is available
	if exec.CommandContext(ctx, "gh", "--version").Run() != nil {
		return "", fmt.Errorf("GitHub CLI (gh) is not installed")
	}

	// Fetch Brewfile content using gh CLI
	cmd := exec.CommandContext(ctx, "gh", "api", "repos/heyjobs/infrastructure/contents/Brewfile", "--jq", ".content")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to fetch Brewfile from GitHub: %w", err)
	}

	// Decode base64 content
	decodeCmd := exec.CommandContext(ctx, "base64", "-d")
	decodeCmd.Stdin = strings.NewReader(strings.TrimSpace(string(output)))
	brewfileContent, err := decodeCmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to decode Brewfile content: %w", err)
	}

	return string(brewfileContent), nil
}
