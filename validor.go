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
	ExceptionList map[string]bool
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
	config := GetConfig()
	config.ParseExceptionList()

	if config.Example == "" {
		t.Fatal(redError("-example flag is not set"))
	}

	exampleList := parseExampleList(config.Example)
	results := NewTestResults()

	for _, ex := range exampleList {
		if config.ExceptionList[ex] {
			t.Logf("Skipping example %s as it is in the exception list", ex)
			continue
		}

		t.Run(ex, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			modulePath := filepath.Join("..", "examples", ex)
			module := NewModule(ex, modulePath)

			if err := module.Apply(ctx, t); err != nil {
				t.Fail()
			} else {
				t.Logf("✓ Module %s applied successfully", module.Name)
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

// TestApplyAllParallel tests all Terraform modules in parallel
func TestApplyAllParallel(t *testing.T) {
	config := GetConfig()
	config.ParseExceptionList()

	manager := NewModuleManager(filepath.Join("..", "examples"))
	manager.SetConfig(config)
	modules, err := manager.DiscoverModules()
	if err != nil {
		errText := fmt.Sprintf("Failed to discover modules: %v", err)
		t.Fatal(redError(errText))
	}

	RunTests(t, modules, true, config)
}

// TestApplyAllSequential tests all Terraform modules sequentially
func TestApplyAllSequential(t *testing.T) {
	config := GetConfig()
	config.ParseExceptionList()

	manager := NewModuleManager(filepath.Join("..", "examples"))
	manager.SetConfig(config)
	modules, err := manager.DiscoverModules()
	if err != nil {
		errText := fmt.Sprintf("Failed to discover modules: %v", err)
		t.Fatal(redError(errText))
	}

	RunTests(t, modules, false, config)
}

// TestApplyAllLocal tests all Terraform modules with local source paths
func TestApplyAllLocal(t *testing.T) {
	ctx := context.Background()
	config := GetConfig()
	config.ParseExceptionList()

	manager := NewModuleManager(filepath.Join("..", "examples"))
	manager.SetConfig(config)
	modules, err := manager.DiscoverModules()
	if err != nil {
		errText := fmt.Sprintf("Failed to discover modules: %v", err)
		t.Fatal(redError(errText))
	}

	// Get the expected module info from repository
	moduleInfo := extractModuleInfoFromRepo()
	if moduleInfo.Name == "" || moduleInfo.Provider == "" {
		t.Fatal(redError("Could not determine module name and provider from repository"))
	}

	// Create converter with registry client
	converter := NewSourceConverter(NewRegistryClient())

	// Convert all modules to use local source
	var allFilesToRestore []FileRestore
	for _, module := range modules {
		filesToRestore, err := converter.ConvertToLocal(ctx, module.Path, moduleInfo)
		if err != nil {
			t.Logf("Warning: Failed to convert module %s to local source: %v", module.Name, err)
			continue
		}
		allFilesToRestore = append(allFilesToRestore, filesToRestore...)
	}

	t.Cleanup(func() {
		if err := converter.RevertToRegistry(context.Background(), allFilesToRestore); err != nil {
			t.Logf("Warning: Failed to revert files to registry source: %v", err)
		}
	})

	RunTests(t, modules, true, config)
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
