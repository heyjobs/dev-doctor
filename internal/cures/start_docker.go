package cures

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

// StartDocker starts the Docker Desktop application
func StartDocker(ctx context.Context) error {
	fmt.Println("  Starting Docker Desktop...")
	fmt.Println()

	// Try to open Docker Desktop application
	cmd := exec.CommandContext(ctx, "open", "-a", "Docker")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start Docker Desktop: %w. Make sure Docker Desktop is installed", err)
	}

	fmt.Println("  ✓ Docker Desktop launched")
	fmt.Println()
	fmt.Println("  ⏱  Waiting for Docker daemon to start (this may take 10-30 seconds)...")
	fmt.Println()

	// Wait for Docker daemon to be ready
	maxAttempts := 30
	for i := 0; i < maxAttempts; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			cmd := exec.CommandContext(ctx, "docker", "info")
			if err := cmd.Run(); err == nil {
				fmt.Println("  ✓ Docker daemon is running")
				return nil
			}
			time.Sleep(2 * time.Second)
		}
	}

	return fmt.Errorf("Docker daemon did not start within expected time. Please check Docker Desktop manually")
}
