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

const dbtAnalyticsProjectPath = "Documents/github_repos/bi_analytics_dbt"

const dbtStagingHost = "main-dwh-staging.heyjobs.de"
const dbtProdHost = "dwh.heyjobs.de"

func dbtProjectDir() string {
	return filepath.Join(os.Getenv("HOME"), dbtAnalyticsProjectPath)
}

// CheckDbtAnalyticsVenvActive verifies the bi_analytics_dbt Python venv is currently activated
func CheckDbtAnalyticsVenvActive(ctx context.Context) (types.Status, string, error) {
	virtualEnv := os.Getenv("VIRTUAL_ENV")
	if virtualEnv == "" {
		return types.StatusCritical, "No Python venv is active - run: source venv/bin/activate", nil
	}

	expectedVenv := filepath.Join(dbtProjectDir(), "venv")
	if virtualEnv != expectedVenv {
		return types.StatusWarning, fmt.Sprintf("Wrong venv is active (%s) - expected bi_analytics_dbt venv", virtualEnv), nil
	}

	return types.StatusHealthy, "bi_analytics_dbt venv is active", nil
}

// CheckDbtInstalled verifies dbt is in PATH and the Redshift adapter is installed
func CheckDbtInstalled(ctx context.Context) (types.Status, string, error) {
	// Use venv dbt binary directly
	dbtBin := filepath.Join(dbtProjectDir(), "venv", "bin", "dbt")
	if _, err := os.Stat(dbtBin); err != nil {
		// Fall back to system dbt if venv doesn't exist
		dbtBin = "dbt"
	}

	cmd := exec.CommandContext(ctx, dbtBin, "--version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return types.StatusCritical, "dbt is not installed or not in PATH - run: pip install dbt-redshift", nil
	}

	if !strings.Contains(string(output), "redshift") {
		return types.StatusCritical, "dbt-redshift adapter is not installed - run: pip install dbt-redshift", nil
	}

	return types.StatusHealthy, "dbt with Redshift adapter is installed", nil
}

// CheckDbtVenv verifies the Python venv exists in the bi_analytics_dbt project
func CheckDbtVenv(ctx context.Context) (types.Status, string, error) {
	activatePath := filepath.Join(dbtProjectDir(), "venv", "bin", "activate")

	if _, err := os.Stat(activatePath); os.IsNotExist(err) {
		return types.StatusCritical, fmt.Sprintf("Python venv not found at %s - run: python -m venv venv && pip install dbt-redshift", activatePath), nil
	}

	return types.StatusHealthy, "Python venv exists and is ready to activate", nil
}

// CheckDbtSecretConfig verifies the dbt environment variables are set and point to Staging.
// Variables are expected to be loaded by running: source secret_config.env
func CheckDbtSecretConfig(ctx context.Context) (types.Status, string, error) {
	requiredVars := []string{"DBT_HOST", "DBT_REDSHIFT_USER", "DBT_REDSHIFT_PASSWORD", "DBT_PROFILES_DIR"}

	var missing []string
	for _, v := range requiredVars {
		if os.Getenv(v) == "" {
			missing = append(missing, v)
		}
	}

	if len(missing) > 0 {
		return types.StatusCritical, fmt.Sprintf("dbt env vars not set (%s) - run: source secret_config.env", strings.Join(missing, ", ")), nil
	}

	if os.Getenv("DBT_HOST") == dbtProdHost {
		return types.StatusWarning, fmt.Sprintf("dbt credentials point to Production (%s) - switch to Staging (%s)", dbtProdHost, dbtStagingHost), nil
	}

	return types.StatusHealthy, fmt.Sprintf("dbt env vars set and pointing to Staging (%s)", os.Getenv("DBT_HOST")), nil
}

// CheckDbtPackages verifies dbt packages have been installed via `dbt deps`
func CheckDbtPackages(ctx context.Context) (types.Status, string, error) {
	packagesDir := filepath.Join(dbtProjectDir(), "dbt_packages")

	entries, err := os.ReadDir(packagesDir)
	if os.IsNotExist(err) {
		return types.StatusWarning, "dbt packages not installed - run: dbt deps", nil
	}
	if err != nil {
		return types.StatusWarning, fmt.Sprintf("Cannot read dbt_packages directory: %v", err), nil
	}

	if len(entries) == 0 {
		return types.StatusWarning, "dbt_packages directory is empty - run: dbt deps", nil
	}

	return types.StatusHealthy, fmt.Sprintf("dbt packages installed (%d packages found)", len(entries)), nil
}
