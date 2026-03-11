package cures

import (
	"fmt"
	"sync"

	"github.com/yourusername/dev-doctor/internal/types"
)

// Registry manages the mapping of cure identifiers to implementations
type Registry struct {
	mu    sync.RWMutex
	cures map[string]types.CureFunc
}

// NewRegistry creates a new cure registry
func NewRegistry() *Registry {
	return &Registry{
		cures: make(map[string]types.CureFunc),
	}
}

// Register adds a cure implementation to the registry
func (r *Registry) Register(id string, fn types.CureFunc) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.cures[id] = fn
}

// Get retrieves a cure implementation by ID
func (r *Registry) Get(id string) (types.CureFunc, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	fn, exists := r.cures[id]
	if !exists {
		return nil, fmt.Errorf("cure not found: %s", id)
	}
	return fn, nil
}

// List returns all registered cure identifiers
func (r *Registry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ids := make([]string, 0, len(r.cures))
	for id := range r.cures {
		ids = append(ids, id)
	}
	return ids
}

// DefaultRegistry returns a registry with all built-in cures registered
func DefaultRegistry() *Registry {
	reg := NewRegistry()

	// dbt analytics cures
	reg.Register("install_dbt_redshift", InstallDbtRedshift)
	reg.Register("setup_dbt_venv", SetupDbtVenv)
	reg.Register("setup_dbt_secret_config", SetupDbtSecretConfig)
	reg.Register("run_dbt_deps", RunDbtDeps)

	// Register all placeholder implementations
	reg.Register("configure_git", ConfigureGit)
	reg.Register("update_opentofu", UpdateOpenTofu)
	reg.Register("connect_to_vpn", ConnectToVPN)

	// Register real cure implementations
	reg.Register("update_python", UpdatePython)
	reg.Register("install_aws_vault", InstallAWSVault)
	reg.Register("setup_aws_sso", SetupAWSSSO)
	reg.Register("install_brewfile", InstallBrewfile)
	reg.Register("install_wasp", InstallWasp)
	reg.Register("install_docker", InstallDocker)
	reg.Register("start_docker", StartDocker)
	reg.Register("install_docker_compose", InstallDockerCompose)
	reg.Register("install_claude_cli", InstallClaudeCLI)
	reg.Register("install_homebrew", InstallHomebrew)
	reg.Register("update_homebrew", UpdateHomebrew)

	return reg
}
