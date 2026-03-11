package cures

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// InstallDocker installs Docker Desktop via Homebrew
func InstallDocker(ctx context.Context) error {
	fmt.Println("  Checking for Docker Desktop...")

	// Check if Docker Desktop app already exists
	if _, err := os.Stat("/Applications/Docker.app"); err == nil {
		fmt.Println("  ✓ Docker Desktop is already installed in Applications")
		fmt.Println()
		fmt.Println("  ℹ If docker commands are not working, try:")
		fmt.Println("  1. Launch Docker Desktop from Applications")
		fmt.Println("  2. Restart your terminal")
		return nil
	}

	fmt.Println("  Checking Homebrew installation...")

	// Check if brew is installed
	if exec.CommandContext(ctx, "brew", "--version").Run() != nil {
		fmt.Println()
		fmt.Println("  ⚠ Homebrew is not installed")
		fmt.Println()
		printDockerManualInstallInstructions()
		return nil
	}

	fmt.Println("  ✓ Homebrew is installed")
	fmt.Println()

	// Check if Docker is already installed via brew
	cmd := exec.CommandContext(ctx, "brew", "list", "--cask", "docker")
	if err := cmd.Run(); err == nil {
		fmt.Println("  ✓ Docker cask is already installed via Homebrew")
		fmt.Println("  ℹ Docker Desktop app should be in /Applications/Docker.app")
		return nil
	}

	fmt.Println("  Installing Docker Desktop via Homebrew...")
	fmt.Println("  ⏱  This may take several minutes...")
	fmt.Println()

	// Install Docker Desktop as a cask
	cmd = exec.CommandContext(ctx, "brew", "install", "--cask", "docker")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println()
		fmt.Println("  ✖ Failed to install Docker via Homebrew")
		fmt.Println()
		printDockerManualInstallInstructions()
		return nil
	}

	fmt.Println()
	fmt.Println("  ✓ Docker Desktop installed successfully")
	fmt.Println()
	fmt.Println("  ⚠ Please launch Docker Desktop from Applications to complete setup")

	return nil
}

// printDockerManualInstallInstructions shows OS-specific Docker installation guidance
func printDockerManualInstallInstructions() {
	// Get macOS version
	cmd := exec.Command("sw_vers", "-productVersion")
	output, err := cmd.Output()
	macOSVersion := ""
	if err == nil {
		macOSVersion = strings.TrimSpace(string(output))
	}

	fmt.Println("  Homebrew installation failed.")
	fmt.Println()

	// Check if running macOS 13 (Ventura)
	if macOSVersion != "" && strings.HasPrefix(macOSVersion, "13.") {
		fmt.Println("  ⚠ You are running macOS 13 (Ventura)")
		fmt.Println("  Latest Docker Desktop requires macOS 14 (Sonoma) or later")
		fmt.Println()
		fmt.Println("  Attempting to download and install compatible version (4.43.0)...")
		fmt.Println()

		// Try to install Docker Desktop 4.43.0 for macOS 13
		if err := installDockerDMG(context.Background(), "4.43.0", "198134"); err != nil {
			fmt.Println()
			fmt.Println("  ✖ Automatic installation failed:", err)
			fmt.Println()
			printManualInstallFallback()
		}
	} else {
		fmt.Println("  Please install Docker Desktop manually:")
		fmt.Println("  https://docs.docker.com/desktop/install/mac-install/")
	}
}

// printManualInstallFallback shows manual installation instructions
func printManualInstallFallback() {
	fmt.Println("  Manual installation options:")
	fmt.Println()
	fmt.Println("  1. Upgrade to macOS 14+ and install latest Docker Desktop:")
	fmt.Println("     https://docs.docker.com/desktop/install/mac-install/")
	fmt.Println()
	fmt.Println("  2. Download Docker Desktop 4.43.0 manually:")
	fmt.Println("     https://desktop.docker.com/mac/main/amd64/198134/Docker.dmg (Intel)")
	fmt.Println("     https://desktop.docker.com/mac/main/arm64/198134/Docker.dmg (Apple Silicon)")
}

