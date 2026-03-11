package diagnostics

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/yourusername/dev-doctor/internal/types"
)

// CheckWaspVersion checks if Wasp is installed
func CheckWaspVersion(ctx context.Context) (types.Status, string, error) {
	// Check if wasp command exists
	cmd := exec.CommandContext(ctx, "wasp", "version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return types.StatusCritical, "Wasp is not installed or not in PATH", nil
	}

	// Parse version from output
	versionStr := strings.TrimSpace(string(output))
	if versionStr != "" {
		return types.StatusHealthy, fmt.Sprintf("Wasp %s is installed", versionStr), nil
	}

	return types.StatusHealthy, "Wasp is installed", nil
}
