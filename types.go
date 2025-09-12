package validor

import (
	"context"
	"fmt"
	"sync"
	"testing"
)

// ModuleProcessor defines methods for processing Terraform modules
type ModuleProcessor interface {
	Apply(ctx context.Context, t *testing.T) error
	Destroy(ctx context.Context, t *testing.T) error
	CleanupFiles(t *testing.T) error
}

// TestResults holds the results of module tests in a thread-safe manner
type TestResults struct {
	mu            sync.RWMutex
	modules       []*Module
	failedModules []*Module
}

// NewTestResults creates a new TestResults instance
func NewTestResults() *TestResults {
	return &TestResults{
		modules:       make([]*Module, 0),
		failedModules: make([]*Module, 0),
	}
}

// AddModule adds a module to the results in a thread-safe way
func (tr *TestResults) AddModule(module *Module) {
	tr.mu.Lock()
	defer tr.mu.Unlock()
	tr.modules = append(tr.modules, module)
	if len(module.Errors) > 0 {
		tr.failedModules = append(tr.failedModules, module)
	}
}

// GetResults returns the modules and failed modules in a thread-safe way
func (tr *TestResults) GetResults() ([]*Module, []*Module) {
	tr.mu.RLock()
	defer tr.mu.RUnlock()
	return tr.modules, tr.failedModules
}

// ModuleInfo holds module and provider information
type ModuleInfo struct {
	Name      string
	Provider  string
	Namespace string
}

// FileRestore holds information needed to restore a file
type FileRestore struct {
	Path            string
	OriginalContent string
	ModuleName      string
	Provider        string
	Namespace       string
}

// TerraformRegistryResponse represents the API response structure
type TerraformRegistryResponse struct {
	Versions []struct {
		Version string `json:"version"`
	} `json:"versions"`
}

// ModuleError represents an error that occurred during module operations
type ModuleError struct {
	ModuleName string
	Operation  string
	Err        error
}

func (e *ModuleError) Error() string {
	return fmt.Sprintf("%s failed for module %s: %v", e.Operation, e.ModuleName, e.Err)
}

func (e *ModuleError) Unwrap() error {
	return e.Err
}
