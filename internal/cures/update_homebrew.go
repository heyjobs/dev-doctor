package cures

import (
	"context"
	"fmt"
	"os"
	"os/exec"
)

// UpdateHomebrew updates Homebrew to the latest version
func UpdateHomebrew(ctx context.Context) error {
	fmt.Println("  Updating Homebrew...")

	cmd := exec.CommandContext(ctx, "brew", "update")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to update Homebrew: %w", err)
	}

	fmt.Println()
	fmt.Println("  ✓ Homebrew updated successfully")

	return nil
}
