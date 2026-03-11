package types

import "context"

// Status represents the health status of a diagnostic test
type Status string

const (
	StatusHealthy  Status = "healthy"
	StatusInfo     Status = "info"
	StatusWarning  Status = "warning"
	StatusCritical Status = "critical"
)

// String returns the string representation of the status
func (s Status) String() string {
	return string(s)
}

// Severity represents how critical a failed test is
type Severity string

const (
	SeverityInfo     Severity = "info"
	SeverityWarning  Severity = "warning"
	SeverityCritical Severity = "critical"
)

// DiagnosticTest represents a test definition from YAML configuration
type DiagnosticTest struct {
	Test        string   `yaml:"test"`
	Description string   `yaml:"description"`
	Diagnostic  string   `yaml:"diagnostic"`
	Cure        string   `yaml:"cure"`
	Severity    Severity `yaml:"severity"`
	Symptom     string   `yaml:"symptom"`
	Profiles    []string `yaml:"profiles"`
}

// DiagnosticResult represents the result of executing a diagnostic test
type DiagnosticResult struct {
	TestID       string
	Description  string
	Status       Status
	Summary      string
	Symptom      string
	CureID       string
	FixAvailable bool
	Severity     Severity
}

// DiagnosticFunc is the function signature for diagnostic implementations
type DiagnosticFunc func(ctx context.Context) (Status, string, error)

// CureFunc is the function signature for cure implementations
type CureFunc func(ctx context.Context) error

// ConsultationMode represents the operating mode selected by the user
type ConsultationMode string

const (
	ModeDiagnosisOnly         ConsultationMode = "diagnosis_only"
	ModeDiagnosisAndTreatment ConsultationMode = "diagnosis_and_treatment"
)

// DiagnosticConfig represents the root configuration structure
type DiagnosticConfig struct {
	Tests []DiagnosticTest `yaml:"tests"`
}

// Summary represents aggregated results from a diagnostic run
type Summary struct {
	Total     int
	Healthy   int
	Info      int
	Warning   int
	Critical  int
	Results   []DiagnosticResult
}
