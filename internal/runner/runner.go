package runner

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/yourusername/dev-doctor/internal/cures"
	"github.com/yourusername/dev-doctor/internal/diagnostics"
	"github.com/yourusername/dev-doctor/internal/types"
)

// Runner orchestrates the execution of diagnostic tests and treatments
type Runner struct {
	diagnosticRegistry *diagnostics.Registry
	cureRegistry       *cures.Registry
	config             *types.DiagnosticConfig
	diagnosticTimeout  time.Duration
	cureTimeout        time.Duration
}

// NewRunner creates a new diagnostic runner
func NewRunner(config *types.DiagnosticConfig) *Runner {
	return &Runner{
		diagnosticRegistry: diagnostics.DefaultRegistry(),
		cureRegistry:       cures.DefaultRegistry(),
		config:             config,
		diagnosticTimeout:  30 * time.Second, // Default timeout per diagnostic
		cureTimeout:        30 * time.Minute, // Default timeout per cure (longer for brew installations which can be slow)
	}
}

// ResultCallback is called after each diagnostic test completes
type ResultCallback func(result types.DiagnosticResult)

// RunDiagnostics executes all diagnostic tests and returns results
func (r *Runner) RunDiagnostics(ctx context.Context) (*types.Summary, error) {
	return r.RunDiagnosticsWithCallback(ctx, nil)
}

// RunDiagnosticsWithCallback executes all diagnostic tests and calls the callback after each one
func (r *Runner) RunDiagnosticsWithCallback(ctx context.Context, callback ResultCallback) (*types.Summary, error) {
	results := make([]types.DiagnosticResult, 0, len(r.config.Tests))

	for _, test := range r.config.Tests {
		result, err := r.runSingleDiagnostic(ctx, test)
		if err != nil {
			// Log error but continue with other tests
			result = types.DiagnosticResult{
				TestID:      test.Test,
				Description: test.Description,
				Status:      types.StatusCritical,
				Summary:     fmt.Sprintf("Test execution failed: %v", err),
				Symptom:     test.Symptom,
				CureID:      test.Cure,
				FixAvailable: false,
				Severity:    test.Severity,
			}
		}
		results = append(results, result)

		// Call callback immediately after each test completes
		if callback != nil {
			callback(result)
		}
	}

	return r.buildSummary(results), nil
}

// runSingleDiagnostic executes a single diagnostic test
func (r *Runner) runSingleDiagnostic(ctx context.Context, test types.DiagnosticTest) (types.DiagnosticResult, error) {
	// Get the diagnostic function
	diagFunc, err := r.diagnosticRegistry.Get(test.Diagnostic)
	if err != nil {
		return types.DiagnosticResult{}, fmt.Errorf("diagnostic not found: %s", test.Diagnostic)
	}

	// Create context with timeout
	testCtx, cancel := context.WithTimeout(ctx, r.diagnosticTimeout)
	defer cancel()

	// Execute the diagnostic
	status, summary, err := diagFunc(testCtx)
	if err != nil {
		return types.DiagnosticResult{}, err
	}

	// Override status based on severity from config when check fails
	if status != types.StatusHealthy {
		status = severityToStatus(test.Severity)
	}

	// Check if cure is available
	fixAvailable := false
	if test.Cure != "" {
		_, err := r.cureRegistry.Get(test.Cure)
		fixAvailable = (err == nil)
	}

	return types.DiagnosticResult{
		TestID:       test.Test,
		Description:  test.Description,
		Status:       status,
		Summary:      summary,
		Symptom:      test.Symptom,
		CureID:       test.Cure,
		FixAvailable: fixAvailable,
		Severity:     test.Severity,
	}, nil
}

// severityToStatus maps severity from config to status
func severityToStatus(severity types.Severity) types.Status {
	switch severity {
	case types.SeverityCritical:
		return types.StatusCritical
	case types.SeverityWarning:
		return types.StatusWarning
	case types.SeverityInfo:
		return types.StatusInfo
	default:
		return types.StatusWarning
	}
}

