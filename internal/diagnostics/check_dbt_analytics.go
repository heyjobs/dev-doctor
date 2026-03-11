package diagnostics

import (
	"bufio"
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
	cmd := exec.CommandContext(ctx, "dbt", "--version")
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

// CheckDbtSecretConfig verifies secret_config.env exists, credentials are configured,
// and the active credentials point to Staging (not Production).
func CheckDbtSecretConfig(ctx context.Context) (types.Status, string, error) {
	configPath := filepath.Join(dbtProjectDir(), "secret_config.env")

	f, err := os.Open(configPath)
	if os.IsNotExist(err) {
		return types.StatusCritical, "secret_config.env not found - copy from secret_config.env.template and fill in credentials", nil
	}
	if err != nil {
		return types.StatusCritical, fmt.Sprintf("Cannot read secret_config.env: %v", err), nil
	}
	defer f.Close()

	var unconfigured []string
	var activeHost string

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		varName := strings.TrimPrefix(parts[0], "export ")
		value := parts[1]

		if strings.Contains(value, "<TO BE SET>") {
			unconfigured = append(unconfigured, varName)
		}

		if varName == "DBT_HOST" {
			activeHost = value
		}
	}

	if len(unconfigured) > 0 {
		return types.StatusCritical, fmt.Sprintf("Credentials not configured in secret_config.env: %s", strings.Join(unconfigured, ", ")), nil
	}

	if activeHost == dbtProdHost {
		return types.StatusWarning, fmt.Sprintf("secret_config.env is pointing to Production (%s) - switch to Staging credentials (%s)", dbtProdHost, dbtStagingHost), nil
	}

	return types.StatusHealthy, fmt.Sprintf("secret_config.env configured with Staging credentials (%s)", activeHost), nil
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