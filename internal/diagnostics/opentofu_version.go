package diagnostics

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/yourusername/dev-doctor/internal/types"
)

// CheckOpenTofuVersion checks if OpenTofu (tofu) is installed
func CheckOpenTofuVersion(ctx context.Context) (types.Status, string, error) {
	// Check if tofu command exists
	cmd := exec.CommandContext(ctx, "tofu", "--version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return types.StatusCritical, "OpenTofu is not installed or not in PATH", nil
	}

	// Parse version from output (e.g., "OpenTofu v1.11.4")
	versionStr := strings.TrimSpace(string(output))
	lines := strings.Split(versionStr, "\n")
	if len(lines) > 0 {
		firstLine := lines[0]
		// Extract version like "v1.11.4" from "OpenTofu v1.11.4"
		if strings.Contains(firstLine, "OpenTofu") {
			return types.StatusHealthy, fmt.Sprintf("%s is installed", firstLine), nil
		}
	}

	return types.StatusHealthy, fmt.Sprintf("OpenTofu is installed (%s)", versionStr), nil
}
