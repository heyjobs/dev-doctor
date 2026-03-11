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

// SetupDbtSecretConfig prints instructions to load dbt environment variables.
// Sourcing a file into the parent shell cannot be done from a subprocess,
// so the cure tells the user exactly what to run.
func SetupDbtSecretConfig(ctx context.Context) error {
	fmt.Println("  dbt environment variables are not set in this shell session.")
	fmt.Println("  They are loaded by sourcing secret_config.env, which is done by start_env.sh.")
	fmt.Println()
	fmt.Println("  Run the following from the bi_analytics_dbt repo root:")
	fmt.Println()
	fmt.Println("      source secret_config.env")
	fmt.Println()
	fmt.Println("  Or run the full environment setup:")
	fmt.Println()
	fmt.Println("      source venv/bin/activate && source secret_config.env")

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
