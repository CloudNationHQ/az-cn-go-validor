// Package validor provides testing utilities for Terraform modules.
package validor

import (
	"context"
	"fmt"
	"path/filepath"
	"slices"
	"testing"
)

func TestApplyNoError(t *testing.T, opts ...Option) {
	config := setupConfigWithOptions(opts...)
	if config.Example == "" {
		t.Fatal(redError("-example flag is not set"))
	}
	modules := createModulesFromNames(parseExampleList(config.Example), getExamplesPath(config))
	sourceType := BoolToStr(config.Local, "local", "registry")
	var setup TestSetupFunc
	if config.Local {
		setup = createLocalSetupFunc(config)
	}
	runModuleTests(t, modules, true, config, setup, sourceType)
}

func TestApplyAllParallel(t *testing.T, opts ...Option) {
	config := setupConfigWithOptions(opts...)
	modules := discoverModules(t, config)
	RunTests(t, modules, true, config)
}

func TestApplyAllSequential(t *testing.T, opts ...Option) {
	config := setupConfigWithOptions(opts...)
	modules := discoverModules(t, config)
	RunTests(t, modules, false, config)
}

func TestApplyAllLocal(t *testing.T, opts ...Option) {
	config := setupConfigWithOptions(opts...)
	modules := discoverModules(t, config)
	runModuleTests(t, modules, true, config, createLocalSetupFunc(config), "local")
}

func RunTests(t *testing.T, modules []*Module, parallel bool, config *Config) {
	runModuleTests(t, modules, parallel, config, nil, "registry")
}

func RunTestsWithOptions(t *testing.T, opts ...TestOption) {
	tc := &TestConfig{
		Parallel: true,
	}

	for _, opt := range opts {
		opt(tc)
	}

	if tc.Config == nil {
		tc.Config = NewConfig()
	}

	if tc.ExamplesPath != "" {
		tc.Config.ExamplesPath = tc.ExamplesPath
	}

	modules := createModulesFromNames(tc.ModuleNames, getExamplesPath(tc.Config))
	sourceType := BoolToStr(tc.UseLocal, "local", "registry")
	var setup TestSetupFunc
	if tc.UseLocal {
		setup = createLocalSetupFunc(tc.Config)
	}
	runModuleTestsFn(t, modules, tc.Parallel, tc.Config, setup, sourceType)
}

func runModuleTests(t *testing.T, modules []*Module, parallel bool, config *Config, setup TestSetupFunc, sourceType string) {
	ctx := context.Background()
	results := NewTestResults()

	if setup != nil {
		if err := setup(ctx, t, modules); err != nil {
			t.Fatal(redError(fmt.Sprintf("Setup failed: %v", err)))
		}
	}

	for _, module := range modules {
		if slices.Contains(config.ExceptionList, module.Name) {
			t.Logf("Skipping example %s as it is in the exception list", module.Name)
			continue
		}

		t.Run(module.Name, func(t *testing.T) {
			if parallel {
				t.Parallel()
			}

			if err := module.Apply(ctx, t); err != nil {
				t.Fail()
			} else {
				t.Logf("✓ Module %s applied successfully with %s source", module.Name, sourceType)
			}

			if !config.SkipDestroy {
				if err := module.Destroy(ctx, t); err != nil && !module.ApplyFailed {
					t.Logf("Cleanup failed for module %s: %v", module.Name, err)
				}
			}

			results.AddModule(module)
		})
	}

	t.Cleanup(func() {
		modules, _ := results.GetResults()
		PrintModuleSummary(t, modules)
	})
}

func discoverModules(t *testing.T, config *Config) []*Module {
	examplesPath := getExamplesPath(config)
	manager := NewModuleManager(examplesPath)
	manager.SetConfig(config)
	modules, err := manager.DiscoverModules()
	if err != nil {
		errText := fmt.Sprintf("Failed to discover modules: %v", err)
		t.Fatal(redError(errText))
	}
	return modules
}

func extractModuleNames(modules []*Module) []string {
	moduleNames := make([]string, 0, len(modules))
	for _, module := range modules {
		moduleNames = append(moduleNames, module.Name)
	}
	return moduleNames
}

func createModulesFromNames(moduleNames []string, basePath string) []*Module {
	modules := make([]*Module, 0, len(moduleNames))
	for _, name := range moduleNames {
		path := filepath.Join(basePath, name)
		modules = append(modules, NewModule(name, path))
	}
	return modules
}

var runModuleTestsFn = runModuleTests
