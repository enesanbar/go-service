# Go Service Development Guidelines

This document provides guidelines and information for developers working on the Go Service project.

## Build/Configuration Instructions

### Project Structure

The Go Service project is organized as a collection of Go modules, each with its own specific functionality:

- **core**: Contains the core functionality of the service, including configuration, logging, error handling, and dependency injection.
- **cache**: Provides caching functionality (currently in-memory, with Redis and Memcached support planned).
- **messaging**: Implements messaging functionality, including RabbitMQ support.
- **persistence**: Provides database connectivity (MongoDB, MySQL).
- **protocol**: Implements HTTP (REST) and gRPC server functionality.
- **cron**: Provides scheduling functionality for cron jobs.

### Configuration

The project uses [Viper](https://github.com/spf13/viper) for configuration, supporting both configuration files and environment variables.

#### Configuration Sources

- **File (default)**: Uses YAML files in the `config` directory.
- **Consul**: Retrieves configuration from Consul.

Set the configuration source using the `CONFIG_SOURCE` environment variable:

```shell
CONFIG_SOURCE=file go run *.go  # Default
CONFIG_SOURCE=consul go run *.go
```

#### Environment Configuration

Supported environments:
- dev (default)
- test
- staging
- prod

Set the environment using the `DEPLOY_TYPE` environment variable:

```shell
DEPLOY_TYPE=prod go run *.go
```

#### Configuration Files

For file-based configuration, create the following files:
- `${PROJECT_DIR}/config/base.yaml`: Base configuration
- `${PROJECT_DIR}/config/[environment].yaml`: Environment-specific configuration

For Consul-based configuration, create:
- `go-config/my-service.yaml`: Base configuration
- `go-config/my-service_[environment].yaml`: Environment-specific configuration

### Building a New Service

1. Create a new Go module:
   ```shell
   mkdir my-service
   cd my-service
   go mod init github.com/my-org/my-service
   ```

2. Create a main.go file with the service configuration (see README.md for examples).

3. Install dependencies:
   ```shell
   go mod tidy
   ```

4. Run the service:
   ```shell
   DEPLOY_TYPE=dev go run main.go
   ```

## Testing Information

### Running Tests

Tests follow the standard Go testing convention with files named `*_test.go`. To run tests for a specific package:

```shell
cd path/to/package
go test -v
```

To run all tests in the project with coverage:

```shell
go test ./... -v -cover
```

### Writing Tests

1. Create a test file with the same package name as the file being tested, with a `_test.go` suffix.
2. Import the `testing` package and any other required packages.
3. Write test functions with the prefix `Test` followed by the name of the function or method being tested.
4. Use the `t.Error` or `t.Errorf` methods to report test failures.

### Example Test

Here's an example test for the `Error` type in the `core/errors` package:

```go
package errors

import (
	"errors"
	"testing"
)

func TestNewError(t *testing.T) {
	// Test case 1: Create a new error with all fields
	err := NewError("ERR001", "Something went wrong", "TestOperation", errors.New("underlying error"))
	
	if err.Code != "ERR001" {
		t.Errorf("Expected Code to be 'ERR001', got '%s'", err.Code)
	}
	
	// Additional assertions...
}
```

## Code Style and Development Practices

### Linting

The project uses [golangci-lint](https://golangci-lint.run/) for code quality checks. The configuration is in `.golangci.yml`.

To run the linter:

```shell
golangci-lint run
```

Enabled linters include:
- govet: Reports suspicious constructs
- errcheck: Checks for unchecked errors
- staticcheck: Finds bugs and performance issues
- unused: Checks for unused variables, constants, etc.
- ineffassign: Detects ineffectual assignments
- gocritic: Detects various code issues
- gosec: Security-related issues
- prealloc: Suggests preallocating slices
- whitespace: Detects unnecessary whitespace
- misspell: Detects spelling errors
- revive: Extensible, configurable, and faster linter

### Dependency Injection

The project uses [uber-go/fx](https://github.com/uber-go/fx) for dependency injection and application lifecycle management. When adding new components:

1. Create a module that exports the component.
2. Register the module in the service's dependency injection container.

### Error Handling

Use the custom error types defined in `core/errors` for consistent error handling:

```go
import "github.com/enesanbar/go-service/core/errors"

func SomeFunction() error {
	// Create a new error
	return errors.NewError("ERR001", "Something went wrong", "SomeFunction", nil)
}
```

### Logging

Use the structured logging provided by `core/log` for consistent logging:

```go
import "github.com/enesanbar/go-service/core/log"

func SomeFunction() {
	logger := log.GetLogger()
	logger.Info("Something happened", "key", "value")
}
```

### Health Checks

Implement health checks for all components that connect to external services:

```go
import "github.com/enesanbar/go-service/core/healthchecker"

type MyHealthChecker struct {}

func (c *MyHealthChecker) Check() error {
	// Check if the component is healthy
	return nil
}

func RegisterHealthChecker() healthchecker.HealthChecker {
	return &MyHealthChecker{}
}
```