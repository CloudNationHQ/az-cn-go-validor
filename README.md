# validor [![Go Reference](https://pkg.go.dev/badge/github.com/cloudnationhq/az-cn-go-validor.svg)](https://pkg.go.dev/github.com/cloudnationhq/az-cn-go-validor)

A terraform module testing tool that validates infrastructure configurations by executing real apply and destroy operations.

Ensures your modules deploy successfully, handles cleanup automatically, and provides detailed error reporting for reliable infrastructure testing.

## Why validor?

Terraform modules can fail in production due to untested configurations, provider incompatibilities, or incomplete setups.

Manual testing is time-consuming and error-prone.

Validor helps you:

Test modules in isolated environments before production.

Validate apply/destroy cycles with real provider interactions.

Run multiple modules concurrently for faster CI/CD pipelines.

Test with local sources, exceptions, and custom configurations.

Automate testing across teams and large codebases.

## Installation

`go get github.com/cloudnationhq/az-cn-go-validor`

## Usage

See the [examples/](examples/) directory for sample Terraform modules and test configurations.

Run unit tests:

```
go test ./...
```

## Features

`Module Testing`

Executes full Terraform apply/destroy cycles for real validation.

Supports parallel and sequential execution modes.

Handles local source testing for module development.

Provides detailed error reporting with actionable feedback.

`Flexible Configuration`

Command-line flags for runtime configuration (`-example`, `-exception`, `-local`, `-namespace`, `-examples-path`).

Optional pattern support for programmatic configuration (e.g., `WithExamplesPath`, `WithExample`).

Environment variable support for CI/CD integration.

Exception lists to skip problematic modules.

Configurable namespace for custom registry sources.

Customizable examples directory path for flexible project structures.

`Advanced Terraform Support`

Works with all major cloud providers and custom modules.

Respects Terraform state and resource dependencies.

Handles complex module structures and submodules.

Automatic cleanup of generated files and states.

`Error Reporting & Logging`

Structured error types for better debugging.

Outputstest summaries with failure details.

Integration with Go testing framework for CI/CD.

## Configuration

`Command-Line Flags`

`-example`: Comma-separated list of specific examples to test.

`-exception`: Comma-separated list of examples to exclude.

`-local`: Use local source paths instead of registry.

`-namespace`: Terraform registry namespace (default: "cloudnationhq").

`-skip-destroy`: Skip destroy operations after apply.

`-examples-path`: Path to examples directory (defaults to '../examples').

`Environment Variables`

For CI/CD pipelines, configure via environment variables:

`VALIDOR_EXAMPLE`: Specific examples to test.

`VALIDOR_EXCEPTION`: Examples to exclude.

`VALIDOR_LOCAL`: Enable local source testing (true/false).

`VALIDOR_NAMESPACE`: Registry namespace.

`VALIDOR_SKIP_DESTROY`: Skip destroy (true/false).

### Notes

Local testing requires the module repository to be properly structured.

Namespace configuration allows testing against custom registries.

## Contributors

We welcome contributions from the community! Whether it's reporting a bug, suggesting a new feature, or submitting a pull request, your input is highly valued. <br><br>

<a href="https://github.com/cloudnationhq/az-cn-go-validor/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=cloudnationhq/az-cn-go-validor" />
</a>
