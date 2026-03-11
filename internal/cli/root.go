package cli

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/yourusername/dev-doctor/internal/config"
	"github.com/yourusername/dev-doctor/internal/runner"
	"github.com/yourusername/dev-doctor/internal/types"
)

var (
	configPath  string
	quietMode   bool
	modeFlag    string
	profileFlag string
)

// NewRootCommand creates the root command for dev-doctor
func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dev-doctor",
		Short: "Developer workstation health check and remediation tool",
		Long: `dev-doctor performs diagnostic tests on your development environment
and identifies common setup issues that can slow you down.

It checks for problems like expired AWS credentials, outdated tools,
misconfigured services, and more.`,
		RunE: runDiagnostics,
	}

	cmd.Flags().StringVarP(&configPath, "config", "c", "", "Path to diagnostics configuration file")
	cmd.Flags().BoolVarP(&quietMode, "quiet", "q", false, "Suppress progress messages")
	cmd.Flags().StringVarP(&modeFlag, "mode", "m", "", "Consultation mode: 'diagnosis' or 'treatment' (treatment spawns Claude if cures fail)")
	cmd.Flags().StringVarP(&profileFlag, "profile", "p", "", "Profile to run: 'basic', 'infrastructure', or 'data' (skips interactive prompt)")

	return cmd
}

// runDiagnostics is the main execution function
func runDiagnostics(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Print welcome message
	if !quietMode {
		printWelcome()
	}

	// Load configuration
	loader := config.NewLoader(configPath)
	cfg, err := loader.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Get profile (from flag or interactive prompt)
	var profile string
	if profileFlag != "" {
		// Validate profile from flag
		validProfiles := []string{"basic", "infrastructure", "data"}
		valid := false
		for _, vp := range validProfiles {
			if profileFlag == vp {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("invalid profile: %s (must be 'basic', 'infrastructure', or 'data')", profileFlag)
		}
		profile = profileFlag
	} else {
		// Prompt interactively
		var err error
		profile, err = promptProfile()
		if err != nil {
			return fmt.Errorf("failed to get profile: %w", err)
		}
	}

	// Filter configuration by profile
	cfg = config.FilterByProfile(cfg, profile)
	if len(cfg.Tests) == 0 {
		return fmt.Errorf("no tests found for profile: %s", profile)
	}

	// Get consultation mode (from flag or interactive prompt)
	var mode types.ConsultationMode
	if modeFlag != "" {
		// Use mode from flag
		switch modeFlag {
		case "diagnosis":
			mode = types.ModeDiagnosisOnly
		case "treatment":
			mode = types.ModeDiagnosisAndTreatment
		default:
			return fmt.Errorf("invalid mode: %s (must be 'diagnosis' or 'treatment')", modeFlag)
		}
	} else {
		// Prompt interactively
		var err error
		mode, err = promptConsultationMode()
		if err != nil {
			return fmt.Errorf("failed to get consultation mode: %w", err)
		}
		// Add immediate visual feedback after prompt
		fmt.Println()
	}

	// Create runner
	r := runner.NewRunner(cfg)
	var summary *types.Summary

	// Handle different modes
	if mode == types.ModeDiagnosisAndTreatment {
		// Treatment mode: show diagnostic chart first, then apply cures with Claude help
		if !quietMode {
			printSectionHeader("Running Diagnostic Chart")
			fmt.Println()
		}

		treatmentHeaderShown := false
		summary, err = r.RunWithClaudeAssist(ctx,
			// Diagnostic callback - show results as they complete
			func(result types.DiagnosticResult) {
				if !quietMode {
					printResult(result)
				}
			},
			// Treatment callback - show treatment headers
			func(result types.DiagnosticResult) {
				if !quietMode {
					if !treatmentHeaderShown {
						fmt.Println()
						printSectionHeader("Applying Treatments")
						fmt.Println()
						treatmentHeaderShown = true
					}
					printTreatmentHeader(result)
				}
			},
		)
		if err != nil {
			return fmt.Errorf("treatment execution failed: %w", err)
		}
	} else {
		// Diagnosis-only mode: just check without fixing
		if !quietMode {
			printSectionHeader("Running Diagnostic Chart")
			fmt.Println()
		}

		summary, err = r.RunDiagnosticsWithCallback(ctx, func(result types.DiagnosticResult) {
			if !quietMode {
				printResult(result)
			}
		})
		if err != nil {
			return fmt.Errorf("diagnostic execution failed: %w", err)
		}
	}

	// Print summary
	fmt.Println()
	printSummary(summary)

	return nil
}

// printWelcome displays the welcome message
func printWelcome() {
	banner := color.New(color.FgCyan, color.Bold)
	banner.Println("\n╔════════════════════════════════════════════════╗")
	banner.Println("║           Welcome to dev-doctor                ║")
	banner.Println("╚════════════════════════════════════════════════╝")

	fmt.Println("\nRunning a diagnostic chart for your developer environment.")
	fmt.Println()
}

// promptProfile asks the user to select a profile
func promptProfile() (string, error) {
	var profile string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select diagnostic profile").
				Description("Choose which tests to run based on your role").
				Options(
					huh.NewOption("Basic - core developer tools only", "basic"),
					huh.NewOption("Infrastructure - basic + platform tools (Docker, OpenTofu)", "infrastructure"),
					huh.NewOption("Data - basic + data engineering tools (Python)", "data"),
				).
				Value(&profile),
		),
	)

	err := form.Run()
	if err != nil {
		return "", err
	}

	return profile, nil
}

