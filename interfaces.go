package validor

import (
	"context"
)

type Logger interface {
	Helper()
	Logf(format string, args ...any)
	Log(args ...any)
}

type ModuleRunner interface {
	Apply(ctx context.Context, logger Logger) error
	Destroy(ctx context.Context, logger Logger) error
	Cleanup(ctx context.Context, logger Logger) error
}

type ModuleDiscoverer interface {
	DiscoverModules() ([]*Module, error)
	SetConfig(config *Config)
}

type SourceConverter interface {
	ConvertToLocal(ctx context.Context, modulePath string, moduleInfo ModuleInfo) ([]FileRestore, error)
	RevertToRegistry(ctx context.Context, filesToRestore []FileRestore) error
}

type RegistryClient interface {
	GetLatestVersion(ctx context.Context, namespace, name, provider string) (string, error)
}

type TestRunner interface {
	RunTests(logger Logger, modules []*Module, parallel bool, config *Config)
	RunLocalTests(logger Logger, examplesPath string) error
}