// installDockerDMG downloads and installs Docker Desktop from DMG
func installDockerDMG(ctx context.Context, version string, buildNumber string) error {
	// Detect CPU architecture
	cmd := exec.CommandContext(ctx, "uname", "-m")
	archOutput, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to detect CPU architecture: %w", err)
	}

	arch := strings.TrimSpace(string(archOutput))
	var downloadArch string
	if arch == "arm64" {
		downloadArch = "arm64"
		fmt.Println("  Detected: Apple Silicon (ARM64)")
	} else {
		downloadArch = "amd64"
		fmt.Println("  Detected: Intel (AMD64)")
	}

	// Construct download URL
	downloadURL := fmt.Sprintf("https://desktop.docker.com/mac/main/%s/%s/Docker.dmg", downloadArch, buildNumber)
	fmt.Println("  Download URL:", downloadURL)
	fmt.Println()

	// Create temp directory for download
	tempDir, err := os.MkdirTemp("", "docker-install-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	dmgPath := filepath.Join(tempDir, "Docker.dmg")

	// Download DMG
	fmt.Println("  Downloading Docker Desktop", version, "...")
	fmt.Println("  ⏱  This may take several minutes (file is ~500MB)...")
	fmt.Println()

	if err := downloadFile(ctx, downloadURL, dmgPath); err != nil {
		return fmt.Errorf("failed to download DMG: %w", err)
	}

	fmt.Println("  ✓ Download complete")
	fmt.Println()

	// Mount DMG
	fmt.Println("  Mounting DMG...")
	mountCmd := exec.CommandContext(ctx, "hdiutil", "attach", dmgPath, "-nobrowse", "-quiet")
	mountOutput, err := mountCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to mount DMG: %w\n%s", err, string(mountOutput))
	}

	// Parse mount point
	mountPoint := ""
	lines := strings.Split(string(mountOutput), "\n")
	for _, line := range lines {
		if strings.Contains(line, "/Volumes/Docker") {
			fields := strings.Fields(line)
			if len(fields) > 0 {
				mountPoint = fields[len(fields)-1]
				break
			}
		}
	}

	if mountPoint == "" {
		return fmt.Errorf("failed to determine mount point")
	}

	defer func() {
		exec.Command("hdiutil", "detach", mountPoint, "-quiet").Run()
	}()

	fmt.Println("  ✓ DMG mounted at", mountPoint)
	fmt.Println()

	// Copy Docker.app to /Applications
	fmt.Println("  Installing Docker.app to /Applications...")
	fmt.Println("  (This may require sudo password)")
	fmt.Println()

	dockerAppSource := filepath.Join(mountPoint, "Docker.app")
	dockerAppDest := "/Applications/Docker.app"

	// Remove existing Docker.app if present
	if _, err := os.Stat(dockerAppDest); err == nil {
		fmt.Println("  Removing existing Docker.app...")
		removeCmd := exec.CommandContext(ctx, "sudo", "rm", "-rf", dockerAppDest)
		removeCmd.Stdout = os.Stdout
		removeCmd.Stderr = os.Stderr
		if err := removeCmd.Run(); err != nil {
			return fmt.Errorf("failed to remove existing Docker.app: %w", err)
		}
	}

	// Copy new Docker.app
	copyCmd := exec.CommandContext(ctx, "sudo", "cp", "-R", dockerAppSource, dockerAppDest)
	copyCmd.Stdout = os.Stdout
	copyCmd.Stderr = os.Stderr
	if err := copyCmd.Run(); err != nil {
		return fmt.Errorf("failed to copy Docker.app: %w", err)
	}

	fmt.Println()
	fmt.Println("  ✓ Docker Desktop installed successfully!")
	fmt.Println()
	fmt.Println("  ⚠ Please launch Docker Desktop from Applications to complete setup")

	return nil
}

// downloadFile downloads a file from URL to destination path
func downloadFile(ctx context.Context, url string, dest string) error {
	// Create the file
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	// Create HTTP request with context
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	// Download the file
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Write the body to file with progress
	_, err = io.Copy(out, resp.Body)
	return err
}
