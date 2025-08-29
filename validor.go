// Package validor provides testing utilities for Terraform modules.
// It offers functionality to apply, destroy and validate Terraform configurations
// in parallel or sequential mode with comprehensive error reporting.
package validor

import (
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
	config := GetConfig()
	config.ParseExceptionList()

	manager := NewModuleManager(filepath.Join("..", "examples"))
	manager.SetConfig(config)
	modules, err := manager.DiscoverModules()
	if err != nil {
		errText := fmt.Sprintf("Failed to discover modules: %v", err)
		t.Fatal(redError(errText))
	}

	// Get the expected module name from repository
	expectedModuleName := extractModuleNameFromRepo()
	if expectedModuleName == "" {
		t.Fatal(redError("Could not determine module name from repository"))
	}

	// Convert all modules to use local source
	var modifiedFiles []string
	for _, module := range modules {
		files, err := convertToLocalSource(module.Path, expectedModuleName)
		if err != nil {
			t.Logf("Warning: Failed to convert module %s to local source: %v", module.Name, err)
			continue
		}
		modifiedFiles = append(modifiedFiles, files...)
	}

	// Ensure cleanup happens regardless of test outcome
	t.Cleanup(func() {
		for _, file := range modifiedFiles {
			if err := revertLocalSource(file); err != nil {
				t.Logf("Warning: Failed to revert %s: %v", file, err)
			}
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

// extractModuleNameFromRepo extracts module name from repository name
// Examples: terraform-azure-vnet -> vnet, terraform-azure-sql -> sql
func extractModuleNameFromRepo() string {
	// Get current working directory name (repository name)
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}
	repoName := filepath.Base(wd)

	// Extract module name from terraform-azure-{MODULE} pattern
	re := regexp.MustCompile(`^terraform-azure-(.+)$`)
	if matches := re.FindStringSubmatch(repoName); len(matches) > 1 {
		return matches[1]
	}

	// If pattern doesn't match, return empty (will cause test to fail with clear error)
	return ""
}

// convertToLocalSource converts module blocks in Terraform files to use local source
func convertToLocalSource(modulePath, expectedModuleName string) ([]string, error) {
	var modifiedFiles []string

	// Find all .tf files in the module path
	files, err := filepath.Glob(filepath.Join(modulePath, "*.tf"))
	if err != nil {
		return nil, err
	}

	modulePattern := fmt.Sprintf(`(?m)^(\s*module\s+"[^"]*"\s*\{[^}]*source\s*=\s*)"cloudnationhq/%s/azure"([^}]*version\s*=\s*"[^"]*")?([^}]*\})`, expectedModuleName)
	re := regexp.MustCompile(modulePattern)

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			continue
		}

		originalContent := string(content)

		// Replace module source and remove version
		newContent := re.ReplaceAllStringFunc(originalContent, func(match string) string {
			// Extract the module block parts
			parts := re.FindStringSubmatch(match)
			if len(parts) < 4 {
				return match
			}

			moduleStart := parts[1]
			moduleEnd := parts[3]

			// Remove version line if present
			moduleEnd = regexp.MustCompile(`(?m)^\s*version\s*=\s*"[^"]*"\s*\n?`).ReplaceAllString(moduleEnd, "")

			return fmt.Sprintf(`%s"../../"%s`, moduleStart, moduleEnd)
		})

		if newContent != originalContent {
			// Create backup
			backupFile := file + ".backup"
			if err := os.WriteFile(backupFile, content, 0644); err != nil {
				return modifiedFiles, err
			}

			// Write modified content
			if err := os.WriteFile(file, []byte(newContent), 0644); err != nil {
				return modifiedFiles, err
			}

			modifiedFiles = append(modifiedFiles, file)
		}
	}

	return modifiedFiles, nil
}

// revertLocalSource reverts a file from its backup
func revertLocalSource(file string) error {
	backupFile := file + ".backup"

	// Read backup content
	content, err := os.ReadFile(backupFile)
	if err != nil {
		return err
	}

	// Restore original content
	if err := os.WriteFile(file, content, 0644); err != nil {
		return err
	}

	// Remove backup file
	return os.Remove(backupFile)
}
