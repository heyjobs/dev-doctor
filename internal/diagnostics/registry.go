package diagnostics

import (
	"fmt"
	"sync"

	"github.com/yourusername/dev-doctor/internal/types"
)

// Registry manages the mapping of diagnostic identifiers to implementations
type Registry struct {
	mu          sync.RWMutex
	diagnostics map[string]types.DiagnosticFunc
}

// NewRegistry creates a new diagnostic registry
func NewRegistry() *Registry {
	return &Registry{
		diagnostics: make(map[string]types.DiagnosticFunc),
	}
}

// Register adds a diagnostic implementation to the registry
func (r *Registry) Register(id string, fn types.DiagnosticFunc) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.diagnostics[id] = fn
}

// Get retrieves a diagnostic implementation by ID
func (r *Registry) Get(id string) (types.DiagnosticFunc, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	fn, exists := r.diagnostics[id]
	if !exists {
		return nil, fmt.Errorf("diagnostic not found: %s", id)
	}
	return fn, nil
}

// List returns all registered diagnostic identifiers
func (r *Registry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ids := make([]string, 0, len(r.diagnostics))
	for id := range r.diagnostics {
		ids = append(ids, id)
	}
	return ids
}

// DefaultRegistry returns a registry with all built-in diagnostics registered
func DefaultRegistry() *Registry {
	reg := NewRegistry()

	// Register all mock implementations
	reg.Register("git_config", CheckGitConfiguration)
	reg.Register("opentofu_version", CheckOpenTofuVersion)

	// Register real implementations
	reg.Register("python_version", CheckPythonVersion)
	reg.Register("aws_vault", CheckAWSVault)
	reg.Register("aws_sso_setup", CheckAWSSSOSetup)
	reg.Register("brewfile", CheckBrewfile)
	reg.Register("vpn_connection", CheckVPNConnection)
	reg.Register("wasp_version", CheckWaspVersion)
	reg.Register("valde_test", CheckValdeTest)
	reg.Register("valde_warning", CheckValdeWarning)
	reg.Register("valde_critical", CheckValdeCritical)

	return reg
}
