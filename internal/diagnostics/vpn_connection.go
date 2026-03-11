package diagnostics

import (
	"context"
	"os/exec"
	"strings"

	"github.com/yourusername/dev-doctor/internal/types"
)

// CheckVPNConnection checks if connected to VPN
func CheckVPNConnection(ctx context.Context) (types.Status, string, error) {
	// Check for active VPN connections using scutil
	cmd := exec.CommandContext(ctx, "scutil", "--nc", "list")
	output, err := cmd.Output()
	if err != nil {
		return types.StatusWarning, "Unable to check VPN status", nil
	}

	outputStr := string(output)
	lines := strings.Split(outputStr, "\n")

	// Look for connected VPN connections
	for _, line := range lines {
		if strings.Contains(line, "Connected") {
			// Extract connection name if possible
			return types.StatusHealthy, "Connected to VPN", nil
		}
	}

	// Also check for VPN interfaces (utun, ppp, etc) with actual IPv4 addresses
	ifconfigCmd := exec.CommandContext(ctx, "ifconfig")
	ifconfigOutput, err := ifconfigCmd.Output()
	if err == nil {
		ifconfigStr := string(ifconfigOutput)
		// Parse ifconfig output to check if any utun interface has an inet address
		currentInterface := ""
		for _, line := range strings.Split(ifconfigStr, "\n") {
			// Track which interface we're looking at
			if strings.HasPrefix(line, "utun") {
				currentInterface = "utun"
			} else if !strings.HasPrefix(line, "\t") && !strings.HasPrefix(line, " ") {
				// New interface section started
				currentInterface = ""
			}

			// If we're in a utun interface section and see an inet (IPv4) address
			if currentInterface == "utun" && strings.Contains(line, "\tinet ") {
				return types.StatusHealthy, "VPN interface detected", nil
			}
		}
	}

	return types.StatusWarning, "Not connected to VPN", nil
}
