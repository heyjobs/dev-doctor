# dev-doctor

> A health check and automated remediation tool for developer workstations

---

## What is dev-doctor?

**dev-doctor** is like a physician for your development environment. It:

1. **Runs diagnostic tests** to check your developer tools (Docker, Homebrew, Python, AWS, etc.)
2. **Identifies issues** (missing tools, outdated versions, broken configurations)
3. **Fixes problems automatically** with built-in automated cures

Think of it as `brew doctor` but for your entire development environment, not just Homebrew.

### The Problem It Solves

Developers waste time troubleshooting environment issues:
- "Docker commands aren't working"
- "My Python version is wrong"
- "AWS credentials expired"
- "Homebrew is outdated"

Instead of manually checking each tool, **dev-doctor runs all checks at once and can fix issues automatically**.

---

## Two Ways to Use dev-doctor

### 1. **Standalone Command** (Quick automated fixes)

Run `dev-doctor` directly for fast automated diagnostics and cures:

```bash
# Check status of all tools
dev-doctor --profile basic --mode diagnosis

# Check status AND automatically fix issues
dev-doctor --profile basic --mode treatment
```

**Best for:** Quick checks and automated fixes that work 95% of the time.

### 2. **Claude Skill** (Interactive fixing with AI assistance)

Run `/dev-doctor` inside Claude for interactive, intelligent problem-solving:

```bash
claude
> /dev-doctor
```

Claude will:
1. Ask which profile to check (basic/infrastructure/data)
2. Run diagnostics and show results
3. Ask which issue you want to fix
4. Try automated cure first
5. If cure fails, Claude manually fixes it with creative problem-solving

**Best for:** Complex issues where automated cures fail and you need intelligent help.

---

## Installation

### Prerequisites

- Go 1.21 or later
- macOS (currently tested on macOS 13+)

### Quick Install

```bash
git clone <repository-url>
cd dev-doctor
make install
```

This installs `dev-doctor` to `~/bin`. Add it to your PATH if needed:

```bash
echo 'export PATH="$HOME/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

### Alternative: System-wide Install

```bash
make install-global  # Installs to /usr/local/bin (requires sudo)
```

### Verify Installation

```bash
which dev-doctor
dev-doctor --help
```

---

## Usage

### Standalone Command Usage

#### Interactive Mode (Recommended)

Simply run:

```bash
dev-doctor
```

You'll be prompted to:
1. Select a profile (basic/infrastructure/data)
2. Choose mode (diagnosis or treatment)

#### Command-Line Mode

For scripts or automation:

```bash
# Diagnosis only (check status without fixing)
dev-doctor --profile basic --mode diagnosis

# Treatment mode (check and fix issues)
dev-doctor --profile infrastructure --mode treatment

# Custom config file
dev-doctor --config /path/to/config.yaml --profile data --mode treatment

# Quiet mode (minimal output)
dev-doctor --profile basic --mode diagnosis --quiet
```

### Claude Skill Usage

Inside Claude, run:

```bash
/dev-doctor
```

Follow Claude's interactive prompts:

```
Claude: Which profile would you like me to check?
- basic (Homebrew, AWS, VPN, Claude CLI)
- infrastructure (Docker, docker-compose, OpenTofu + basic)
- data (Python, Wasp + basic)

You: infrastructure

Claude: [Runs diagnostics and shows results]

Claude: I found 2 issues that need attention:
1. Docker is not installed [CRITICAL]
2. docker-compose is not installed [WARNING]

Which issue would you like me to help fix?

You: Docker

