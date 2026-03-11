package diagnostics

import (
	"context"
	"os/exec"
	"strings"

	"github.com/yourusername/dev-doctor/internal/types"
)

// CheckHomebrewInstalled verifies that Homebrew is installed on macOS
func CheckHomebrewInstalled(ctx context.Context) (types.Status, string, error) {
	cmd := exec.CommandContext(ctx, "brew", "--version")
	output, err := cmd.Output()
	if err != nil {
		return types.StatusCritical, "Homebrew is not installed", nil
	}

	version := strings.TrimSpace(string(output))
	return types.StatusHealthy, version, nil
}
