package validor

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

// DefaultSourceConverter implements SourceConverter
type DefaultSourceConverter struct {
	registryClient RegistryClient
}

// NewSourceConverter creates a new source converter
func NewSourceConverter(client RegistryClient) SourceConverter {
	return &DefaultSourceConverter{
		registryClient: client,
	}
}

// ConvertToLocal converts module blocks in Terraform files to use local source
func (c *DefaultSourceConverter) ConvertToLocal(ctx context.Context, modulePath string, moduleInfo ModuleInfo) ([]FileRestore, error) {
	var filesToRestore []FileRestore

	files, err := filepath.Glob(filepath.Join(modulePath, "*.tf"))
	if err != nil {
		return nil, fmt.Errorf("failed to find terraform files: %w", err)
	}

	modulePattern := fmt.Sprintf(`(?m)^(\s*module\s+"[^"]*"\s*\{[^}]*source\s*=\s*)"%s/%s/%s"([^}]*version\s*=\s*"[^"]*")?([^}]*\})`,
		regexp.QuoteMeta(moduleInfo.Namespace), regexp.QuoteMeta(moduleInfo.Name), regexp.QuoteMeta(moduleInfo.Provider))
	re := regexp.MustCompile(modulePattern)

	submodulePattern := fmt.Sprintf(`(?m)^(\s*module\s+"[^"]*"\s*\{[^}]*source\s*=\s*)"%s/%s/%s//modules/([^"]*)"([^}]*version\s*=\s*"[^"]*")?([^}]*\})`,
		regexp.QuoteMeta(moduleInfo.Namespace), regexp.QuoteMeta(moduleInfo.Name), regexp.QuoteMeta(moduleInfo.Provider))
	subRe := regexp.MustCompile(submodulePattern)

	for _, file := range files {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		content, err := os.ReadFile(file)
		if err != nil {
			continue
		}

		originalContent := string(content)
		newContent := c.processContent(originalContent, re)
		newContent = c.processSubmoduleContent(newContent, subRe)

		if newContent != originalContent {
			if err := os.WriteFile(file, []byte(newContent), 0644); err != nil {
				return filesToRestore, fmt.Errorf("failed to write file %s: %w", file, err)
			}

			filesToRestore = append(filesToRestore, FileRestore{
				Path:            file,
				OriginalContent: originalContent,
				ModuleName:      moduleInfo.Name,
				Provider:        moduleInfo.Provider,
				Namespace:       moduleInfo.Namespace,
			})
		}
	}

	return filesToRestore, nil
}

// processContent replaces module source and removes version
func (c *DefaultSourceConverter) processContent(content string, re *regexp.Regexp) string {
	return re.ReplaceAllStringFunc(content, func(match string) string {
		parts := re.FindStringSubmatch(match)
		if len(parts) < 4 {
			return match
		}

		moduleStart := parts[1]
		moduleEnd := parts[3]

		// Remove version line if present
		versionRegex := regexp.MustCompile(`(?m)^\s*version\s*=\s*"[^"]*"\s*\n?`)
		moduleEnd = versionRegex.ReplaceAllString(moduleEnd, "")

		return fmt.Sprintf(`%s"../../"%s`, moduleStart, moduleEnd)
	})
}

// processSubmoduleContent replaces submodule source and removes version
func (c *DefaultSourceConverter) processSubmoduleContent(content string, re *regexp.Regexp) string {
	return re.ReplaceAllStringFunc(content, func(match string) string {
		parts := re.FindStringSubmatch(match)
		if len(parts) < 6 {
			return match
		}

		moduleStart := parts[1]
		submoduleName := parts[4]
		moduleEnd := parts[6]

		// Remove version line if present
		versionRegex := regexp.MustCompile(`(?m)^\s*version\s*=\s*"[^"]*"\s*\n?`)
		moduleEnd = versionRegex.ReplaceAllString(moduleEnd, "")

		return fmt.Sprintf(`%s"../../modules/%s"%s`, moduleStart, submoduleName, moduleEnd)
	})
}

// RevertToRegistry reverts files back to registry source with latest version
func (c *DefaultSourceConverter) RevertToRegistry(ctx context.Context, filesToRestore []FileRestore) error {
	for _, restore := range filesToRestore {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Get latest version from registry
		latestVersion, err := c.registryClient.GetLatestVersion(ctx, restore.Namespace, restore.ModuleName, restore.Provider)
		if err != nil {
			// If we can't get the latest version, fall back to original content
			if writeErr := os.WriteFile(restore.Path, []byte(restore.OriginalContent), 0644); writeErr != nil {
				return fmt.Errorf("failed to restore file %s: %w", restore.Path, writeErr)
			}
			continue
		}

		// Create updated content with latest version
		updatedContent := c.updateVersionInContent(restore.OriginalContent, latestVersion)

		if err := os.WriteFile(restore.Path, []byte(updatedContent), 0644); err != nil {
			return fmt.Errorf("failed to write updated file %s: %w", restore.Path, err)
		}
	}
	return nil
}

// updateVersionInContent updates the version in the content if it exists
func (c *DefaultSourceConverter) updateVersionInContent(content, latestVersion string) string {
	versionRegex := regexp.MustCompile(`(version\s*=\s*")[^"]*(")`)
	if versionRegex.MatchString(content) {
		return versionRegex.ReplaceAllString(content, fmt.Sprintf("${1}~> %s${2}", latestVersion))
	}
	return content
}
