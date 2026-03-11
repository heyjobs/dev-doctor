package diagnostics

import (
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/yourusername/dev-doctor/internal/types"
)

// CheckPythonVersion checks if Python 3.10.x is installed
func CheckPythonVersion(ctx context.Context) (types.Status, string, error) {
	// Try python command
	cmd := exec.CommandContext(ctx, "python", "--version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return types.StatusCritical, "Python is not installed or not in PATH", nil
	}

	// Parse version from output (e.g., "Python 3.10.0")
	versionStr := strings.TrimSpace(string(output))
	versionRegex := regexp.MustCompile(`Python (\d+)\.(\d+)\.(\d+)`)
	matches := versionRegex.FindStringSubmatch(versionStr)

	if len(matches) != 4 {
		return types.StatusCritical, fmt.Sprintf("Unable to parse Python version from: %s", versionStr), nil
	}

	major := matches[1]
	minor := matches[2]
	patch := matches[3]

	// Check if version is 3.10.x
	if major == "3" && minor == "10" {
		return types.StatusHealthy, fmt.Sprintf("Python %s.%s.%s meets requirements (3.10.x)", major, minor, patch), nil
	}

	// Version exists but is not 3.10.x
	if major == "3" {
		if minor < "10" {
			return types.StatusWarning, fmt.Sprintf("Python %s.%s.%s is outdated (requires 3.10.x)", major, minor, patch), nil
		}
		return types.StatusWarning, fmt.Sprintf("Python %s.%s.%s is newer than required 3.10.x", major, minor, patch), nil
	}

	return types.StatusCritical, fmt.Sprintf("Python %s.%s.%s does not meet requirements (requires 3.10.x)", major, minor, patch), nil
}