// TreatmentCallback is called before each treatment is applied
type TreatmentCallback func(result types.DiagnosticResult)

// ApplyTreatments applies cures to failed diagnostics
func (r *Runner) ApplyTreatments(ctx context.Context, results []types.DiagnosticResult) error {
	return r.ApplyTreatmentsWithCallback(ctx, results, nil)
}

// ApplyTreatmentsWithCallback applies cures and calls callback before each treatment
func (r *Runner) ApplyTreatmentsWithCallback(ctx context.Context, results []types.DiagnosticResult, callback TreatmentCallback) error {
	for _, result := range results {
		// Only apply treatments to non-healthy tests with available fixes
		if result.Status == types.StatusHealthy || !result.FixAvailable {
			continue
		}

		// Get the cure function
		cureFunc, err := r.cureRegistry.Get(result.CureID)
		if err != nil {
			return fmt.Errorf("cure not found for %s: %w", result.TestID, err)
		}

		// Call callback before applying treatment
		if callback != nil {
			callback(result)
		}

		// Create context with timeout (longer for cures like Python installation)
		cureCtx, cancel := context.WithTimeout(ctx, r.cureTimeout)
		defer cancel()

		// Apply the cure
		if err := cureFunc(cureCtx); err != nil {
			return fmt.Errorf("failed to apply cure for %s: %w", result.TestID, err)
		}

		// Add visual separator after treatment
		fmt.Println()
	}

	return nil
}

// buildSummary aggregates results into a summary
func (r *Runner) buildSummary(results []types.DiagnosticResult) *types.Summary {
	summary := &types.Summary{
		Total:   len(results),
		Results: results,
	}

	for _, result := range results {
		switch result.Status {
		case types.StatusHealthy:
			summary.Healthy++
		case types.StatusInfo:
			summary.Info++
		case types.StatusWarning:
			summary.Warning++
		case types.StatusCritical:
			summary.Critical++
		}
	}

	return summary
}

// SetDiagnosticTimeout configures the timeout for diagnostic tests
func (r *Runner) SetDiagnosticTimeout(timeout time.Duration) {
	r.diagnosticTimeout = timeout
}

// SetCureTimeout configures the timeout for cure operations
func (r *Runner) SetCureTimeout(timeout time.Duration) {
	r.cureTimeout = timeout
}

