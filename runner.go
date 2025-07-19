package validor

import (
	"context"
	"testing"
	"time"
)

// RunTests executes tests for multiple modules
func RunTests(t *testing.T, modules []*Module, parallel bool, config *Config) {
	results := NewTestResults()
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 30*time.Minute)
	defer cancel()

	for _, module := range modules {
		t.Run(module.Name, func(t *testing.T) {
			if parallel {
				t.Parallel()
			}

			if !config.SkipDestroy {
				defer func() {
					if err := module.Destroy(ctx, t); err != nil && !module.ApplyFailed {
						t.Logf("Warning: Cleanup for module %s failed: %v", module.Name, err)
					}
				}()
			}

			if err := module.Apply(ctx, t); err != nil {
				t.Fail()
			} else {
				t.Logf("✓ Module %s applied successfully", module.Name)
			}

			results.AddModule(module)
		})
	}

	t.Cleanup(func() {
		modules, _ := results.GetResults()
		PrintModuleSummary(t, modules)
	})
}
