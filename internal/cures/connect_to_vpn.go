package cures

import (
	"context"
	"fmt"
)

// ConnectToVPN provides instructions for connecting to VPN
func ConnectToVPN(ctx context.Context) error {
	fmt.Println("  ⚠ Check here to set up VPN:")
	fmt.Println("  https://heyjobs.atlassian.net/wiki/spaces/dnp/pages/2722070581/How+to+set+up+HeyJobs+VPN")
	fmt.Println()
	fmt.Println("  VPN is needed for:")
	fmt.Println("  - Connecting to Redshift via IDE (using Wasp)")
	fmt.Println("  - Connecting to Airflow instances")
	return nil
}
