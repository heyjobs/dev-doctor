package diagnostics

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/yourusername/dev-doctor/internal/types"
)

// CheckAWSSSOSetup checks if AWS SSO is properly configured
func CheckAWSSSOSetup(ctx context.Context) (types.Status, string, error) {
	var issues []string

	// Check 1: AWS CLI installed
	if exec.CommandContext(ctx, "aws", "--version").Run() != nil {
		return types.StatusCritical, "AWS CLI is not installed", nil
	}

	// Check 2: JQ installed
	if exec.CommandContext(ctx, "jq", "--version").Run() != nil {
		issues = append(issues, "jq is not installed")
	}

	// Check 3: AWS SSO cache exists
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return types.StatusCritical, "Cannot determine home directory", nil
	}

	ssoCacheDir := filepath.Join(homeDir, ".aws", "sso", "cache")
	cacheFiles, err := filepath.Glob(filepath.Join(ssoCacheDir, "*.json"))
	if err != nil || len(cacheFiles) == 0 {
		issues = append(issues, "AWS SSO not configured (no cache files found)")
	}

	// Check 4: AWS config has SSO profiles
	configFile := filepath.Join(homeDir, ".aws", "config")
	configData, err := os.ReadFile(configFile)
	if err != nil {
		issues = append(issues, "AWS config file not found")
	} else {
		// Check if config has sso_start_url (indicates SSO profiles)
		if !strings.Contains(string(configData), "sso_start_url") {
			issues = append(issues, "No SSO profiles found in AWS config")
		}
	}

	// Return status based on issues found
	if len(issues) == 0 {
		return types.StatusHealthy, "AWS SSO is properly configured", nil
	}

	if len(issues) == 1 {
		return types.StatusWarning, issues[0], nil
	}

	return types.StatusWarning, fmt.Sprintf("Multiple issues: %s", strings.Join(issues, "; ")), nil
}
