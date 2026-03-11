package cures

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// SetupDbtVenv creates the Python venv and installs dbt-redshift into it.
// Assumes the current working directory is the bi_analytics_dbt repo root.
// Used as cure for both check_dbt_analytics_venv_active and check_dbt_venv.
func SetupDbtVenv(ctx context.Context) error {
	activatePath := filepath.Join("venv", "bin", "activate")

	// Create venv if it doesn't already exist
	if _, err := os.Stat(activatePath); os.IsNotExist(err) {
		fmt.Println("  Creating Python venv...")
		cmd := exec.CommandContext(ctx, "python3", "-m", "venv", "venv")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to create venv: %w", err)
		}
		fmt.Println("  ✓ venv created")
	} else {
		fmt.Println("  ✓ venv already exists")
	}

	// Install dbt-redshift using the venv's pip directly
	fmt.Println("  Installing dbt-redshift into venv (this may take a minute)...")
	cmd := exec.CommandContext(ctx, "venv/bin/pip", "install", "dbt-redshift")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install dbt-redshift: %w", err)
	}

	fmt.Println("  ✓ dbt-redshift installed")
	fmt.Println()
	fmt.Println("  ℹ Activate the venv with: source venv/bin/activate")

	return nil
}

// InstallDbtRedshift installs dbt-redshift into the active venv.
// Falls back to the local venv pip if no venv is currently active.
func InstallDbtRedshift(ctx context.Context) error {
	pipBin := "pip"
	if os.Getenv("VIRTUAL_ENV") == "" {
		if _, err := os.Stat("venv/bin/pip"); err != nil {
			return fmt.Errorf("no venv is active and venv/bin/pip not found — run the venv cure first")
		}
		fmt.Println("  No venv active, using local venv pip...")
		pipBin = "venv/bin/pip"
	}

	fmt.Println("  Installing dbt-redshift...")
	cmd := exec.CommandContext(ctx, pipBin, "install", "dbt-redshift")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install dbt-redshift: %w", err)
	}

	fmt.Println("  ✓ dbt-redshift installed successfully")
	return nil
}

// SetupDbtSecretConfig activates and verifies the dbt environment with AWS credentials.
// It checks that venv, aws-vault profile, and secret_config.env all exist and work together.
func SetupDbtSecretConfig(ctx context.Context) error {
	fmt.Println("  Setting up dbt environment with AWS credentials...")
	fmt.Println()

	// 1. Verify venv exists
	venvPath := filepath.Join("venv", "bin", "activate")
	if _, err := os.Stat(venvPath); os.IsNotExist(err) {
		return fmt.Errorf("venv not found - run the venv cure first")
	}
	fmt.Println("  ✓ venv exists")

	// 2. Verify aws-vault profile exists
	fmt.Println("  Checking aws-vault profile...")
	checkCmd := exec.CommandContext(ctx, "aws-vault", "list")
	output, err := checkCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("aws-vault not working: %w", err)
	}
	if !strings.Contains(string(output), "data_platform_sso_staging") {
		return fmt.Errorf("aws-vault profile 'data_platform_sso_staging' not found - configure it first")
	}
	fmt.Println("  ✓ aws-vault profile 'data_platform_sso_staging' exists")

	// 3. Verify secret_config.env exists, or create it from template
	secretConfigPath := "secret_config.env"
	templatePath := "secret_config.env.template"

	if _, err := os.Stat(secretConfigPath); os.IsNotExist(err) {
		// Check if template exists
		if _, err := os.Stat(templatePath); os.IsNotExist(err) {
			return fmt.Errorf("secret_config.env not found and no template available")
		}

		// Copy template to secret_config.env
		fmt.Println("  Creating secret_config.env from template...")
		cmd := exec.CommandContext(ctx, "cp", templatePath, secretConfigPath)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to copy template: %w", err)
		}
		fmt.Println("  ✓ secret_config.env created from template")
		fmt.Println()
		fmt.Println("  ⚠️  IMPORTANT: Edit secret_config.env and configure your credentials!")
		fmt.Println()
	} else {
		fmt.Println("  ✓ secret_config.env exists")
	}

	// 4. Test the full environment activation
	fmt.Println("  Testing environment activation with aws-vault...")
	testCmd := exec.CommandContext(ctx, "aws-vault", "exec", "data_platform_sso_staging", "--",
		"bash", "-c", "source secret_config.env && echo $DBT_HOST")
	testOutput, err := testCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to activate environment: %w", err)
	}

	dbHost := strings.TrimSpace(string(testOutput))
	if dbHost == "" {
		return fmt.Errorf("DBT_HOST not set in secret_config.env")
	}
	fmt.Printf("  ✓ Environment activated successfully (DBT_HOST=%s)\n", dbHost)
	fmt.Println()
	fmt.Println("  ℹ  To use dbt interactively, run:")
	fmt.Println("      source venv/bin/activate")
	fmt.Println("      aws-vault exec data_platform_sso_staging -- bash")
	fmt.Println("      source secret_config.env")

	return nil
}

// RunDbtDeps runs `dbt deps` in the current directory (bi_analytics_dbt repo root),
// using aws-vault to activate AWS credentials and sourcing secret_config.env.
func RunDbtDeps(ctx context.Context) error {
	// Use venv dbt binary directly so it works regardless of whether venv is active
	dbtBin := "venv/bin/dbt"
	if _, err := os.Stat(dbtBin); err != nil {
		dbtBin = "dbt"
	}

	// Ensure secret_config.env exists (copy from template if needed)
	secretConfigPath := "secret_config.env"
	templatePath := "secret_config.env.template"

	if _, err := os.Stat(secretConfigPath); os.IsNotExist(err) {
		if _, err := os.Stat(templatePath); os.IsNotExist(err) {
			return fmt.Errorf("secret_config.env not found and no template available")
		}

		fmt.Println("  Creating secret_config.env from template...")
		cpCmd := exec.CommandContext(ctx, "cp", templatePath, secretConfigPath)
		if err := cpCmd.Run(); err != nil {
			return fmt.Errorf("failed to copy template: %w", err)
		}
		fmt.Println("  ✓ secret_config.env created from template")
		fmt.Println("  ⚠️  Remember to configure your credentials in secret_config.env")
		fmt.Println()
	}

	fmt.Println("  Running dbt deps with AWS credentials (data_platform_sso_staging)...")
	cmd := exec.CommandContext(ctx, "aws-vault", "exec", "data_platform_sso_staging", "--",
		"bash", "-c", fmt.Sprintf("source secret_config.env && %s deps", dbtBin))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("dbt deps failed: %w", err)
	}

	fmt.Println("  ✓ dbt packages installed successfully")
	return nil
}
