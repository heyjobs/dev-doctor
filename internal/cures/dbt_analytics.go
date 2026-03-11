package cures

import "context"

// InstallDbtRedshift placeholder: install dbt-redshift via pip
func InstallDbtRedshift(ctx context.Context) error {
	// pip install dbt-redshift
	return nil
}

// SetupDbtVenv placeholder: create venv and install dbt-redshift
func SetupDbtVenv(ctx context.Context) error {
	// python -m venv venv && venv/bin/pip install dbt-redshift
	return nil
}

// SetupDbtSecretConfig placeholder: guide user to configure secret_config.env
func SetupDbtSecretConfig(ctx context.Context) error {
	// cp secret_config.env.template secret_config.env
	return nil
}

// RunDbtDeps placeholder: run dbt deps to install packages
func RunDbtDeps(ctx context.Context) error {
	// source venv/bin/activate && source secret_config.env && dbt deps
	return nil
}