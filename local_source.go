package validor

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"
)

func convertModulesToLocal(ctx context.Context, t *testing.T, converter SourceConverter, moduleNames []string, exceptionList []string, moduleInfo ModuleInfo, examplesPath string) []FileRestore {
	var allFilesToRestore []FileRestore

	for _, moduleName := range moduleNames {
		if slices.Contains(exceptionList, moduleName) {
			continue
		}

		modulePath := filepath.Join(examplesPath, moduleName)
		filesToRestore, err := converter.ConvertToLocal(ctx, modulePath, moduleInfo)
		if err != nil {
			t.Logf("Warning: Failed to convert module %s to local source: %v", moduleName, err)
			continue
		}
		allFilesToRestore = append(allFilesToRestore, filesToRestore...)
	}

	return allFilesToRestore
}

func createLocalSetupFunc(config *Config) TestSetupFunc {
	return func(ctx context.Context, t *testing.T, modules []*Module) error {
		moduleInfo := extractModuleInfoFromRepo()
		if moduleInfo.Name == "" || moduleInfo.Provider == "" {
			return fmt.Errorf("could not determine module name and provider from repository")
		}
		moduleInfo.Namespace = config.Namespace

		converter := NewSourceConverter(NewRegistryClient())
		moduleNames := extractModuleNames(modules)
		allFilesToRestore := convertModulesToLocal(ctx, t, converter, moduleNames, config.ExceptionList, moduleInfo, getExamplesPath(config))

		t.Cleanup(func() {
			cleanupCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			if err := converter.RevertToRegistry(cleanupCtx, allFilesToRestore); err != nil {
				t.Logf("Warning: Failed to revert files to registry source: %v", err)
			}
		})
		return nil
	}
}

func extractModuleInfoFromRepo() ModuleInfo {
	wd, err := os.Getwd()
	if err != nil {
		return ModuleInfo{}
	}

	if filepath.Base(wd) == "tests" {
		wd = filepath.Dir(wd)
	}

	if repoName := getRepoNameFromGit(wd); repoName != "" {
		if info, ok := parseModuleName(repoName); ok {
			return info
		}
	}

	repoName := filepath.Base(wd)
	if info, ok := parseModuleName(repoName); ok {
		return info
	}
	return ModuleInfo{}
}

func parseModuleName(repoName string) (ModuleInfo, bool) {
	const prefix = "terraform-"
	if !strings.HasPrefix(repoName, prefix) {
		return ModuleInfo{}, false
	}

	parts := strings.SplitN(repoName[len(prefix):], "-", 2)
	if len(parts) != 2 {
		return ModuleInfo{}, false
	}

	return ModuleInfo{
		Provider: parts[0],
		Name:     parts[1],
	}, true
}

func getRepoNameFromGit(dir string) string {
	output, err := gitRemoteURL(dir)
	if err != nil {
		return ""
	}

	url := strings.TrimSpace(string(output))
	parts := strings.Split(url, "/")
	if len(parts) > 0 {
		repoName := parts[len(parts)-1]
		return strings.TrimSuffix(repoName, ".git")
	}
	return ""
}

var gitRemoteURL = func(dir string) ([]byte, error) {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	cmd.Dir = dir
	return cmd.Output()
}
