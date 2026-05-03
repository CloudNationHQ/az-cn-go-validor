package validor

import (
	"context"
)

// Logger is the minimal test logger contract used by runner abstractions.
type Logger interface {
	Helper()
	Logf(format string, args ...any)
	Log(args ...any)
}

// ModuleRunner applies, destroys, and cleans up a Terraform module.
type ModuleRunner interface {
	Apply(ctx context.Context, logger Logger) error
	Destroy(ctx context.Context, logger Logger) error
	Cleanup(ctx context.Context, logger Logger) error
}

// ModuleDiscoverer discovers Terraform example modules from a configured source.
type ModuleDiscoverer interface {
	DiscoverModules() ([]*Module, error)
	SetConfig(config *Config)
}

// SourceConverter converts Terraform module sources between registry and local paths.
type SourceConverter interface {
	ConvertToLocal(ctx context.Context, modulePath string, moduleInfo ModuleInfo) ([]FileRestore, error)
	RevertToRegistry(ctx context.Context, filesToRestore []FileRestore) error
}

// RegistryClient fetches Terraform Registry metadata.
type RegistryClient interface {
	GetLatestVersion(ctx context.Context, namespace, name, provider string) (string, error)
}

// TestRunner runs Terraform module tests.
type TestRunner interface {
	RunTests(logger Logger, modules []*Module, parallel bool, config *Config)
	RunLocalTests(logger Logger, examplesPath string) error
}
