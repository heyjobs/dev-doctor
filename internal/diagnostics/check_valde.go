package diagnostics

import (
	"context"

	"github.com/yourusername/dev-doctor/internal/types"
)

// CheckValdeTest is a sanity-check diagnostic that always passes
func CheckValdeTest(ctx context.Context) (types.Status, string, error) {
	return types.StatusHealthy, "valde_test diagnostic is working correctly", nil
}

// CheckValdeWarning is a sanity-check diagnostic that always returns a warning
func CheckValdeWarning(ctx context.Context) (types.Status, string, error) {
	return types.StatusWarning, "valde_warning: this is what a warning looks like", nil
}

// CheckValdeCritical is a sanity-check diagnostic that always returns critical
func CheckValdeCritical(ctx context.Context) (types.Status, string, error) {
	return types.StatusCritical, "valde_critical: this is what a critical failure looks like", nil
}