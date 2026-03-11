package diagnostics

import (
	"context"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/yourusername/dev-doctor/internal/types"
)

// CheckHomebrewUpdated verifies that Homebrew is up to date
func CheckHomebrewUpdated(ctx context.Context) (types.Status, string, error) {
	// First check brew is installed
	if exec.CommandContext(ctx, "brew", "--version").Run() != nil {
		return types.StatusCritical, "Homebrew is not installed", nil
	}

	cmd := exec.CommandContext(ctx, "brew", "outdated", "--verbose")
	output, err := cmd.Output()
	if err != nil {
		return types.StatusWarning, "Could not check Homebrew update status", nil
	}

	outdated := strings.TrimSpace(string(output))
	if outdated != "" {
		return types.StatusWarning, "Homebrew has outdated packages", nil
	}

	// Check if brew's formula index is stale (not updated in 7 days)
	repoCmd := exec.CommandContext(ctx, "brew", "--repository")
	repoOutput, err := repoCmd.Output()
	if err != nil {
		return types.StatusWarning, "Could not determine Homebrew repository path", nil
	}

	fetchHead := strings.TrimSpace(string(repoOutput)) + "/.git/FETCH_HEAD"
	info, err := os.Stat(fetchHead)
	if err != nil {
		return types.StatusWarning, "Homebrew formulae have never been updated", nil
	}

	if time.Since(info.ModTime()) > 7*24*time.Hour {
		return types.StatusWarning, "Homebrew formulae are not up to date", nil
	}

	return types.StatusHealthy, "Homebrew is up to date", nil
}
