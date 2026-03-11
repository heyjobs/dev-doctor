package cures

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

// InstallWasp installs wasp-go from GitHub
func InstallWasp(ctx context.Context) error {
	fmt.Println("  ⚠️  IMPORTANT: Wasp requires prerequisites to work properly!")
	fmt.Println()
	fmt.Println("  Before installing Wasp, you MUST complete ALL prerequisite steps:")
	fmt.Println("  https://heyjobs.atlassian.net/wiki/spaces/dnp/pages/3158278303/How-to+-+Create+Redshift+Auth+With+Temp+Credentials+via+IDE+and+Wasp#1.-Ask-to-be-added-to-SSO-groups")
	fmt.Println()
	fmt.Println("  Prerequisites include:")
	fmt.Println("  1. Being added to SSO groups")
	fmt.Println("  2. Configuring AWS SSO")
	fmt.Println("  3. Setting up Redshift access")
	fmt.Println()
	fmt.Println("  ⚠️  Wasp WILL NOT WORK without these prerequisites!")
	fmt.Println()
	fmt.Println("  Proceeding with Wasp installation...")
	fmt.Println()

	// Check if Go is installed
	fmt.Println("  Checking for Go installation...")
	if _, err := exec.LookPath("go"); err != nil {
		return fmt.Errorf("Go is not installed. Please install Go first: brew install go")
	}
	fmt.Println("  ✓ Go is installed")
	fmt.Println()

	// Check if git is configured to use SSH for GitHub
	fmt.Println("  Checking git configuration for private repository access...")
	checkSSHCmd := exec.CommandContext(ctx, "git", "config", "--global", "--get", "url.git@github.com:.insteadOf")
	output, err := checkSSHCmd.Output()

	needsSSHConfig := err != nil || string(output) != "https://github.com/\n"

	if needsSSHConfig {
		fmt.Println("  Configuring git to use SSH for GitHub private repositories...")
		configCmd := exec.CommandContext(ctx, "git", "config", "--global", "url.git@github.com:.insteadOf", "https://github.com/")
		if err := configCmd.Run(); err != nil {
			fmt.Println("  ⚠️  Warning: Could not configure git SSH rewrite")
			fmt.Println()
		} else {
			fmt.Println("  ✓ Git configured to use SSH for GitHub")
		}
	} else {
		fmt.Println("  ✓ Git already configured to use SSH")
	}
	fmt.Println()

	// Install wasp-go from GitHub (private repository)
	// Note: We need to clone and build because the module path doesn't match the repo path
	fmt.Println("  Installing wasp from github.com/heyjobs/wasp-go...")

	// Create temporary directory for cloning
	tmpDir, err := os.MkdirTemp("", "wasp-go-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	// Clone the repository
	fmt.Println("  Cloning repository...")
	cloneCmd := exec.CommandContext(ctx, "git", "clone", "git@github.com:heyjobs/wasp-go.git", tmpDir)
	cloneCmd.Stdout = os.Stdout
	cloneCmd.Stderr = os.Stderr
	if err := cloneCmd.Run(); err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	// Build the binary
	fmt.Println("  Building wasp...")
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", "wasp", ".")
	buildCmd.Dir = tmpDir
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("failed to build wasp: %w", err)
	}

	// Determine installation directory
	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		goPath = filepath.Join(os.Getenv("HOME"), "go")
	}
	goBin := filepath.Join(goPath, "bin")

	// Create bin directory if it doesn't exist
	if err := os.MkdirAll(goBin, 0755); err != nil {
		return fmt.Errorf("failed to create bin directory: %w", err)
	}

	// Install the binary
	fmt.Println("  Installing binary...")
	srcBinary := filepath.Join(tmpDir, "wasp")
	dstBinary := filepath.Join(goBin, "wasp")

	// Copy the binary
	if err := copyFile(srcBinary, dstBinary); err != nil {
		return fmt.Errorf("failed to install binary: %w", err)
	}

	// Make it executable
	if err := os.Chmod(dstBinary, 0755); err != nil {
		return fmt.Errorf("failed to make binary executable: %w", err)
	}

	fmt.Println("  ✓ Wasp installed successfully")
	fmt.Println()

	// Check if $GOPATH/bin is in PATH (goPath and goBin already defined above)
	fmt.Println("  Checking PATH configuration...")
	pathEnv := os.Getenv("PATH")
	if !contains(pathEnv, goBin) {
		fmt.Println()
		fmt.Println("  ⚠️  WARNING: Go binaries directory is not in your PATH!")
		fmt.Println()
		fmt.Println("  Add this to your ~/.zshrc or ~/.bashrc:")
		fmt.Printf("    export PATH=\"$PATH:%s\"\n", goBin)
		fmt.Println()
		fmt.Println("  Then run: source ~/.zshrc")
		fmt.Println()
	} else {
		fmt.Println("  ✓ Go binaries directory is in PATH")
		fmt.Println()
	}

	// Verify wasp is now accessible
	fmt.Println("  Verifying wasp installation...")
	waspPath := filepath.Join(goBin, "wasp")
	if _, err := os.Stat(waspPath); err == nil {
		fmt.Printf("  ✓ Wasp installed at: %s\n", waspPath)
		fmt.Println()

		// Test wasp command
		testCmd := exec.CommandContext(ctx, "wasp", "version")
		if err := testCmd.Run(); err == nil {
			fmt.Println("  ✓ Wasp command is working!")
		} else {
			fmt.Println()
			fmt.Println("  ⚠️  Wasp is installed but not accessible via 'wasp' command")
			fmt.Println("  You may need to restart your terminal or run: source ~/.zshrc")
		}
	} else {
		return fmt.Errorf("wasp installation completed but binary not found at %s", waspPath)
	}

	fmt.Println()
	fmt.Println("  ℹ️  Remember: Wasp requires the prerequisites mentioned above to function!")
	fmt.Println()

	return nil
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}
