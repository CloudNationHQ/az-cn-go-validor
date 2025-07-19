// Package validor provides testing utilities for Terraform modules.
// It offers functionality to apply, destroy and validate Terraform configurations
// in parallel or sequential mode with comprehensive error reporting.
package validor

import (
	"flag"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

var (
	configOnce   sync.Once
	globalConfig *Config
)

// Config holds the configuration for validor tests
type Config struct {
	SkipDestroy   bool
	Exception     string
	Example       string
	ExceptionList map[string]bool
}

// GetConfig returns a singleton configuration instance
func GetConfig() *Config {
	configOnce.Do(func() {
		globalConfig = &Config{}
		flag.BoolVar(&globalConfig.SkipDestroy, "skip-destroy", false, "Skip running terraform destroy after apply")
		flag.StringVar(&globalConfig.Exception, "exception", "", "Comma-separated list of examples to exclude")
		flag.StringVar(&globalConfig.Example, "example", "", "Specific example(s) to test (comma-separated)")
	})
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
	flag.Parse()
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
			modulePath := filepath.Join("..", "examples", ex)
			module := NewModule(ex, modulePath)

			if err := module.Apply(t); err != nil {
				t.Fail()
			} else {
				t.Logf("✓ Module %s applied successfully", module.Name)
			}

			if !config.SkipDestroy {
				if err := module.Destroy(t); err != nil && !module.ApplyFailed {
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
	flag.Parse()
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
