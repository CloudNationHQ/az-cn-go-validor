package validor

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type mockTB struct {
	logs  []string
	fatal bool
}

func (m *mockTB) Helper() {}

func (m *mockTB) Log(args ...any) {
	m.logs = append(m.logs, strings.TrimSpace(fmt.Sprint(args...)))
}

func (m *mockTB) Logf(format string, args ...any) {
	m.logs = append(m.logs, fmt.Sprintf(format, args...))
}

func (m *mockTB) Fatal(args ...any) {
	m.fatal = true
	m.logs = append(m.logs, strings.TrimSpace(fmt.Sprint(args...)))
}

type mockRegistryClient struct {
	latestVersion string
	err           error
}

func (m *mockRegistryClient) GetLatestVersion(ctx context.Context, namespace, name, provider string) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.latestVersion, nil
}

type roundTripperFunc func(req *http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func testContext(t *testing.T) context.Context {
	t.Helper()
	return context.Background()
}

func setupMockExamplesDir(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()

	for _, example := range []string{"example1", "example2", "example3"} {
		exampleDir := filepath.Join(tmpDir, example)
		if err := os.MkdirAll(exampleDir, 0o755); err != nil {
			t.Fatalf("failed to create example dir: %v", err)
		}
		mainTf := filepath.Join(exampleDir, "main.tf")
		if err := os.WriteFile(mainTf, []byte("# mock terraform file"), 0o644); err != nil {
			t.Fatalf("failed to create main.tf: %v", err)
		}
	}

	return tmpDir
}

func createMockModules(names []string, basePath string) []*Module {
	modules := make([]*Module, len(names))
	for i, name := range names {
		modules[i] = NewModule(name, filepath.Join(basePath, name))
		modules[i].applyHook = func(ctx context.Context, tb *testing.T, m *Module) error {
			return nil
		}
		modules[i].destroyHook = func(ctx context.Context, tb *testing.T, m *Module) error {
			return nil
		}
	}
	return modules
}
