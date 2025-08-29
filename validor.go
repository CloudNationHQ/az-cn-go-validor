// Package validor provides testing utilities for Terraform modules.
// It offers functionality to apply, destroy and validate Terraform configurations
// in parallel or sequential mode with comprehensive error reporting.
package validor

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"
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

	// Get the expected module info from repository
	moduleInfo := extractModuleInfoFromRepo()
	if moduleInfo.Name == "" || moduleInfo.Provider == "" {
		t.Fatal(redError("Could not determine module name and provider from repository"))
	}

	// Convert all modules to use local source
	var allFilesToRestore []FileRestore
	for _, module := range modules {
		filesToRestore, err := convertToLocalSource(module.Path, moduleInfo)
		if err != nil {
			t.Logf("Warning: Failed to convert module %s to local source: %v", module.Name, err)
			continue
		}
		allFilesToRestore = append(allFilesToRestore, filesToRestore...)
	}

	// Ensure cleanup happens regardless of test outcome
	t.Cleanup(func() {
		if err := revertToRegistrySource(allFilesToRestore); err != nil {
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

// ModuleInfo holds module and provider information
type ModuleInfo struct {
	Name     string
	Provider string
}

// extractModuleInfoFromRepo extracts module name and provider from repository name
// Examples: terraform-azure-vnet -> {vnet, azure}, terraform-azuread-groups -> {groups, azuread}
func extractModuleInfoFromRepo() ModuleInfo {
	wd, err := os.Getwd()
	if err != nil {
		return ModuleInfo{}
	}

	// If we're in tests directory, go up one level to get repo name
	if filepath.Base(wd) == "tests" {
		wd = filepath.Dir(wd)
	}
	repoName := filepath.Base(wd)

	// Extract module name and provider from terraform-{PROVIDER}-{MODULE} pattern
	re := regexp.MustCompile(`^terraform-([^-]+)-(.+)$`)
	if matches := re.FindStringSubmatch(repoName); len(matches) > 2 {
		return ModuleInfo{
			Name:     matches[2], // MODULE part
			Provider: matches[1], // PROVIDER part
		}
	}

	// If pattern doesn't match, return empty (will cause test to fail with clear error)
	return ModuleInfo{}
}

// convertToLocalSource converts module blocks in Terraform files to use local source
func convertToLocalSource(modulePath string, moduleInfo ModuleInfo) ([]FileRestore, error) {
	var filesToRestore []FileRestore

	// Find all .tf files in the module path
	files, err := filepath.Glob(filepath.Join(modulePath, "*.tf"))
	if err != nil {
		return nil, err
	}

	modulePattern := fmt.Sprintf(`(?m)^(\s*module\s+"[^"]*"\s*\{[^}]*source\s*=\s*)"cloudnationhq/%s/%s"([^}]*version\s*=\s*"[^"]*")?([^}]*\})`, moduleInfo.Name, moduleInfo.Provider)
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

			moduleEnd = regexp.MustCompile(`(?m)^\s*version\s*=\s*"[^"]*"\s*\n?`).ReplaceAllString(moduleEnd, "")

			return fmt.Sprintf(`%s"../../"%s`, moduleStart, moduleEnd)
		})

		if newContent != originalContent {
			// Write modified content
			if err := os.WriteFile(file, []byte(newContent), 0644); err != nil {
				return filesToRestore, err
			}

			// Store restoration info
			filesToRestore = append(filesToRestore, FileRestore{
				Path:            file,
				OriginalContent: originalContent,
				ModuleName:      moduleInfo.Name,
				Provider:        moduleInfo.Provider,
			})
		}
	}

	return filesToRestore, nil
}

// getLatestModuleVersion fetches the latest version from Terraform Registry
func getLatestModuleVersion(namespace, name, provider string) (string, error) {
	url := fmt.Sprintf("https://registry.terraform.io/v1/modules/%s/%s/%s/versions", namespace, name, provider)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch module versions: HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var registryResp TerraformRegistryResponse
	if err := json.Unmarshal(body, &registryResp); err != nil {
		return "", err
	}

	if len(registryResp.Versions) == 0 {
		return "", fmt.Errorf("no versions found for module")
	}

	// Return the latest version (first in the list)
	return registryResp.Versions[0].Version, nil
}

// revertToRegistrySource reverts files back to registry source with latest version
func revertToRegistrySource(filesToRestore []FileRestore) error {
	for _, restore := range filesToRestore {
		// Get latest version from registry
		latestVersion, err := getLatestModuleVersion("cloudnationhq", restore.ModuleName, restore.Provider)
		if err != nil {
			// If we can't get the latest version, fall back to original content
			if writeErr := os.WriteFile(restore.Path, []byte(restore.OriginalContent), 0644); writeErr != nil {
				return writeErr
			}
			continue
		}

		// Create updated content with latest version
		updatedContent := restore.OriginalContent

		// Update version if it exists in original content
		versionRegex := regexp.MustCompile(`(version\s*=\s*")[^"]*(")`)
		if versionRegex.MatchString(updatedContent) {
			updatedContent = versionRegex.ReplaceAllString(updatedContent, fmt.Sprintf("${1}~> %s${2}", latestVersion))
		}

		if err := os.WriteFile(restore.Path, []byte(updatedContent), 0644); err != nil {
			return err
		}
	}
	return nil
}
