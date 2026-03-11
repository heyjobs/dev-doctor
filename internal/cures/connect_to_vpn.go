package cures

import (
	"context"
	"fmt"
)

// ConnectToVPN provides instructions for connecting to VPN
func ConnectToVPN(ctx context.Context) error {
	fmt.Println("  To connect to VPN:")
	fmt.Println("  1. Open your VPN client application")
	fmt.Println("  2. Connect to the company VPN")
	fmt.Println("  3. Wait for connection to establish")
	fmt.Println()
	fmt.Println("  ⚠ Cannot auto-connect - please connect manually")
	return nil
}