// RunWithClaudeAssist processes diagnostic-cure pairs sequentially,
// spawning Claude sessions when cures fail
func (r *Runner) RunWithClaudeAssist(ctx context.Context, callback TreatmentCallback) (*types.Summary, error) {
	results := make([]types.DiagnosticResult, 0, len(r.config.Tests))

	for _, test := range r.config.Tests {
		fmt.Printf("\n🔍 Checking: %s\n", test.Description)

		// Run diagnostic
		result, err := r.runSingleDiagnostic(ctx, test)
		if err != nil {
			result = types.DiagnosticResult{
				TestID:      test.Test,
				Description: test.Description,
				Status:      types.StatusCritical,
				Summary:     fmt.Sprintf("Test execution failed: %v", err),
				Symptom:     test.Symptom,
				CureID:      test.Cure,
				FixAvailable: false,
				Severity:    test.Severity,
			}
		}
		results = append(results, result)

		// If healthy, move to next
		if result.Status == types.StatusHealthy {
			fmt.Println("  ✓ Healthy")
			continue
		}

		// If no cure available, skip
		if !result.FixAvailable || result.CureID == "" {
			fmt.Printf("  ⚠ %s (no automated cure available)\n", result.Status)
			continue
		}

		// Try automated cure
		fmt.Printf("\n💊 Applying cure: %s\n", result.CureID)
		if callback != nil {
			callback(result)
		}

		cureOutput := &bytes.Buffer{}
		cureErr := r.applyCureWithOutput(ctx, result.CureID, cureOutput)

		if cureErr == nil {
			// Verify cure worked
			verifyResult, err := r.runSingleDiagnostic(ctx, test)
			if err == nil && verifyResult.Status == types.StatusHealthy {
				fmt.Println("  ✓ Cure succeeded!")
				results[len(results)-1] = verifyResult
				continue
			}
		}

		// Cure failed - spawn Claude
		fmt.Println("\n  ✖ Automated cure failed")
		fmt.Println("  🤖 Spawning Claude to fix this issue...")

		claudeErr := r.spawnClaude(ctx, test, result, cureOutput.String())
		if claudeErr != nil {
			fmt.Printf("  ✖ Failed to spawn Claude: %v\n", claudeErr)
			continue
		}

		// Re-run diagnostic after Claude fixes it
		maxAttempts := 3
		for attempt := 1; attempt <= maxAttempts; attempt++ {
			fmt.Printf("\n🔍 Verifying fix (attempt %d/%d)...\n", attempt, maxAttempts)
			verifyResult, err := r.runSingleDiagnostic(ctx, test)
			if err == nil && verifyResult.Status == types.StatusHealthy {
				fmt.Println("  ✓ Issue resolved!")
				results[len(results)-1] = verifyResult
				break
			}

			if attempt < maxAttempts {
				fmt.Println("  ✖ Still failing, spawning Claude again...")
				claudeErr := r.spawnClaude(ctx, test, result, cureOutput.String())
				if claudeErr != nil {
					fmt.Printf("  ✖ Failed to spawn Claude: %v\n", claudeErr)
					break
				}
			} else {
				fmt.Println("  ✖ Max attempts reached, moving to next diagnostic")
			}
		}
	}

	return r.buildSummary(results), nil
}

// applyCureWithOutput runs a cure and captures its output
func (r *Runner) applyCureWithOutput(ctx context.Context, cureID string, output *bytes.Buffer) error {
	cureFunc, err := r.cureRegistry.Get(cureID)
	if err != nil {
		return fmt.Errorf("cure not found: %w", err)
	}

	// Redirect stdout to capture cure output
	oldStdout := os.Stdout
	pr, pw, _ := os.Pipe()
	os.Stdout = pw

	cureCtx, cancel := context.WithTimeout(ctx, r.cureTimeout)
	defer cancel()

	cureErr := cureFunc(cureCtx)

	// Restore stdout
	pw.Close()
	os.Stdout = oldStdout

	// Read captured output
	buf := make([]byte, 1024*1024) // 1MB buffer
	n, _ := pr.Read(buf)
	output.Write(buf[:n])

	return cureErr
}

// spawnClaude spawns a Claude session with context about the failure
func (r *Runner) spawnClaude(ctx context.Context, test types.DiagnosticTest, result types.DiagnosticResult, cureOutput string) error {
	// Create context message for Claude
	contextMsg := fmt.Sprintf(`dev-doctor failed to fix this issue automatically. Please help fix it.

## Diagnostic Information

**Test**: %s
**Description**: %s
**Status**: %s [%s]
**Symptom**: %s

## Automated Cure Attempted

**Cure ID**: %s
**Cure Output**:
%s

## Your Task

The automated cure failed. Please:
1. Analyze what went wrong
2. Fix the issue using alternative methods
3. Verify the fix works

When you're done, dev-doctor will re-run the diagnostic to verify.
`, test.Test, test.Description, result.Status, result.Severity, test.Symptom, result.CureID, cureOutput)

	// Check if claude CLI is available
	if _, err := exec.LookPath("claude"); err != nil {
		fmt.Println("\n⚠ Claude CLI not found. Please install it to enable automatic fixing:")
		fmt.Println("  https://docs.claude.com/claude-code")
		fmt.Println("\nManual fix required:")
		fmt.Println(contextMsg)
		return fmt.Errorf("claude CLI not available")
	}

	// Spawn Claude with context
	cmd := exec.CommandContext(ctx, "claude", "chat", "--message", contextMsg)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
