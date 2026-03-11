package diagnostics

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/yourusername/dev-doctor/internal/types"
)

// CheckClaudeCLI checks if Claude CLI is installed and accessible
func CheckClaudeCLI(ctx context.Context) (types.Status, string, error) {
	cmd := exec.CommandContext(ctx, "claude", "--version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return types.StatusCritical, "Claude CLI is not installed or not in PATH", nil
	}

	versionStr := strings.TrimSpace(string(output))
	if versionStr != "" {
		return types.StatusHealthy, fmt.Sprintf("Claude CLI is installed (%s)", versionStr), nil
	}

	return types.StatusHealthy, "Claude CLI is installed", nil
}
