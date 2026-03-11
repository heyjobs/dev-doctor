package diagnostics

import (
	"context"
	"os/exec"

	"github.com/yourusername/dev-doctor/internal/types"
)

// CheckDockerRunning verifies that Docker daemon is running
func CheckDockerRunning(ctx context.Context) (types.Status, string, error) {
	// Try to run a simple docker command that requires the daemon
	cmd := exec.CommandContext(ctx, "docker", "info")
	err := cmd.Run()
	if err != nil {
		return types.StatusCritical, "Docker daemon is not running", nil
	}

	return types.StatusHealthy, "Docker daemon is running", nil
}