// promptConsultationMode asks the user to select a mode
func promptConsultationMode() (types.ConsultationMode, error) {
	var mode string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select consultation mode").
				Description("Choose how you want dev-doctor to operate").
				Options(
					huh.NewOption("Diagnosis only - identify issues without fixing", "diagnosis_only"),
					huh.NewOption("Treatment - fix issues automatically (spawns Claude if automated cures fail)", "diagnosis_and_treatment"),
				).
				Value(&mode),
		),
	)

	err := form.Run()
	if err != nil {
		return "", err
	}

	return types.ConsultationMode(mode), nil
}

// printSectionHeader prints a formatted section header
func printSectionHeader(title string) {
	header := color.New(color.FgYellow, color.Bold)
	header.Println(title)
	fmt.Println(strings.Repeat("─", len(title)))
}

// printResults displays diagnostic results (deprecated - results now shown in real-time)
func printResults(summary *types.Summary) {
	for _, result := range summary.Results {
		printResult(result)
	}
}

// printResult displays a single diagnostic result
func printResult(result types.DiagnosticResult) {
	var icon string
	var statusColor *color.Color

	switch result.Status {
	case types.StatusHealthy:
		icon = "✔"
		statusColor = color.New(color.FgGreen)
	case types.StatusInfo:
		icon = "ℹ"
		statusColor = color.New(color.FgBlue)
	case types.StatusWarning:
		icon = "⚠"
		statusColor = color.New(color.FgYellow)
	case types.StatusCritical:
		icon = "✖"
		statusColor = color.New(color.FgRed)
	}

	statusColor.Printf("%s ", icon)
	fmt.Printf("%-45s ", result.Description)
	statusColor.Printf("[%s]\n", strings.ToUpper(result.Status.String()))

	if result.Status != types.StatusHealthy {
		dimmed := color.New(color.Faint)
		dimmed.Printf("  └─ %s\n", result.Summary)
		if result.Symptom != "" {
			italic := color.New(color.Faint, color.Italic)
			italic.Printf("     Impact: %s\n", result.Symptom)
		}
	}
}

// printTreatmentHeader displays a visual header before each treatment
func printTreatmentHeader(result types.DiagnosticResult) {
	cyan := color.New(color.FgCyan, color.Bold)
	dimmed := color.New(color.Faint)

	// Print box border
	fmt.Println("┌" + strings.Repeat("─", 78) + "┐")

	// Print treatment title
	cyan.Printf("│ 💊 Treatment: %-63s │\n", result.CureID)

	// Print issue description
	fmt.Printf("│ Issue: %-68s │\n", result.Description)

	// Print bottom border
	fmt.Println("└" + strings.Repeat("─", 78) + "┘")
	fmt.Println()

	dimmed.Println("Applying treatment...")
	fmt.Println()
}

// printSummary displays the final summary
func printSummary(summary *types.Summary) {
	printSectionHeader("Chart Complete")
	fmt.Println()

	green := color.New(color.FgGreen)
	blue := color.New(color.FgBlue)
	yellow := color.New(color.FgYellow)
	red := color.New(color.FgRed)

	fmt.Printf("Total tests:     %d\n", summary.Total)
	green.Printf("Healthy:         %d\n", summary.Healthy)
	blue.Printf("Info:            %d\n", summary.Info)
	yellow.Printf("Warning:         %d\n", summary.Warning)
	red.Printf("Critical:        %d\n", summary.Critical)

	fmt.Println()

	if summary.Critical > 0 {
		red.Println("⚠ Critical issues detected. Your environment may not function correctly.")
	} else if summary.Warning > 0 {
		yellow.Println("⚠ Some warnings detected. Consider addressing them to avoid future issues.")
	} else if summary.Info > 0 {
		blue.Println("ℹ Some informational items detected. Review them when convenient.")
	} else {
		green.Println("✓ All systems healthy! Your development environment is in good shape.")
	}
}

// getUnhealthyResults filters results to only unhealthy ones
func getUnhealthyResults(results []types.DiagnosticResult) []types.DiagnosticResult {
	unhealthy := make([]types.DiagnosticResult, 0)
	for _, result := range results {
		if result.Status != types.StatusHealthy {
			unhealthy = append(unhealthy, result)
		}
	}
	return unhealthy
}

// Execute runs the CLI application
func Execute() {
	cmd := NewRootCommand()
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