Claude: [Tries automated cure, or manually fixes if it fails]
```

---

## Profiles

dev-doctor organizes diagnostics into **profiles** based on your role:

### `basic` (Core Developer Tools)

**Who:** All developers
**Checks:**
- ✓ Homebrew installed and up to date
- ✓ VPN connection status
- ✓ AWS Vault installed
- ✓ AWS SSO configured
- ✓ Claude CLI installed

### `infrastructure` (Platform Engineering)

**Who:** DevOps, Platform Engineers
**Includes:** Everything in `basic` +
- ✓ Docker installed and running
- ✓ docker-compose available
- ✓ OpenTofu version requirements

### `data` (Data Engineering)

**Who:** Data Engineers, ML Engineers
**Includes:** Everything in `basic` +
- ✓ Python version (3.10.x)
- ✓ Wasp installed (for Redshift access)

---

## Example Output

### Diagnosis Mode

```
╔════════════════════════════════════════════════╗
║           Welcome to dev-doctor                ║
╚════════════════════════════════════════════════╝

Running a diagnostic chart for your developer environment.

Running Diagnostic Chart
────────────────────────

✔ Verify Homebrew is installed                  [HEALTHY]
✔ Verify Homebrew is up to date                 [HEALTHY]
ℹ Verify connected to VPN                       [INFO]
  └─ Not connected to VPN
     Impact: Cannot access internal resources without VPN
✖ Verify Docker is installed                    [CRITICAL]
  └─ Docker is not installed
     Impact: Cannot run Docker containers
⚠ Verify Python version is 3.10.x               [WARNING]
  └─ Python 3.9.7 does not meet requirements
     Impact: Python scripts may fail

Chart Complete
──────────────

Total tests:     5
Healthy:         2
Info:            1
Warning:         1
Critical:        1

✖ Critical issues detected. Your environment may not function correctly.
```

### Treatment Mode

```
Running Diagnostic Chart
────────────────────────
[Shows diagnostic results]

Applying Treatments
───────────────────

┌──────────────────────────────────────────────────────────────────────────────┐
│ 💊 Treatment: install_docker                                                  │
│ Issue: Verify Docker is installed                                           │
└──────────────────────────────────────────────────────────────────────────────┘

Applying treatment...

💊 Applying cure: install_docker

  Checking Homebrew installation...
  ✓ Homebrew is installed

  Installing Docker Desktop via Homebrew...
  [Installation output]

  ✓ Diagnostic now passing!

┌──────────────────────────────────────────────────────────────────────────────┐
│ 💊 Treatment: update_python                                                   │
│ Issue: Verify Python version is 3.10.x                                       │
└──────────────────────────────────────────────────────────────────────────────┘

[Continues with next treatment]
```

---

## How dev-doctor Works (Architecture)

### Modular Design

dev-doctor is built with modularity in mind. Each component is independent and extensible:

```
┌─────────────────────────────────────────────────────────────┐
│                         CLI Layer                           │
│  (User interaction, output formatting, prompts)             │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    Configuration Loader                      │
│  (Loads YAML config, embedded by default)                   │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                         Runner                              │
│  (Orchestrates diagnostic execution and cure application)   │
└─────────────────────────────────────────────────────────────┘
                              │
                 ┌────────────┴────────────┐
                 ▼                         ▼
┌──────────────────────────┐  ┌──────────────────────────┐
│  Diagnostic Registry     │  │     Cure Registry        │
│  (Maps IDs → Functions)  │  │  (Maps IDs → Functions)  │
└──────────────────────────┘  └──────────────────────────┘
                 │                         │
                 ▼                         ▼
┌──────────────────────────┐  ┌──────────────────────────┐
│  Diagnostic Functions    │  │     Cure Functions       │
│  - check_docker()        │  │  - install_docker()      │
│  - check_python()        │  │  - update_python()       │
│  - check_homebrew()      │  │  - update_homebrew()     │
│  - ...                   │  │  - ...                   │
└──────────────────────────┘  └──────────────────────────┘
```

### Key Components

#### 1. **YAML Configuration** (`configs/diagnostics.yaml`)

Defines all tests declaratively:

```yaml
tests:
  - test: check_docker_installed
    description: Verify Docker is installed
    diagnostic: docker_installed    # ID to look up function
    cure: install_docker           # ID to look up cure
    severity: critical
    symptom: Cannot run Docker containers
    profiles:
      - infrastructure
