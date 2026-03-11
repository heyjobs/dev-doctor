# dev-doctor

A health check and remediation tool for developer workstations.

## Overview

`dev-doctor` acts like a physician for your development environment. It runs diagnostic tests to identify common setup issues (expired AWS credentials, outdated tools, misconfigured services) and can optionally apply treatments to fix them automatically.

This reduces environment-related support requests and helps developers maintain healthy local setups.

## Features

- **Modular diagnostic system** - easily extensible with new checks
- **YAML-based configuration** - define tests declaratively
- **Interactive CLI** - polished user experience with colored output
- **Two consultation modes:**
  - **Diagnosis only** - identify issues without making changes
  - **Diagnosis + treatments** - identify and automatically fix issues
- **Mock diagnostics** - this first iteration uses mocked results to demonstrate the system

## Installation

### Prerequisites

- Go 1.21 or later

### Build from source

```bash
cd dev-doctor
go mod download
go build -o dev-doctor cmd/dev-doctor/main.go
```

### Run

```bash
./dev-doctor
```

Or install globally:

```bash
go install ./cmd/dev-doctor
```

## Usage

### Basic usage

Run with default configuration:

```bash
dev-doctor
```

### Specify custom configuration

```bash
dev-doctor --config /path/to/diagnostics.yaml
```

### Quiet mode

Suppress progress messages:

```bash
dev-doctor --quiet
```

## Project Structure

```
dev-doctor/
├── cmd/
│   └── dev-doctor/          # Main entry point
│       └── main.go
├── internal/
│   ├── cli/                 # CLI interface and user interaction
│   │   └── root.go
│   ├── config/              # Configuration loading and parsing
│   │   └── loader.go
│   ├── diagnostics/         # Diagnostic implementations
│   │   ├── registry.go      # Diagnostic registry
│   │   └── implementations.go
│   ├── cures/               # Treatment implementations
│   │   ├── registry.go      # Cure registry
│   │   └── implementations.go
│   ├── runner/              # Orchestrates diagnostic execution
│   │   └── runner.go
│   └── types/               # Core type definitions
│       └── types.go
├── configs/
│   └── diagnostics.yaml     # Test configuration
├── go.mod
├── go.sum
└── README.md
```

## Configuration

Diagnostics are defined in YAML format. Each test specifies:

- **test**: Unique identifier for the test
- **description**: Human-readable description
- **diagnostic**: Identifier for the diagnostic implementation
- **cure**: Identifier for the treatment implementation
- **severity**: `info`, `warning`, or `critical`
- **symptom**: Description of what the user experiences when this test fails

### Example configuration

```yaml
tests:
  - test: check_docker_runtime
    description: Verify Docker daemon is running
    diagnostic: docker_runtime
    cure: start_docker
    severity: critical
    symptom: Docker commands fail with daemon connection errors

  - test: check_aws_credentials
    description: Verify AWS credentials are active
    diagnostic: aws_credentials
    cure: refresh_aws_sso
    severity: warning
    symptom: AWS commands fail due to expired authentication
```

## Extending dev-doctor

The system is designed to be easily extensible. To add a new diagnostic:

### 1. Add YAML configuration

Edit `configs/diagnostics.yaml`:

```yaml
tests:
  - test: check_node_version
    description: Verify Node.js version meets requirements
    diagnostic: node_version
    cure: update_node
    severity: warning
    symptom: npm commands fail or produce unexpected results
```

### 2. Implement the diagnostic

Add to `internal/diagnostics/implementations.go`:

```go
func CheckNodeVersion(ctx context.Context) (types.Status, string, error) {
    // Execute: node --version
    // Parse version
    // Compare against requirements

    // Return status, summary message, and error
    return types.StatusHealthy, "Node.js v20.10.0 meets requirements", nil
}
```

### 3. Register the diagnostic

Add to `internal/diagnostics/registry.go` in the `DefaultRegistry()` function:

```go
reg.Register("node_version", CheckNodeVersion)
```

### 4. Implement the cure

Add to `internal/cures/implementations.go`:

```go
func UpdateNode(ctx context.Context) error {
    // Execute: brew upgrade node (or appropriate package manager)
    return nil
}
```

### 5. Register the cure

Add to `internal/cures/registry.go` in the `DefaultRegistry()` function:

```go
reg.Register("update_node", UpdateNode)
```

That's it! The CLI will automatically:
- Load your new test from YAML
- Execute the diagnostic
- Display results
- Apply the cure if requested

## Architecture

### Component Overview

1. **CLI Layer** (`internal/cli/`)
   - Handles user interaction
   - Formats output with colors and icons
   - Orchestrates the diagnostic flow

2. **Configuration** (`internal/config/`)
   - Loads YAML configuration
   - Validates schema
   - Provides defaults

3. **Type System** (`internal/types/`)
   - Core data structures
   - Function signatures
   - Enums for status and severity

4. **Diagnostic Registry** (`internal/diagnostics/`)
   - Maps diagnostic IDs to implementations
   - Provides lookup for the runner
   - Thread-safe registry

5. **Cure Registry** (`internal/cures/`)
   - Maps cure IDs to implementations
   - Provides lookup for treatment application
   - Thread-safe registry

6. **Runner** (`internal/runner/`)
   - Executes diagnostics with timeout
   - Collects and aggregates results
   - Applies treatments when requested

### Design Principles

- **Separation of concerns** - each component has a single responsibility
- **Registry pattern** - diagnostics and cures are registered by identifier
- **Declarative configuration** - tests are defined in YAML, not code
- **Extensibility** - adding new diagnostics requires minimal code changes
- **Type safety** - strong typing throughout the system

## Current Status

This is the **first iteration** focused on:
- Clean architecture
- Modular design
- Polished CLI experience
- YAML-based configuration
- Project structure

**Mock diagnostics** are currently in place. They return randomized results to simulate realistic scenarios.

## Next Steps

To move from mocked diagnostics to real implementations:

1. Replace mock implementations with real system checks
2. Add error handling for various failure modes
3. Implement actual remediation logic in cures
4. Add more comprehensive test coverage
5. Consider adding configuration for version requirements, paths, etc.
6. Add logging for troubleshooting
7. Add support for custom diagnostic plugins

## Example Output

```
╔════════════════════════════════════════════════╗
║           Welcome to dev-doctor                ║
╚════════════════════════════════════════════════╝

Running a diagnostic chart for your developer environment.

? Select consultation mode: Diagnosis only

Running Diagnostic Chart
────────────────────────

✔ Verify Docker daemon is running              [HEALTHY]
⚠ Verify AWS credentials are active            [WARNING]
  └─ AWS credentials expired 2 hours ago
✔ Validate AWS configuration file structure    [HEALTHY]
✔ Verify Git user configuration is set         [HEALTHY]
⚠ Verify OpenTofu version meets requirements   [WARNING]
  └─ OpenTofu version 1.4.0 is outdated

Chart Complete
──────────────

Total tests:     5
Healthy:         3
Warning:         2
Critical:        0

⚠ Some warnings detected. Consider addressing them to avoid future issues.
```

## Contributing

When adding new diagnostics:

1. Follow the extension guide above
2. Use meaningful test and diagnostic identifiers
3. Provide clear descriptions and symptom messages
4. Set appropriate severity levels
5. Implement both diagnostic and cure functions
6. Test thoroughly before committing

## License

[Add your license here]
