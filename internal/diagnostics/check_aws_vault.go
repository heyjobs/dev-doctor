package diagnostics

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/yourusername/dev-doctor/internal/types"
)

// CheckAWSVault checks if aws-vault is installed
func CheckAWSVault(ctx context.Context) (types.Status, string, error) {
	// Check if aws-vault command exists
	cmd := exec.CommandContext(ctx, "aws-vault", "--version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Command not found or failed to execute
		return types.StatusCritical, "aws-vault is not installed or not in PATH", nil
	}

	// aws-vault is installed, get version info
	versionStr := strings.TrimSpace(string(output))
	return types.StatusHealthy, fmt.Sprintf("aws-vault is installed (%s)", versionStr), nil
}
