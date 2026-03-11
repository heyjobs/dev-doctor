package cures

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// InstallBrewfile installs missing packages from the infrastructure Brewfile
func InstallBrewfile(ctx context.Context) error {
	fmt.Println("  Checking Homebrew installation...")

	// Check if brew is installed
	if exec.CommandContext(ctx, "brew", "--version").Run() != nil {
		return fmt.Errorf("Homebrew is not installed. Install from https://brew.sh")
	}

	fmt.Println("  ✓ Homebrew is installed")
	fmt.Println()

	// Fetch required packages from infrastructure Brewfile
	fmt.Println("  Fetching package list from infrastructure repository...")
	requiredPackages, err := fetchBrewfilePackages(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch Brewfile: %w", err)
	}
	fmt.Println("  ✓ Package list fetched successfully")
	fmt.Println()

	// Check which packages are missing
	fmt.Println("  Checking installed packages...")
	var missing []string
	for _, pkg := range requiredPackages {
		cmd := exec.CommandContext(ctx, "brew", "list", pkg)
		if err := cmd.Run(); err != nil {
			missing = append(missing, pkg)
		}
	}

	if len(missing) == 0 {
		fmt.Println("  ✓ All infrastructure tools are already installed")
		return nil
	}

	fmt.Printf("  Found %d missing package(s): %s\n", len(missing), strings.Join(missing, ", "))
	fmt.Println()
	fmt.Println("  Installing missing packages...")
	fmt.Println("  ⏱  This may take several minutes (brew installations can be slow)...")
	fmt.Println()

	// Install each missing package individually
	for i, pkg := range missing {
		fmt.Printf("  [%d/%d] Installing %s...\n", i+1, len(missing), pkg)

		cmd := exec.CommandContext(ctx, "brew", "install", pkg)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to install %s: %w", pkg, err)
		}
	}

	fmt.Println()
	fmt.Println("  ✓ All missing infrastructure tools installed successfully")

	return nil
}

// fetchBrewfilePackages fetches and parses the Brewfile from the infrastructure repo
func fetchBrewfilePackages(ctx context.Context) ([]string, error) {
	// Check if gh CLI is available
	if exec.CommandContext(ctx, "gh", "--version").Run() != nil {
		return nil, fmt.Errorf("GitHub CLI (gh) is not installed")
	}

	// Fetch Brewfile content using gh CLI
	cmd := exec.CommandContext(ctx, "gh", "api", "repos/heyjobs/infrastructure/contents/Brewfile", "--jq", ".content")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Brewfile from GitHub: %w", err)
	}

	// Decode base64 content
	decodeCmd := exec.CommandContext(ctx, "base64", "-d")
	decodeCmd.Stdin = strings.NewReader(strings.TrimSpace(string(output)))
	brewfileContent, err := decodeCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to decode Brewfile content: %w", err)
	}

	// Parse Brewfile to extract package names
	var packages []string
	lines := strings.Split(string(brewfileContent), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Match lines like: brew "package-name"
		if strings.HasPrefix(line, "brew \"") && strings.HasSuffix(line, "\"") {
			pkg := strings.TrimPrefix(line, "brew \"")
			pkg = strings.TrimSuffix(pkg, "\"")
			packages = append(packages, pkg)
		}
	}

	return packages, nil
}
