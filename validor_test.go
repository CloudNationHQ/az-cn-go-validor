package validor

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestParseExampleList(t *testing.T) {
	tests := []struct {
		name    string
		example string
		want    []string
	}{
		{
			name:    "single example",
			example: "example1",
			want:    []string{"example1"},
		},
		{
			name:    "multiple examples",
			example: "example1,example2,example3",
			want:    []string{"example1", "example2", "example3"},
		},
		{
			name:    "examples with spaces",
			example: " example1 , example2 , example3 ",
			want:    []string{"example1", "example2", "example3"},
		},
		{
			name:    "examples with trailing comma",
			example: "example1,example2,",
			want:    []string{"example1", "example2"},
		},
		{
			name:    "examples with empty entries",
			example: "example1,,example2",
			want:    []string{"example1", "example2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseExampleList(tt.example)
			// Handle nil vs empty slice comparison
			if len(got) == 0 && len(tt.want) == 0 {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseExampleList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractModuleInfoFromRepo(t *testing.T) {
	tests := []struct {
		name     string
		repoName string
		want     ModuleInfo
	}{
		{
			name:     "valid terraform-azure module",
			repoName: "terraform-azure-mymodule",
			want: ModuleInfo{
				Name:     "mymodule",
				Provider: "azure",
			},
		},
		{
			name:     "valid terraform-aws module",
			repoName: "terraform-aws-vpc",
			want: ModuleInfo{
				Name:     "vpc",
				Provider: "aws",
			},
		},
		{
			name:     "module with hyphenated name",
			repoName: "terraform-azure-storage-account",
			want: ModuleInfo{
				Name:     "storage-account",
				Provider: "azure",
			},
		},
		{
			name:     "invalid format - no terraform prefix",
			repoName: "azure-mymodule",
			want:     ModuleInfo{},
		},
		{
			name:     "invalid format - no provider",
			repoName: "terraform-mymodule",
			want:     ModuleInfo{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary directory with the repo name
			tmpDir := t.TempDir()
			repoDir := filepath.Join(tmpDir, tt.repoName)
			if err := os.Mkdir(repoDir, 0755); err != nil {
				t.Fatalf("Failed to create test directory: %v", err)
			}

			// Change to the repo directory
			originalWd, err := os.Getwd()
			if err != nil {
				t.Fatalf("Failed to get current directory: %v", err)
			}
			defer os.Chdir(originalWd)

			if err := os.Chdir(repoDir); err != nil {
				t.Fatalf("Failed to change to test directory: %v", err)
			}

			got := extractModuleInfoFromRepo()

			if got.Name != tt.want.Name || got.Provider != tt.want.Provider {
				t.Errorf("extractModuleInfoFromRepo() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestExtractModuleInfoFromRepo_WithTestsSubdir(t *testing.T) {
	// Create a temporary directory structure
	tmpDir := t.TempDir()
	repoName := "terraform-azure-testmodule"
	repoDir := filepath.Join(tmpDir, repoName)
	testsDir := filepath.Join(repoDir, "tests")

	if err := os.MkdirAll(testsDir, 0755); err != nil {
		t.Fatalf("Failed to create test directories: %v", err)
	}

	// Change to the tests directory
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalWd)

	if err := os.Chdir(testsDir); err != nil {
		t.Fatalf("Failed to change to test directory: %v", err)
	}

	got := extractModuleInfoFromRepo()

	want := ModuleInfo{
		Name:     "testmodule",
		Provider: "azure",
	}

	if got.Name != want.Name || got.Provider != want.Provider {
		t.Errorf("extractModuleInfoFromRepo() from tests subdir = %+v, want %+v", got, want)
	}
}

func TestGetRepoNameFromGit(t *testing.T) {
	// This test requires a git repository, so we'll skip it if not in a git repo
	tmpDir := t.TempDir()

	t.Run("non-git directory", func(t *testing.T) {
		result := getRepoNameFromGit(tmpDir)
		if result != "" {
			t.Errorf("getRepoNameFromGit() for non-git dir should return empty string, got %v", result)
		}
	})
}

func TestTestConfig_Options(t *testing.T) {
	t.Run("WithConfig", func(t *testing.T) {
		config := &Config{Example: "test"}
		tc := &TestConfig{}
		WithConfig(config)(tc)
		if tc.Config != config {
			t.Error("WithConfig did not set Config correctly")
		}
	})

	t.Run("WithModules", func(t *testing.T) {
		modules := []string{"mod1", "mod2"}
		tc := &TestConfig{}
		WithModules(modules)(tc)
		if !reflect.DeepEqual(tc.ModuleNames, modules) {
			t.Error("WithModules did not set ModuleNames correctly")
		}
	})

	t.Run("WithLocalSource", func(t *testing.T) {
		tc := &TestConfig{}
		WithLocalSource(true)(tc)
		if !tc.UseLocal {
			t.Error("WithLocalSource did not set UseLocal correctly")
		}
	})

	t.Run("WithParallel", func(t *testing.T) {
		tc := &TestConfig{}
		WithParallel(false)(tc)
		if tc.Parallel {
			t.Error("WithParallel did not set Parallel correctly")
		}
	})

	t.Run("WithTestExamplesPath", func(t *testing.T) {
		tc := &TestConfig{}
		WithTestExamplesPath("/test/path")(tc)
		if tc.ExamplesPath != "/test/path" {
			t.Error("WithTestExamplesPath did not set ExamplesPath correctly")
		}
	})
}

func TestSetupConfigWithOptions(t *testing.T) {
	// Reset global config
	originalConfig := globalConfig
	defer func() { globalConfig = originalConfig }()

	globalConfig = &Config{
		Exception: "ex1,ex2",
	}

	t.Run("apply options to global config", func(t *testing.T) {
		config := setupConfigWithOptions(
			WithSkipDestroy(true),
			WithLocal(true),
		)

		if !config.SkipDestroy {
			t.Error("SkipDestroy should be true")
		}
		if !config.Local {
			t.Error("Local should be true")
		}
		// ExceptionList should be parsed
		if len(config.ExceptionList) != 2 {
			t.Errorf("ExceptionList should have 2 items, got %d", len(config.ExceptionList))
		}
	})
}

func TestConvertModulesToLocal(t *testing.T) {
	// Create test directory structure
	tmpDir := t.TempDir()
	examplesDir := filepath.Join(tmpDir, "examples")
	if err := os.MkdirAll(examplesDir, 0755); err != nil {
		t.Fatalf("Failed to create examples directory: %v", err)
	}

	// Create test modules
	moduleNames := []string{"example1", "example2"}
	for _, modName := range moduleNames {
		modDir := filepath.Join(examplesDir, modName)
		if err := os.Mkdir(modDir, 0755); err != nil {
			t.Fatalf("Failed to create module directory: %v", err)
		}

		// Create a simple terraform file
		tfContent := `
module "test" {
  source  = "cloudnationhq/mymodule/azure"
  version = "~> 1.0"
}
`
		tfFile := filepath.Join(modDir, "main.tf")
		if err := os.WriteFile(tfFile, []byte(tfContent), 0644); err != nil {
			t.Fatalf("Failed to create terraform file: %v", err)
		}
	}

	client := &mockRegistryClient{latestVersion: "1.0.0"}
	converter := NewSourceConverter(client)
	moduleInfo := ModuleInfo{
		Name:      "mymodule",
		Provider:  "azure",
		Namespace: "cloudnationhq",
	}

	ctx := testContext(t)
	mockT := &testing.T{}
	filesToRestore := convertModulesToLocal(ctx, mockT, converter, moduleNames, []string{}, moduleInfo, examplesDir)

	if len(filesToRestore) == 0 {
		t.Error("convertModulesToLocal should return files to restore")
	}

	// Should have converted main.tf in both modules
	if len(filesToRestore) != 2 {
		t.Errorf("Expected 2 files to restore, got %d", len(filesToRestore))
	}
}

func TestConvertModulesToLocal_WithExceptions(t *testing.T) {
	tmpDir := t.TempDir()
	examplesDir := filepath.Join(tmpDir, "examples")
	if err := os.MkdirAll(examplesDir, 0755); err != nil {
		t.Fatalf("Failed to create examples directory: %v", err)
	}

	moduleNames := []string{"example1", "example2"}
	exceptionList := []string{"example2"}

	for _, modName := range moduleNames {
		modDir := filepath.Join(examplesDir, modName)
		if err := os.Mkdir(modDir, 0755); err != nil {
			t.Fatalf("Failed to create module directory: %v", err)
		}

		tfContent := `
module "test" {
  source  = "cloudnationhq/mymodule/azure"
  version = "~> 1.0"
}
`
		tfFile := filepath.Join(modDir, "main.tf")
		if err := os.WriteFile(tfFile, []byte(tfContent), 0644); err != nil {
			t.Fatalf("Failed to create terraform file: %v", err)
		}
	}

	client := &mockRegistryClient{latestVersion: "1.0.0"}
	converter := NewSourceConverter(client)
	moduleInfo := ModuleInfo{
		Name:      "mymodule",
		Provider:  "azure",
		Namespace: "cloudnationhq",
	}

	ctx := testContext(t)
	mockT := &testing.T{}
	filesToRestore := convertModulesToLocal(ctx, mockT, converter, moduleNames, exceptionList, moduleInfo, examplesDir)

	// Should only convert example1, not example2
	if len(filesToRestore) != 1 {
		t.Errorf("Expected 1 file to restore (excluding exception), got %d", len(filesToRestore))
	}
}
