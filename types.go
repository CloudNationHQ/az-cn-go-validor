package validor

import (
	"context"
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
