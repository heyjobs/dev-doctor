package diagnostics

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/yourusername/dev-doctor/internal/types"
)

// CheckBrewfile checks if all packages from the infrastructure Brewfile are installed
func CheckBrewfile(ctx context.Context) (types.Status, string, error) {
	// Check if brew is installed
	if exec.CommandContext(ctx, "brew", "--version").Run() != nil {
		return types.StatusCritical, "Homebrew is not installed", nil
	}

	// Fetch required packages from infrastructure Brewfile
	requiredPackages, err := fetchBrewfilePackages(ctx)
	if err != nil {
		return types.StatusCritical, fmt.Sprintf("Failed to fetch Brewfile: %s", err.Error()), nil
	}

	if len(requiredPackages) == 0 {
		return types.StatusWarning, "No packages found in Brewfile", nil
	}

	var missing []string

	// Check each package
	for _, pkg := range requiredPackages {
		cmd := exec.CommandContext(ctx, "brew", "list", pkg)
		if err := cmd.Run(); err != nil {
			missing = append(missing, pkg)
		}
	}

	if len(missing) == 0 {
		return types.StatusHealthy, "All infrastructure tools are installed", nil
	}

	if len(missing) == len(requiredPackages) {
		return types.StatusCritical, "No infrastructure tools installed", nil
	}

	return types.StatusWarning, fmt.Sprintf("%d tools missing: %s", len(missing), strings.Join(missing, ", ")), nil
}

// fetchBrewfilePackages fetches and parses the Brewfile from the infrastructure repo
func fetchBrewfilePackages(ctx context.Context) ([]string, error) {
	// Check if gh CLI is available
	if exec.CommandContext(ctx, "gh", "--version").Run() != nil {
		return nil, fmt.Errorf("GitHub CLI (gh) is not installed")
	}

	// Fetch Brewfile content using gh CLI (this may take a moment)
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
