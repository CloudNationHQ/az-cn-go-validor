package validor

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
)

// Module represents a Terraform module to test
type Module struct {
	Name        string
	Path        string
	Options     *terraform.Options
	Errors      []string
	ApplyFailed bool
}

// ModuleManager manages Terraform module discovery
type ModuleManager struct {
	BaseExamplesPath string
	Config           *Config
}

// NewModuleManager creates a new ModuleManager
func NewModuleManager(baseExamplesPath string) *ModuleManager {
	return &ModuleManager{
		BaseExamplesPath: baseExamplesPath,
	}
}

// SetConfig sets the configuration for the module manager
func (mm *ModuleManager) SetConfig(config *Config) {
	mm.Config = config
}

// NewModule creates a new Module instance
func NewModule(name, path string) *Module {
	return &Module{
		Name: name,
		Path: path,
		Options: &terraform.Options{
			TerraformDir: path,
			NoColor:      true,
		},
		Errors:      []string{},
		ApplyFailed: false,
	}
}

// DiscoverModules finds all Terraform modules in the examples directory
func (mm *ModuleManager) DiscoverModules() ([]*Module, error) {
	var modules []*Module

	entries, err := os.ReadDir(mm.BaseExamplesPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read examples directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			moduleName := entry.Name()
			if mm.Config != nil && mm.Config.ExceptionList[moduleName] {
				fmt.Printf("Skipping module %s as it is in the exception list\n", moduleName)
				continue
			}
			modulePath := filepath.Join(mm.BaseExamplesPath, moduleName)
			modules = append(modules, NewModule(moduleName, modulePath))
		}
	}

	return modules, nil
}

// Apply deploys a Terraform module with context support
func (m *Module) Apply(ctx context.Context, t *testing.T) error {
	t.Helper()

	if err := ctx.Err(); err != nil {
		return fmt.Errorf("context cancelled before applying module %s: %w", m.Name, err)
	}

	t.Logf("Applying Terraform module: %s", m.Name)
	terraform.WithDefaultRetryableErrors(t, m.Options)

	type result struct {
		err error
	}

	resultCh := make(chan result, 1)

	go func() {
		_, err := terraform.InitAndApplyE(t, m.Options)
		resultCh <- result{err: err}
	}()

	select {
	case res := <-resultCh:
		if res.err != nil {
			m.ApplyFailed = true
			wrappedErr := fmt.Errorf("terraform apply failed for module %s: %w", m.Name, res.err)
			m.Errors = append(m.Errors, wrappedErr.Error())
			t.Log(redError(wrappedErr.Error()))
			return wrappedErr
		}
		return nil
	case <-ctx.Done():
		return fmt.Errorf("context cancelled while applying module %s: %w", m.Name, ctx.Err())
	}
}

// Destroy tears down a deployed Terraform module with context support
func (m *Module) Destroy(ctx context.Context, t *testing.T) error {
	t.Helper()

	if err := ctx.Err(); err != nil {
		return fmt.Errorf("context cancelled before destroying module %s: %w", m.Name, err)
	}

	t.Logf("Destroying Terraform module: %s", m.Name)

	type result struct {
		err error
	}

	resultCh := make(chan result, 1)

	go func() {
		_, err := terraform.DestroyE(t, m.Options)
		resultCh <- result{err: err}
	}()

	var destroyErr error
	select {
	case res := <-resultCh:
		destroyErr = res.err
	case <-ctx.Done():
		return fmt.Errorf("context cancelled while destroying module %s: %w", m.Name, ctx.Err())
	}

	if destroyErr != nil && !m.ApplyFailed {
		wrappedErr := fmt.Errorf("terraform destroy failed for module %s: %w", m.Name, destroyErr)
		m.Errors = append(m.Errors, wrappedErr.Error())
		t.Log(redError(wrappedErr.Error()))
	}

	if err := m.CleanupFiles(t); err != nil && !m.ApplyFailed {
		wrappedErr := fmt.Errorf("cleanup failed for module %s: %w", m.Name, err)
		m.Errors = append(m.Errors, wrappedErr.Error())
		t.Log(redError(wrappedErr.Error()))
	}

	return destroyErr
}

// CleanupFiles removes Terraform-generated files after testing
func (m *Module) CleanupFiles(t *testing.T) error {
	t.Helper()
	t.Logf("Cleaning up in: %s", m.Options.TerraformDir)
	filesToCleanup := []string{"*.terraform*", "*tfstate*", "*.lock.hcl"}

	for _, pattern := range filesToCleanup {
		matches, err := filepath.Glob(filepath.Join(m.Options.TerraformDir, pattern))
		if err != nil {
			return fmt.Errorf("error matching pattern %s: %w", pattern, err)
		}
		for _, filePath := range matches {
			if err := os.RemoveAll(filePath); err != nil {
				return fmt.Errorf("failed to remove %s: %w", filePath, err)
			}
		}
	}
	return nil
}

// PrintModuleSummary prints a formatted summary of module test results
func PrintModuleSummary(t *testing.T, modules []*Module) {
	t.Helper()

	var failedModules []*Module
	for _, module := range modules {
		if len(module.Errors) > 0 {
			failedModules = append(failedModules, module)
		}
	}

	if len(failedModules) > 0 {
		// Print details for each failed module
		for _, module := range failedModules {
			t.Log(redError("Module " + module.Name + " failed with errors:"))
			for i, errMsg := range module.Errors {
				errText := fmt.Sprintf("  %d. %s", i+1, errMsg)
				t.Log(redError(errText))
			}
			t.Log("") // Empty line for better readability
		}

		// Print a count summary at the end
		totalText := fmt.Sprintf("TOTAL: %d of %d modules failed", len(failedModules), len(modules))
		t.Log(redError(totalText))
	} else {
		t.Logf("\n==== SUCCESS: All %d modules applied and destroyed successfully ====", len(modules))
	}
}
