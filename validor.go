// Package validor provides testing utilities for Terraform modules.
// It offers functionality to apply, destroy and validate Terraform configurations
// in parallel or sequential mode with comprehensive error reporting.
package validor

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

var (
	globalConfig *Config
)

// Config holds the configuration for validor tests
type Config struct {
	SkipDestroy   bool
	Exception     string
	Example       string
	Local         bool
	ExceptionList map[string]bool
}

// Option represents a functional option for Config
type Option func(*Config)

// WithSkipDestroy sets the skip destroy option
func WithSkipDestroy(skip bool) Option {
	return func(c *Config) { c.SkipDestroy = skip }
}

// WithException sets the exception list
func WithException(exception string) Option {
	return func(c *Config) {
		c.Exception = exception
		c.ParseExceptionList()
	}
}

// WithExample sets the example list
func WithExample(example string) Option {
	return func(c *Config) { c.Example = example }
}

// WithLocal sets the local testing option
func WithLocal(local bool) Option {
	return func(c *Config) { c.Local = local }
}

// NewConfig creates a new Config with the provided options
func NewConfig(opts ...Option) *Config {
	config := &Config{}
	for _, opt := range opts {
		opt(config)
	}
	return config
}

// ModuleInfo holds module and provider information
type ModuleInfo struct {
	Name     string
	Provider string
}

// FileRestore holds information needed to restore a file
type FileRestore struct {
	Path            string
	OriginalContent string
	ModuleName      string
	Provider        string
}

// TerraformRegistryResponse represents the API response structure
type TerraformRegistryResponse struct {
	Versions []struct {
		Version string `json:"version"`
	} `json:"versions"`
}

// init registers flags before the test framework parses command line arguments
func init() {
	globalConfig = &Config{}
	flag.BoolVar(&globalConfig.SkipDestroy, "skip-destroy", false, "Skip running terraform destroy after apply")
	flag.StringVar(&globalConfig.Exception, "exception", "", "Comma-separated list of examples to exclude")
	flag.StringVar(&globalConfig.Example, "example", "", "Specific example(s) to test (comma-separated)")
	flag.BoolVar(&globalConfig.Local, "local", false, "Use local source for testing")
}

// GetConfig returns the global configuration instance
func GetConfig() *Config {
	return globalConfig
}

// ParseExceptionList converts a comma-separated list to a map
func (c *Config) ParseExceptionList() {
	c.ExceptionList = make(map[string]bool)
	if c.Exception == "" {
		return
	}
	for _, ex := range strings.FieldsFunc(c.Exception, func(r rune) bool { return r == ',' }) {
		c.ExceptionList[strings.TrimSpace(ex)] = true
	}
}

// TestApplyNoError tests one or more specific Terraform modules
func TestApplyNoError(t *testing.T) {
	config := setupConfig()
	if config.Example == "" {
		t.Fatal(redError("-example flag is not set"))
	}
	runTestsWithConfig(t, config, parseExampleList(config.Example), config.Local)
}

// TestApplyAllParallel tests all Terraform modules in parallel
func TestApplyAllParallel(t *testing.T) {
	config := setupConfig()
	modules := discoverModules(t, config)
	RunTests(t, modules, true, config)
}

// TestApplyAllSequential tests all Terraform modules sequentially
func TestApplyAllSequential(t *testing.T) {
	config := setupConfig()
	modules := discoverModules(t, config)
	RunTests(t, modules, false, config)
}

// TestApplyAllLocal tests all Terraform modules with local source paths
func TestApplyAllLocal(t *testing.T) {
	config := setupConfig()
	modules := discoverModules(t, config)
	moduleNames := extractModuleNames(modules)
	runTestsWithConfig(t, config, moduleNames, true)
}

// TestOption represents a functional option for test execution
type TestOption func(*TestConfig)

// TestConfig holds test execution configuration
type TestConfig struct {
	Config      *Config
	ModuleNames []string
	UseLocal    bool
	Parallel    bool
}

// WithConfig sets the validor configuration
func WithConfig(config *Config) TestOption {
	return func(tc *TestConfig) { tc.Config = config }
}

// WithModules sets the module names to test
func WithModules(moduleNames []string) TestOption {
	return func(tc *TestConfig) { tc.ModuleNames = moduleNames }
}

// WithLocalSource enables local source testing
func WithLocalSource(useLocal bool) TestOption {
	return func(tc *TestConfig) { tc.UseLocal = useLocal }
}

// WithParallel enables parallel test execution
func WithParallel(parallel bool) TestOption {
	return func(tc *TestConfig) { tc.Parallel = parallel }
}

