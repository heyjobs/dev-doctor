package diagnostics

import (
	"context"
	"math/rand"

	"github.com/yourusername/dev-doctor/internal/types"
)

// CheckGitConfiguration simulates checking Git user configuration
func CheckGitConfiguration(ctx context.Context) (types.Status, string, error) {
	// Mock: usually healthy
	if rand.Float32() < 0.85 {
		return types.StatusHealthy, "Git user.name and user.email are configured", nil
	}
	return types.StatusWarning, "Git user configuration is missing", nil
}
