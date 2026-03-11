package cures

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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

// SetupDbtSecretConfig copies secret_config.env.template to secret_config.env if missing.
func SetupDbtSecretConfig(ctx context.Context) error {
	if _, err := os.Stat("secret_config.env"); err == nil {
		fmt.Println("  secret_config.env already exists")
		fmt.Println("  ℹ Edit it to fill in DBT_REDSHIFT_USER and DBT_REDSHIFT_PASSWORD")
		return nil
	}

	fmt.Println("  Copying secret_config.env.template → secret_config.env...")
	data, err := os.ReadFile("secret_config.env.template")
	if err != nil {
		return fmt.Errorf("template not found at secret_config.env.template: %w", err)
	}
	if err := os.WriteFile("secret_config.env", data, 0600); err != nil {
		return fmt.Errorf("failed to write secret_config.env: %w", err)
	}

	fmt.Println("  ✓ secret_config.env created")
	fmt.Println()
	fmt.Println("  ℹ Open it and set:")
	fmt.Println("      DBT_REDSHIFT_USER=<your user>")
	fmt.Println("      DBT_REDSHIFT_PASSWORD=<your password>")

	return nil
}

// RunDbtDeps runs `dbt deps` in the current directory (bi_analytics_dbt repo root),
// sourcing secret_config.env so DBT_PROFILES_DIR and credentials are available.
func RunDbtDeps(ctx context.Context) error {
	// Use venv dbt binary directly so it works regardless of whether venv is active
	dbtBin := "venv/bin/dbt"
	if _, err := os.Stat(dbtBin); err != nil {
		dbtBin = "dbt"
	}

	fmt.Println("  Running dbt deps...")
	cmd := exec.CommandContext(ctx, "bash", "-c", fmt.Sprintf("source secret_config.env && %s deps", dbtBin))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("dbt deps failed: %w", err)
	}

	fmt.Println("  ✓ dbt packages installed successfully")
	return nil
}
