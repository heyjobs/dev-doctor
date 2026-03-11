package cures

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// UpdatePython installs Python 3.10.x using pyenv
func UpdatePython(ctx context.Context) error {
	// Step 1: Check if pyenv is installed
	if !isPyenvInstalled(ctx) {
		fmt.Println("  Installing pyenv...")
		if err := installPyenv(ctx); err != nil {
			return fmt.Errorf("failed to install pyenv: %w", err)
		}
		fmt.Println("  ✓ pyenv installed successfully")
	} else {
		fmt.Println("  ✓ pyenv is already installed")
	}

	// Step 2: Find latest Python 3.10.x version
	fmt.Println("  Finding latest Python 3.10.x version...")
	latestVersion, err := getLatestPython310Version(ctx)
	if err != nil {
		return fmt.Errorf("failed to find Python 3.10.x version: %w", err)
	}
	fmt.Printf("  Found: %s\n", latestVersion)

	// Step 3: Install Python 3.10.x
	fmt.Printf("  Installing Python %s (this may take a few minutes)...\n", latestVersion)
	if err := installPythonVersion(ctx, latestVersion); err != nil {
		return fmt.Errorf("failed to install Python %s: %w", latestVersion, err)
	}
	fmt.Printf("  ✓ Python %s installed successfully\n", latestVersion)

	// Step 4: Set as global default
	fmt.Println("  Setting as global default...")
	if err := setPyenvGlobal(ctx, latestVersion); err != nil {
		return fmt.Errorf("failed to set Python version: %w", err)
	}
	fmt.Println("  ✓ Python version set as global default")

	// Step 5: Check for PYENV_VERSION override
	if envVersion := os.Getenv("PYENV_VERSION"); envVersion != "" && envVersion != latestVersion {
		fmt.Println()
		fmt.Println("  ⚠ WARNING: PYENV_VERSION environment variable is set to", envVersion)
		fmt.Println("  This will override the global setting in your current shell.")
		fmt.Println()
		fmt.Println("  To use Python", latestVersion, "immediately, run:")
		fmt.Println("    unset PYENV_VERSION")
		fmt.Println()
		fmt.Println("  Or start a new shell session.")
	}

	return nil
}

// isPyenvInstalled checks if pyenv is available in PATH
func isPyenvInstalled(ctx context.Context) bool {
	cmd := exec.CommandContext(ctx, "pyenv", "--version")
	return cmd.Run() == nil
}

// installPyenv installs pyenv based on the operating system
func installPyenv(ctx context.Context) error {
	switch runtime.GOOS {
	case "darwin": // macOS
		// Try homebrew first
		if exec.CommandContext(ctx, "brew", "--version").Run() == nil {
			cmd := exec.CommandContext(ctx, "brew", "install", "pyenv")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			return cmd.Run()
		}
		// Fall back to curl installation
		return installPyenvWithCurl(ctx)

	case "linux":
		// Use curl installation script
		return installPyenvWithCurl(ctx)

	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

// installPyenvWithCurl installs pyenv using the official installer
func installPyenvWithCurl(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "bash", "-c",
		"curl https://pyenv.run | bash")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// getLatestPython310Version finds the latest Python 3.10.x version available
func getLatestPython310Version(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "pyenv", "install", "--list")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	// Parse output to find latest 3.10.x
	lines := strings.Split(string(output), "\n")
	var latest string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Match versions like "3.10.0", "3.10.1", etc. (pure numeric only)
		if strings.HasPrefix(line, "3.10.") && !strings.Contains(line, "-") {
			// Ensure it's only digits after "3.10." (no letter suffixes like "3.10.1t")
			version := strings.TrimPrefix(line, "3.10.")
			isNumeric := true
			for _, c := range version {
				if c < '0' || c > '9' {
					isNumeric = false
					break
				}
			}
			if isNumeric {
				latest = line
			}
		}
	}

	if latest == "" {
		return "", fmt.Errorf("no Python 3.10.x version found")
	}

	return latest, nil
}

// installPythonVersion installs a specific Python version using pyenv
func installPythonVersion(ctx context.Context, version string) error {
	cmd := exec.CommandContext(ctx, "pyenv", "install", "-s", version)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// setPyenvGlobal sets the global Python version
func setPyenvGlobal(ctx context.Context, version string) error {
	cmd := exec.CommandContext(ctx, "pyenv", "global", version)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
