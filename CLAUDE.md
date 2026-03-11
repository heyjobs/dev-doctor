# Claude AI Assistant Configuration

This document defines the behavior and guidelines for Claude when working with the dev-doctor codebase.

## Project Overview

**dev-doctor** is a developer workstation health check and remediation tool for HeyJobs engineering teams. It performs diagnostic tests on development environments and automatically fixes common setup issues that slow down developers.

The tool follows a "diagnostic + cure" pattern where each check has an associated automated fix, making developer onboarding and environment maintenance seamless.

### Key Technologies
- **Go (Golang)** - Core implementation language
- **Cobra CLI** - Command-line interface framework
- **YAML Configuration** - Declarative diagnostic definitions
- **Homebrew** - macOS package management integration
- **GitHub CLI (gh)** - Repository integration

## Architecture

### Core Components

1. **Diagnostics** (`internal/diagnostics/`) - Health check implementations
   - Each diagnostic checks one aspect of the dev environment
   - Returns: Status (healthy/info/warning/critical), Summary message
   - Examples: check_docker_installed, check_vpn_connection, check_aws_sso_setup

2. **Cures** (`internal/cures/`) - Automated remediation implementations
   - Each cure fixes issues detected by diagnostics
   - Idempotent - safe to run multiple times
   - Graceful fallback - shows manual instructions if automatic fix fails
   - Examples: install_docker, connect_to_vpn, setup_aws_sso

3. **Configuration** (`configs/diagnostics.yaml`) - Test definitions
   - Maps diagnostics to cures
   - Defines severity levels (info/warning/critical)
   - Organizes tests into profiles (basic/infrastructure/data)
   - Contains symptoms (user-facing impact descriptions)

4. **Runner** (`internal/runner/`) - Orchestration engine
   - Executes diagnostics with timeouts
   - Applies cures sequentially
   - Maps severity to status for display

### Status vs Severity

**Important distinction:**
- **Severity** (in YAML config): How critical the issue is (info/warning/critical)
- **Status** (displayed to user): Derived from severity when check fails
- Healthy checks always show StatusHealthy regardless of configured severity

Mapping:
```
severity: info     → [INFO]     (blue)
severity: warning  → [WARNING]  (yellow)
severity: critical → [CRITICAL] (red)
```

### Profiles

Tests are organized into profiles based on developer roles:
- **basic** - Core tools everyone needs (Git, AWS, VPN)
- **infrastructure** - Platform engineers (Docker, OpenTofu, Brewfile)
- **data** - Data engineers (Python, Wasp)

## Development Commands

### Building
```bash
go build -o dev-doctor ./cmd/dev-doctor
```

### Running Locally
```bash
# Interactive mode (prompts for profile and mode)
./dev-doctor

# With flags (non-interactive)
./dev-doctor --profile infrastructure --mode treatment

# Available profiles: basic, infrastructure, data
# Available modes: diagnosis, treatment
```

### Testing Individual Components
```bash
# Run specific diagnostic
go run cmd/dev-doctor/main.go --profile basic --mode diagnosis

# Test cure implementation
# Cures are Go functions in internal/cures/
```

## Creating New Diagnostics and Cures

### Step 1: Create Diagnostic File
Create `internal/diagnostics/check_<name>.go`:
```go
package diagnostics

import (
    "context"
    "github.com/yourusername/dev-doctor/internal/types"
)

func CheckMyTool(ctx context.Context) (types.Status, string, error) {
    // Check logic here
    // Return StatusCritical/StatusWarning/StatusInfo if failed
    // Return StatusHealthy if passed
    return types.StatusHealthy, "Tool is installed", nil
}
```

### Step 2: Create Cure File
Create `internal/cures/<action>_<name>.go`:
```go
package cures

import (
    "context"
    "fmt"
)

func InstallMyTool(ctx context.Context) error {
    fmt.Println("  Installing my-tool...")
    // Installation logic here
    // Return nil on success
    // Show helpful error message on failure
    return nil
}
```

### Step 3: Register in Registries
Add to `internal/diagnostics/registry.go`:
```go
reg.Register("my_tool", CheckMyTool)
```

Add to `internal/cures/registry.go`:
```go
reg.Register("install_my_tool", InstallMyTool)
```

### Step 4: Add to Config
Add to `configs/diagnostics.yaml`:
```yaml
- test: check_my_tool
  description: Verify my-tool is installed
  diagnostic: my_tool
  cure: install_my_tool
  severity: warning
  symptom: Cannot use my-tool features
  profiles:
    - basic
```

## Cure Best Practices

### Graceful Error Handling
Cures should NEVER fail the entire treatment process. Instead:
```go
if err := someOperation(); err != nil {
    fmt.Println("  ✖ Automatic installation failed")
    fmt.Println()
    fmt.Println("  Please install manually:")
    fmt.Println("  https://example.com/install-guide")
    return nil  // Don't return error!
}
```

### Progress Feedback
Show clear progress indicators:
```go
fmt.Println("  Checking prerequisites...")
fmt.Println("  ✓ Prerequisites met")
fmt.Println()
fmt.Println("  Downloading package...")
fmt.Println("  ⏱  This may take several minutes...")
```

### Idempotency
Always check if already installed:
```go
if alreadyInstalled() {
    fmt.Println("  ✓ Tool already installed")
    return nil
}
```

## Branch Management

1. Work off `main` branch
2. Create feature branches for new work: `feature/<description>`
3. Push feature branches and create PRs to merge into main
4. Example: `feature/docker-auto-install-macos13`

## Skills

### dev-doctor Skill

The `dev-doctor` skill allows Claude to interactively run diagnostics and apply cures with intelligent problem-solving.

**Usage:**
```
/dev-doctor
```

**Behavior:**
1. Prompts user to select profile (basic/infrastructure/data)
2. Runs diagnostics for that profile
3. For each failing diagnostic:
   - Attempts to apply the associated cure
   - If cure succeeds: moves to next diagnostic
   - If cure fails: Claude analyzes the error and tries alternative solutions beyond the coded cure
4. Reports final status

**When cure fails, Claude should:**
- Analyze error messages to understand root cause
- Check system prerequisites
- Try alternative installation methods
- Search for known issues and solutions
- Provide step-by-step manual instructions if automation isn't possible
- Update the cure code if a better solution is found

This makes dev-doctor a collaborative debugging tool rather than just a scripted checker.
