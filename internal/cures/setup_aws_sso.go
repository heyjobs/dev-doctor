package cures

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// SetupAWSSSO configures AWS SSO profiles
func SetupAWSSSO(ctx context.Context) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("cannot determine home directory: %w", err)
	}

	// Step 1: Check and install AWS CLI if missing
	if exec.CommandContext(ctx, "aws", "--version").Run() != nil {
		fmt.Println("  Installing AWS CLI...")
		if err := installAWSCLI(ctx); err != nil {
			return err
		}
		fmt.Println("  ✓ AWS CLI installed")
	} else {
		fmt.Println("  ✓ AWS CLI is already installed")
	}

	// Step 2: Check and install JQ if missing
	if exec.CommandContext(ctx, "jq", "--version").Run() != nil {
		fmt.Println("  Installing jq...")
		if err := installJQ(ctx); err != nil {
			return err
		}
		fmt.Println("  ✓ jq installed")
	} else {
		fmt.Println("  ✓ jq is already installed")
	}

	// Step 3: Check if SSO cache exists, if not provide instructions
	ssoCacheDir := filepath.Join(homeDir, ".aws", "sso", "cache")
	cacheFiles, _ := filepath.Glob(filepath.Join(ssoCacheDir, "*.json"))

	// Filter out botocore cache files
	var validCacheFiles []string
	for _, f := range cacheFiles {
		if !strings.Contains(f, "botocore") {
			validCacheFiles = append(validCacheFiles, f)
		}
	}

	if len(validCacheFiles) == 0 {
		fmt.Println()
		fmt.Println("  ⚠ AWS SSO session not found.")
		fmt.Println()
		fmt.Println("  Please run the following command:")
		fmt.Println("    aws configure sso")
		fmt.Println()
		fmt.Println("  When prompted, enter these values:")
		fmt.Println()
		fmt.Println("    SSO session name (Recommended): heyjobs")
		fmt.Println("    SSO start URL [None]: https://d-99671e8831.awsapps.com/start")
		fmt.Println("    SSO Region [None]: eu-central-1")
		fmt.Println()
		fmt.Println("  A browser window will open - click 'Allow' to authorize.")
		fmt.Println("  After authorizing, you can press Ctrl+C to exit the prompt.")
		fmt.Println()
		fmt.Println("  Then re-run: dev-doctor --mode treatment")
		return nil
	}

	fmt.Println("  ✓ AWS SSO session found")
	fmt.Println()

	// Step 4: Check if we have any SSO profiles in config for login
	configFile := filepath.Join(homeDir, ".aws", "config")
	configData, _ := os.ReadFile(configFile)
	hasSSOProfile := strings.Contains(string(configData), "sso_start_url")

	if !hasSSOProfile {
		// No SSO profile exists, create a bootstrap profile automatically
		fmt.Println("  Creating bootstrap SSO profile...")

		bootstrapProfile := `
[profile heyjobs-bootstrap]
sso_start_url = https://d-99671e8831.awsapps.com/start
sso_region = eu-central-1
sso_account_id = 000000000000
sso_role_name = bootstrap
region = eu-central-1
output = json
`
		f, err := os.OpenFile(configFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			return fmt.Errorf("failed to create bootstrap profile: %w", err)
		}
		if _, err := f.WriteString(bootstrapProfile); err != nil {
			f.Close()
			return fmt.Errorf("failed to write bootstrap profile: %w", err)
		}
		f.Close()

		fmt.Println("  ✓ Bootstrap profile created")

		// Re-read config to get the profile
		configData, _ = os.ReadFile(configFile)
	}

	// Step 5: Refresh SSO session with existing profile
	fmt.Println("  Refreshing AWS SSO session...")

	// Find first SSO profile
	lines := strings.Split(string(configData), "\n")
	var ssoProfile string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "[profile ") {
			ssoProfile = strings.TrimPrefix(line, "[profile ")
			ssoProfile = strings.TrimSuffix(ssoProfile, "]")
			break
		}
	}

	if ssoProfile == "" {
		return fmt.Errorf("could not find SSO profile in config")
	}

	fmt.Printf("  Using profile: %s\n", ssoProfile)
	fmt.Println("  A browser window will open for authentication.")
	fmt.Println()

	loginCmd := exec.CommandContext(ctx, "aws", "sso", "login", "--profile", ssoProfile)
	loginCmd.Stdout = os.Stdout
	loginCmd.Stderr = os.Stderr

	if err := loginCmd.Run(); err != nil {
		return fmt.Errorf("failed to refresh SSO session: %w", err)
	}

	fmt.Println("  ✓ SSO session refreshed")
	fmt.Println()

	// Step 6: Create and run the setup script
	fmt.Println("  Scanning AWS accounts and creating profiles...")

	profilesFile := filepath.Join(homeDir, ".aws", "config.devdoctor")
	scriptPath, err := createSSOSetupScript(homeDir, profilesFile)
	if err != nil {
		return fmt.Errorf("failed to create setup script: %w", err)
	}
	defer os.Remove(scriptPath) // Clean up after

	// Run the script non-interactively (exports to separate file)
	cmd := exec.CommandContext(ctx, "bash", scriptPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("setup script failed: %w", err)
	}

	// Check if profiles were generated
	profileData, err := os.ReadFile(profilesFile)
	if err != nil || len(profileData) == 0 {
		fmt.Println()
		fmt.Println("  ⚠ No new profiles were generated")
		return nil
	}

	// Show what was generated and ask to merge
	fmt.Println()
	fmt.Println("  ✓ Profiles generated successfully!")
	fmt.Println()
	fmt.Println("  Generated profiles saved to: ~/.aws/config.devdoctor")
	fmt.Println()
	fmt.Println("  Preview (first 20 lines):")
	fmt.Println("  " + strings.Repeat("─", 60))

	lines = strings.Split(string(profileData), "\n")
	previewLines := lines
	if len(lines) > 20 {
		previewLines = lines[:20]
	}
	for _, line := range previewLines {
		fmt.Println("  " + line)
	}
	if len(lines) > 20 {
		fmt.Printf("  ... (%d more lines)\n", len(lines)-20)
	}
	fmt.Println("  " + strings.Repeat("─", 60))

	fmt.Println()
	fmt.Print("  Merge these profiles into ~/.aws/config? [y/n]: ")

	var response string
	fmt.Scanln(&response)

	if strings.ToLower(strings.TrimSpace(response)) == "y" {
		// Append to main config
		configF, err := os.OpenFile(configFile, os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			return fmt.Errorf("failed to open config file: %w", err)
		}
		defer configF.Close()

		if _, err := configF.Write(profileData); err != nil {
			return fmt.Errorf("failed to merge profiles: %w", err)
		}

		fmt.Println()
		fmt.Println("  ✓ Profiles merged into ~/.aws/config")
		fmt.Println("  ✓ Backup kept at ~/.aws/config.devdoctor")
	} else {
		fmt.Println()
		fmt.Println("  Profiles saved to ~/.aws/config.devdoctor")
		fmt.Println("  You can manually merge them later if needed.")
	}

	fmt.Println()
	fmt.Println("  ✓ AWS SSO profiles configured successfully")
	return nil
}