```

**Embedded in binary** - no separate config file needed for distribution.

#### 2. **Diagnostic Functions** (`internal/diagnostics/`)

Each diagnostic is a Go function that returns status:

```go
func CheckDockerInstalled(ctx context.Context) (types.Status, string, error) {
    // Check if Docker.app exists
    if _, err := os.Stat("/Applications/Docker.app"); err == nil {
        return types.StatusHealthy, "Docker Desktop is installed", nil
    }
    return types.StatusCritical, "Docker is not installed", nil
}
```

#### 3. **Cure Functions** (`internal/cures/`)

Each cure is a Go function that fixes the issue:

```go
func InstallDocker(ctx context.Context) error {
    fmt.Println("  Installing Docker Desktop via Homebrew...")
    cmd := exec.CommandContext(ctx, "brew", "install", "--cask", "docker")
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    return cmd.Run()
}
```

#### 4. **Registries** (Maps IDs → Functions)

The registries connect YAML config to Go functions:

```go
// Diagnostic Registry
func DefaultRegistry() *Registry {
    reg := NewRegistry()
    reg.Register("docker_installed", CheckDockerInstalled)
    reg.Register("python_version", CheckPythonVersion)
    // ...
    return reg
}

// Cure Registry
func DefaultRegistry() *Registry {
    reg := NewRegistry()
    reg.Register("install_docker", InstallDocker)
    reg.Register("update_python", UpdatePython)
    // ...
    return reg
}
```

#### 5. **Runner** (Execution Engine)

The runner:
1. Loads config
2. Looks up diagnostic functions
3. Executes them with timeout
4. Aggregates results
5. Applies cures when requested
6. Re-runs diagnostics to verify

---

## Adding New Diagnostics

Adding a new check is simple and requires **no changes to core code**:

### Step 1: Add to YAML Config

Edit `configs/diagnostics.yaml`:

```yaml
tests:
  - test: check_node_version
    description: Verify Node.js version meets requirements
    diagnostic: node_version
    cure: update_node
    severity: warning
    symptom: npm commands fail or produce unexpected results
    profiles:
      - basic
```

### Step 2: Implement Diagnostic Function

Create `internal/diagnostics/check_node_version.go`:

```go
package diagnostics

import (
    "context"
    "os/exec"
    "strings"
    "github.com/yourusername/dev-doctor/internal/types"
)

func CheckNodeVersion(ctx context.Context) (types.Status, string, error) {
    cmd := exec.CommandContext(ctx, "node", "--version")
    output, err := cmd.Output()
    if err != nil {
        return types.StatusCritical, "Node.js is not installed", nil
    }

    version := strings.TrimSpace(string(output))
    // Add version comparison logic here

    return types.StatusHealthy, version, nil
}
```

### Step 3: Register Diagnostic

Add to `internal/diagnostics/registry.go`:

```go
func DefaultRegistry() *Registry {
    reg := NewRegistry()
    // ... existing registrations ...
    reg.Register("node_version", CheckNodeVersion)
    return reg
}
```

### Step 4: Implement Cure Function

Create `internal/cures/update_node.go`:

```go
package cures

import (
    "context"
    "fmt"
    "os"
    "os/exec"
)

