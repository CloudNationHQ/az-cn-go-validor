# CLAUDE.md - az-cn-go-validor

Guidelines for Claude Code to work effectively on this Terraform module testing framework.

## Project Overview

**validor** is a Go testing framework for validating Terraform modules by executing real apply/destroy cycles. Single-package library with comprehensive test coverage (83%+).

- **Go Version**: 1.24.2
- **Test Framework**: Standard library + gruntwork-io/terratest
- **Architecture**: Flat single-package design (no internal/ or cmd/)
- **Coding Style**: Idiomatic Go with functional options pattern

## Key Patterns & Conventions

### Configuration
- Use **functional options pattern** exclusively: `NewConfig(WithSkipDestroy(true), WithLocal(true))`
- For CLI apps: `NewConfigFromFlags()` (registers flags and parses automatically)
- No global state; each config instance is independent

### Error Handling
- Wrap errors with context: `fmt.Errorf("operation %s: %w", name, err)`
- Store errors as `[]error` not `[]string` to preserve error chains
- Custom `ModuleError` type used throughout for structured errors

### Interfaces & Decoupling
- All interfaces in `interfaces.go` (contracts only, no impl)
- Small, focused: `Logger`, `ModuleRunner`, `SourceConverter`, `RegistryClient`
- No coupling to `testing.T`; use `Logger` interface instead

### Concurrency & Context
- **Always** context as first parameter: `func(ctx context.Context, ...)`
- Context **never** embedded in structs
- Use `context.WithTimeout()` for bounded operations
- TestResults is thread-safe (sync.RWMutex)

### Performance
- Regexes cached at package level: `var versionRegex = regexp.MustCompile(...)`
- HTTP client has 10-second timeout
- Cleanup operations timeout after 30 seconds

### Testing
- Table-driven tests throughout
- Subtests via `t.Run()` for hierarchy
- Mock interfaces (MockRegistryClient, mockTB) for dependencies
- Helper functions marked with `t.Helper()`
- Store test errors as `error`: `append(module.Errors, fmt.Errorf("msg"))`

## File Organization

```
interfaces.go      → Interface definitions only
types.go           → Data types, ModuleError, TestResults
module.go          → Module lifecycle (Apply, Destroy, Cleanup)
converter.go       → HCL source conversion (registry ↔ local)
registry.go        → Terraform Registry API client
validor.go         → Main test framework & options
runner.go          → Test runner wrapper
utils.go           → Helpers (BoolToStr, redError)
*_test.go          → Test files
```

## Common Tasks

### Adding a Configuration Option
```go
// In validor.go
func WithMyOption(value string) Option {
    return func(c *Config) { c.MyOption = value }
}

// Add to Config struct
type Config struct {
    // ...
    MyOption string
}
```

### Storing Errors from Operations
```go
// Always store error objects, not strings
wrappedErr := &ModuleError{
    ModuleName: m.Name,
    Operation: "my operation",
    Err: err,
}
m.Errors = append(m.Errors, wrappedErr)  // Not wrappedErr.Error()
```

### Adding Context Cancellation
```go
for _, item := range items {
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
    }
    // process item
}
```

### Writing Tests
```go
func TestMyFeature(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {"case 1", "in1", "out1"},
        {"case 2", "in2", "out2"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := MyFunction(tt.input)
            if got != tt.expected {
                t.Errorf("got %v, want %v", got, tt.expected)
            }
        })
    }
}
```

## Standards & Best Practices

✅ **Always**
- Use pointer receivers on types with state
- Return errors, don't panic
- Wrap errors with `fmt.Errorf("...: %w", err)`
- Pass context as first parameter
- Use `make([]*Module, 0)` for slices
- Mark test helpers with `t.Helper()`
- Store errors as `error` interface, not strings

❌ **Never**
- Embed context in structs
- Use global mutable state
- Ignore error returns
- Mix value and pointer receivers
- Use naked returns
- Panic for control flow
- Build regexes in loops (cache at package level)

## Testing

- **Coverage**: Maintain 83%+ statement coverage
- **Race Detector**: Must pass `go test -race ./...`
- **Table-Driven Tests**: Use for comprehensive coverage
- **Mocking**: Implement interfaces (no external mock libraries)

## Before Committing

```bash
go fmt ./...
go vet ./...
go test -race -cover ./...
go tool cover -func=coverage.out
```
