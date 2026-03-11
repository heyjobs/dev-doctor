package cures

import (
	"context"
	"fmt"
)

// StartDocker provides instructions for starting the Docker daemon
func StartDocker(ctx context.Context) error {
	fmt.Println("  ℹ Docker daemon is not running")
	fmt.Println()
	fmt.Println("  To start Docker:")
	fmt.Println("  1. Open Docker Desktop from Applications")
	fmt.Println("  2. Wait for Docker to start (you'll see the whale icon in the menu bar)")
	fmt.Println("  3. Once running, you can use docker commands")
	fmt.Println()
	fmt.Println("  Or run: open -a Docker")

	return nil
}