// installAWSCLI installs AWS CLI based on the operating system
func installAWSCLI(ctx context.Context) error {
	switch runtime.GOOS {
	case "darwin": // macOS
		if exec.CommandContext(ctx, "brew", "--version").Run() != nil {
			return fmt.Errorf("Homebrew is required but not installed")
		}
		cmd := exec.CommandContext(ctx, "brew", "install", "awscli")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()

	case "linux":
		return fmt.Errorf("please install AWS CLI manually: https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html")

	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

// installJQ installs jq based on the operating system
func installJQ(ctx context.Context) error {
	switch runtime.GOOS {
	case "darwin": // macOS
		if exec.CommandContext(ctx, "brew", "--version").Run() != nil {
			return fmt.Errorf("Homebrew is required but not installed")
		}
		cmd := exec.CommandContext(ctx, "brew", "install", "jq")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()

	case "linux":
		return fmt.Errorf("please install jq manually: sudo apt-get install jq or sudo yum install jq")

	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

// createSSOSetupScript creates a bash script for scanning AWS accounts and creating profiles
func createSSOSetupScript(homeDir string, outputFile string) (string, error) {
	scriptContent := fmt.Sprintf(`#!/bin/bash -e
OUTPUT_FILE="%s"
rm -rf "$OUTPUT_FILE"`, outputFile) + `

# Hardcoded SSO configuration values
start_url="https://d-99671e8831.awsapps.com/start"
region="eu-central-1"

# Extract access token from cache
at_filename=$(ls ~/.aws/sso/cache/*.json | grep -v botocore | head -n 1)
at=$(cat $at_filename | jq -r '.accessToken')

# Iterate account list
available_accounts=$(aws sso list-accounts --region "$region" --access-token "$at")
n_accounts=$(echo $available_accounts | jq '.accountList | length')
echo "Accounts found: $n_accounts"
account_list=$(echo $available_accounts | jq -r '.accountList | .[] | .accountId')

for account_id in $account_list; do
    echo "account: $account_id"
    account_data=$( echo $available_accounts | jq -r ".accountList | .[] | select( .accountId == \"$account_id\" )" )
    account_name=$(echo $account_data | jq -r '.accountName // .accountId' | xargs | tr -d "[:space:]")
    account_roles=$(aws sso list-account-roles --region "$region" --access-token "$at" --account-id $account_id)
    role_names=$(echo $account_roles | jq -r '.roleList | .[] | .roleName')

    for role_name in $role_names; do
        echo "  role: $role_name"
        config_profile_name=$(echo ${role_name%???})
        hit=$(cat ~/.aws/config | grep $config_profile_name || echo "")

        if [ -z "$hit" ] ; then
            echo "    profile: $config_profile_name - adding..."
            cat << EOF >> "$OUTPUT_FILE"
[profile $config_profile_name]
sso_start_url = $start_url
sso_region = $region
sso_account_id = $account_id
sso_role_name = $role_name
sts_regional_endpoints = regional
region = $region

EOF
        else
            echo "    profile: $config_profile_name - already exists, skipping..."
        fi
    done
done

echo ""
if [ -f "$OUTPUT_FILE" ] && [ -s "$OUTPUT_FILE" ]; then
    profile_count=$(grep -c "^\[profile " "$OUTPUT_FILE" || echo "0")
    echo "✓ Generated $profile_count new profiles"
else
    echo "No new profiles generated (all already exist)"
fi
echo "Done!"
`

	scriptPath := filepath.Join(homeDir, ".dev-doctor-aws-sso-setup.sh")
	err := os.WriteFile(scriptPath, []byte(scriptContent), 0755)
	if err != nil {
		return "", err
	}

	return scriptPath, nil
}
