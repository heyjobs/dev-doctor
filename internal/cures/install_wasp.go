package cures

import (
	"context"
	"fmt"
)

// InstallWasp provides instructions for installing Wasp
func InstallWasp(ctx context.Context) error {
	fmt.Println("  To install Wasp, please follow the instructions here:")
	fmt.Println()
	fmt.Println("  https://heyjobs.atlassian.net/wiki/spaces/dnp/pages/3158278303/How-to+-+Create+Redshift+Auth+With+Temp+Credentials+via+IDE+and+Wasp")
	fmt.Println()
	return nil
}
