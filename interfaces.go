package validor

import (
	"context"
	"testing"
)

// ModuleRunner defines the interface for running Terraform operations
type ModuleRunner interface {
	Apply(ctx context.Context, t *testing.T) error
	Destroy(ctx context.Context, t *testing.T) error
	Cleanup(ctx context.Context, t *testing.T) error
}

// ModuleDiscoverer defines the interface for discovering modules
type ModuleDiscoverer interface {
	DiscoverModules(ctx context.Context) ([]ModuleRunner, error)
	SetConfig(config *Config)
}

// SourceConverter defines the interface for converting module sources
type SourceConverter interface {
	ConvertToLocal(ctx context.Context, modulePath string, moduleInfo ModuleInfo) ([]FileRestore, error)
	RevertToRegistry(ctx context.Context, filesToRestore []FileRestore) error
}

// RegistryClient defines the interface for Terraform Registry operations
type RegistryClient interface {
	GetLatestVersion(ctx context.Context, namespace, name, provider string) (string, error)
}

// TestRunner defines the interface for running test scenarios
type TestRunner interface {
	RunTests(ctx context.Context, t *testing.T, modules []ModuleRunner, parallel bool, config *Config)
	RunLocalTests(ctx context.Context, t *testing.T, examplesPath string) error
}
