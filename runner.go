package validor

import (
	"testing"
)

// RunTests executes tests for multiple modules
func RunTests(t *testing.T, modules []*Module, parallel bool, config *Config) {
	runModuleTests(t, modules, parallel, config, nil, "registry")
}