// RunTestsWithOptions executes tests with the provided options
func RunTestsWithOptions(t *testing.T, opts ...TestOption) {
	tc := &TestConfig{
		Parallel: true, // default to parallel
	}
	
	for _, opt := range opts {
		opt(tc)
	}
	
	if tc.Config == nil {
		tc.Config = GetConfig()
		tc.Config.ParseExceptionList()
	}
	
	runTestsWithConfig(t, tc.Config, tc.ModuleNames, tc.UseLocal)
}

// runTestsWithConfig runs tests for specific modules with optional local source conversion
func runTestsWithConfig(t *testing.T, config *Config, moduleNames []string, useLocal bool) {
	ctx := context.Background()
	results := NewTestResults()
	
	var converter SourceConverter
	var allFilesToRestore []FileRestore

	// Setup local source conversion if requested
	if useLocal {
		moduleInfo := extractModuleInfoFromRepo()
		if moduleInfo.Name == "" || moduleInfo.Provider == "" {
			t.Fatal(redError("Could not determine module name and provider from repository"))
		}

		converter = NewSourceConverter(NewRegistryClient())
		allFilesToRestore = convertModulesToLocal(ctx, t, converter, moduleNames, config.ExceptionList, moduleInfo)

		// Ensure cleanup happens regardless of test outcome
		t.Cleanup(func() {
			if err := converter.RevertToRegistry(context.Background(), allFilesToRestore); err != nil {
				t.Logf("Warning: Failed to revert files to registry source: %v", err)
			}
		})
	}

	// Run tests for each module
	for _, moduleName := range moduleNames {
		if config.ExceptionList[moduleName] {
			t.Logf("Skipping example %s as it is in the exception list", moduleName)
			continue
		}

		t.Run(moduleName, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			modulePath := filepath.Join("..", "examples", moduleName)
			module := NewModule(moduleName, modulePath)

			sourceType := map[bool]string{true: "local", false: "registry"}[useLocal]

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

// setupConfig initializes and returns the global configuration
func setupConfig() *Config {
	config := GetConfig()
	config.ParseExceptionList()
	return config
}

// discoverModules discovers all modules in the examples directory
func discoverModules(t *testing.T, config *Config) []*Module {
	manager := NewModuleManager(filepath.Join("..", "examples"))
	manager.SetConfig(config)
	modules, err := manager.DiscoverModules()
	if err != nil {
		errText := fmt.Sprintf("Failed to discover modules: %v", err)
		t.Fatal(redError(errText))
	}
	return modules
}

// extractModuleNames extracts module names from discovered modules
func extractModuleNames(modules []*Module) []string {
	var moduleNames []string
	for _, module := range modules {
		moduleNames = append(moduleNames, module.Name)
	}
	return moduleNames
}

// convertModulesToLocal converts specified modules to local source paths
func convertModulesToLocal(ctx context.Context, t *testing.T, converter SourceConverter, moduleNames []string, exceptionList map[string]bool, moduleInfo ModuleInfo) []FileRestore {
	var allFilesToRestore []FileRestore
	
	for _, moduleName := range moduleNames {
		if exceptionList[moduleName] {
			continue
		}

		modulePath := filepath.Join("..", "examples", moduleName)
		filesToRestore, err := converter.ConvertToLocal(ctx, modulePath, moduleInfo)
		if err != nil {
			t.Logf("Warning: Failed to convert module %s to local source: %v", moduleName, err)
			continue
		}
		allFilesToRestore = append(allFilesToRestore, filesToRestore...)
	}
	
	return allFilesToRestore
}

// parseExampleList parses a comma-separated list of examples
func parseExampleList(example string) []string {
	var examples []string
	for ex := range strings.SplitSeq(example, ",") {
		if trimmed := strings.TrimSpace(ex); trimmed != "" {
			examples = append(examples, trimmed)
		}
	}
	return examples
}

// extractModuleInfoFromRepo extracts module name and provider from repository name
func extractModuleInfoFromRepo() ModuleInfo {
	wd, err := os.Getwd()
	if err != nil {
		return ModuleInfo{}
	}

	if filepath.Base(wd) == "tests" {
		wd = filepath.Dir(wd)
	}
	repoName := filepath.Base(wd)

	re := regexp.MustCompile(`^terraform-([^-]+)-(.+)$`)
	if matches := re.FindStringSubmatch(repoName); len(matches) > 2 {
		return ModuleInfo{
			Name:     matches[2],
			Provider: matches[1],
		}
	}
	return ModuleInfo{}
}
