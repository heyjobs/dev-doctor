package cures

import "context"

// ValdeTestCure is a no-op cure for the valde_test diagnostic
func ValdeTestCure(ctx context.Context) error {
	return nil
}