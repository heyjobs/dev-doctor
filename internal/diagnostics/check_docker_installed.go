package diagnostics

import (
	"context"
	"os/exec"

	"github.com/yourusername/dev-doctor/internal/types"
)

// CheckDockerInstalled verifies that Docker is installed
func CheckDockerInstalled(ctx context.Context) (types.Status, string, error) {
	// Check if docker command exists
	cmd := exec.CommandContext(ctx, "docker", "--version")
	output, err := cmd.Output()
	if err != nil {
		return types.StatusCritical, "Docker is not installed", nil
	}

	return types.StatusHealthy, string(output), nil
}
