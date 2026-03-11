package diagnostics

import (
	"context"
	"os/exec"
	"strings"

	"github.com/yourusername/dev-doctor/internal/types"
)

// CheckDockerCompose verifies that docker-compose is installed
func CheckDockerCompose(ctx context.Context) (types.Status, string, error) {
	// Try docker compose (v2 plugin syntax)
	cmd := exec.CommandContext(ctx, "docker", "compose", "version")
	output, err := cmd.Output()
	if err == nil {
		version := strings.TrimSpace(string(output))
		return types.StatusHealthy, version, nil
	}

	// Try docker-compose (v1 standalone syntax)
	cmd = exec.CommandContext(ctx, "docker-compose", "--version")
	output, err = cmd.Output()
	if err == nil {
		version := strings.TrimSpace(string(output))
		return types.StatusHealthy, version, nil
	}

	return types.StatusCritical, "docker-compose is not installed", nil
}