func UpdateNode(ctx context.Context) error {
    fmt.Println("  Updating Node.js via Homebrew...")
    cmd := exec.CommandContext(ctx, "brew", "upgrade", "node")
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    return cmd.Run()
}
```

### Step 5: Register Cure

Add to `internal/cures/registry.go`:

```go
func DefaultRegistry() *Registry {
    reg := NewRegistry()
    // ... existing registrations ...
    reg.Register("update_node", UpdateNode)
    return reg
}
```

### Step 6: Rebuild

```bash
make install
```

**That's it!** Your new diagnostic is now available.

---

## Understanding Status vs Severity

### Status (What the diagnostic returns)

- **HEALTHY** ✔ - Everything is working correctly
- **INFO** ℹ - FYI, no action needed (e.g., "VPN not connected")
- **WARNING** ⚠ - Issue detected, should be fixed
- **CRITICAL** ✖ - Serious issue, must be fixed

### Severity (Defined in YAML config)

Controls what status is returned when a check fails:

```yaml
severity: info       # Failed check → INFO status
severity: warning    # Failed check → WARNING status
severity: critical   # Failed check → CRITICAL status
```

### Treatment Behavior

- **INFO issues** - Shown in diagnostics, **never treated**
- **WARNING issues** - Treated in treatment mode
- **CRITICAL issues** - Treated in treatment mode

---

## Project Structure

```
dev-doctor/
├── cmd/
│   └── dev-doctor/              # Main entry point
│       └── main.go
├── internal/
│   ├── cli/                     # CLI interface, prompts, output
│   │   └── root.go
│   ├── config/                  # YAML loading, validation
│   │   ├── loader.go
│   │   └── embedded/
│   │       └── diagnostics.yaml # Embedded config
│   ├── diagnostics/             # All diagnostic implementations
│   │   ├── registry.go
│   │   ├── check_docker_installed.go
│   │   ├── check_docker_running.go
│   │   ├── check_homebrew_installed.go
│   │   └── ...
│   ├── cures/                   # All cure implementations
│   │   ├── registry.go
│   │   ├── install_docker.go
│   │   ├── start_docker.go
│   │   ├── install_homebrew.go
│   │   └── ...
│   ├── runner/                  # Execution orchestration
│   │   └── runner.go
│   └── types/                   # Core types and interfaces
│       └── types.go
├── configs/
│   └── diagnostics.yaml         # Source config (embedded during build)
├── .claude/
│   └── commands/
│       └── dev-doctor.md        # Claude skill definition
├── Makefile                     # Build and install commands
├── go.mod
└── README.md
```

---

## Makefile Commands

```bash
make help            # Show all commands
make build           # Build binary
make install         # Install to ~/bin
make install-global  # Install to /usr/local/bin (sudo)
make uninstall       # Remove from ~/bin
make uninstall-global# Remove from /usr/local/bin
make clean           # Remove build artifacts
make test            # Run tests
```

---

## Contributing

When adding new diagnostics:

1. **Follow the extension guide** above
2. **Use meaningful identifiers** - `check_docker_installed` not `test1`
3. **Write clear descriptions** - Users see these in output
4. **Set appropriate severity** - Info/Warning/Critical
5. **Provide helpful symptom messages** - Explain the impact
6. **Test thoroughly** - Both diagnostic and cure
7. **Handle errors gracefully** - Return helpful error messages

---

## FAQ

### Q: Why would I use this instead of just checking tools manually?

**A:** Manually checking takes time:
- Is Docker installed? Running?
- Is Homebrew up to date?
- Is my Python version correct?
- Are my AWS credentials valid?

dev-doctor checks **everything at once** and can **fix issues automatically**.

### Q: What's the difference between `--mode diagnosis` and `--mode treatment`?

**A:**
- **Diagnosis mode**: Only checks status, doesn't change anything
- **Treatment mode**: Checks status AND runs automated fixes

### Q: When should I use the Claude skill vs the command?

**A:**
- **Use command** for quick automated checks/fixes
- **Use Claude skill** when:
  - Automated cures fail
  - You need help understanding the issue
  - The fix requires creative problem-solving

### Q: Can I run this in CI/CD?

**A:** Yes! Use command-line flags for non-interactive mode:

```bash
dev-doctor --profile infrastructure --mode diagnosis --quiet
```

### Q: How do I add my own custom checks?

**A:** Follow the "Adding New Diagnostics" section above. It requires:
1. Adding YAML config
2. Writing diagnostic function
3. Writing cure function
4. Registering both in registries

### Q: Why are some issues marked INFO instead of WARNING?

**A:** INFO issues are informational only (e.g., "Not connected to VPN"). They're shown in diagnostics but **never automatically fixed** because they may be intentional.

---

## License

[Add your license here]
