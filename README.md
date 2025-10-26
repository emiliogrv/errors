<div align="center">

[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.18-blue)](https://go.dev/)
[![License](https://img.shields.io/badge/license-BSD--3--Clause-green)](LICENSE)
</div>

## Overview

This **ready-to-use** library provides **structured error handling** for Go applications, designed as a complete drop-in replacement for the standard `errors` package. It extends standard error functionality with:

- **Structured attributes** for rich error context
- **Error wrapping and joining** with full compatibility with `errors.Is`, `errors.As`, and `errors.Join`
- **Stack trace capture** for debugging
- **JSON serialization** for structured logging
- **Direct integration** with popular logging frameworks (Zap, Zerolog, Logrus, slog)
- **Customizable code generation** for your own logging frameworks

## Performance

This library is designed with **performance as a core principle**. The primary goal is to provide meaningful, structured error messages with **maximum performance** by:

- **Avoiding reflection** wherever possible - type-safe attribute helpers eliminate runtime reflection overhead
- **Zero-allocation string building** for common error formatting paths
- **Efficient marshaling** - custom marshalers for each logging framework avoid generic reflection-based serialization
- **Lazy evaluation** - error details are only formatted when actually logged or serialized
- **Minimal allocations** - careful memory management to reduce GC pressure

## Features

### Core Functionality

The library maintains full compatibility with Go's standard `errors` package while adding powerful structured error handling capabilities:

```go
import errors "github.com/emiliogrv/errors/pkg/full"

err := errors.New("database connection failed").
    WithAttrs(
        errors.String("host", "localhost"),
        errors.Int("port", 5432),
        errors.Duration("timeout", 30*time.Second),
    ).
    WithTags("database", "critical").
    WithStack(debug.Stack())
```

### Partial Imports

**Import only what you need!** To avoid pulling all logging framework dependencies into your `go.mod`, you can import specific logger packages:

```go
// Import only Zap support
import errors "github.com/emiliogrv/errors/pkg/zap"

// Import only slog support
import errors "github.com/emiliogrv/errors/pkg/slog"

// Import only Logrus support
import errors "github.com/emiliogrv/errors/pkg/logrus"

// Import only Zerolog support
import errors "github.com/emiliogrv/errors/pkg/zerolog"

// Import only core functionality (no logger integrations)
import errors "github.com/emiliogrv/errors/pkg/core"
```

Each partial import includes the core templates plus the specific logger integration, keeping your dependencies minimal.

### Custom Template Generation

The library supports **custom template generation**, allowing you to:

1. **Generate your own logger integrations** by providing custom templates
2. **Override default templates** to customize behavior
3. **Add new functionality** by simply using the `-input-dir` flag with your template directory

Example:

```bash
go run github.com/emiliogrv/errors/cmd/errors_generator \
    -input-dir ./my-templates \
    -output-dir ./generated \
    -formats mylogger,zap
```

**Core templates are always generated** regardless of which formats you specify, ensuring base functionality is always available.

## Installation

```bash
go get github.com/emiliogrv/errors
```

## Usage

### Basic Error Creation

```go
import errors "github.com/emiliogrv/errors/pkg/full"

// Simple error
err := errors.New("something went wrong")

// Error with attributes
err := errors.New("failed to process request").
    WithAttrs(
        errors.String("user_id", "12345"),
        errors.Int("retry_count", 3),
        errors.Bool("recoverable", true),
    )
```

### Error Wrapping

```go
// Wrap errors with context
if err := doSomething(); err != nil {
    return errors.New("operation failed").
        WithErrors(err).
        WithAttrs(errors.String("operation", "doSomething"))
}

// Join multiple errors
err := errors.Join(
    errors.New("validation failed"),
    errors.New("missing required field: email"),
    errors.New("missing required field: name"),
)
```

### Logging Integration

#### Full Import Example

```go
import (
    "go.uber.org/zap"
    errors "github.com/emiliogrv/errors/pkg/full"
)

func main() {
    err := errors.New("operation failed").
        WithAttrs(errors.Int("code", 500))

    // Zap
    logger, _ := zap.NewProduction()
    logger.Error("error occurred", zap.Any("err", err))

    // JSON marshaling
    jsonData, _ := json.Marshal(err)

    // Logrus
    logrus.WithFields(logrus.Fields{
            "err": err.(*errors.StructuredError).MarshalLogrusFields(),
        }).
        Error("error occurred")
}
```

#### Partial Import Example

```go
import (
    "go.uber.org/zap"
    errors "github.com/emiliogrv/errors/pkg/zap" // Only Zap dependency
)

func main() {
    err := errors.New("operation failed").
        WithAttrs(errors.Int("code", 500))

    logger, _ := zap.NewProduction()
    logger.Error("error occurred", zap.Any("err", err))
}
```

### Custom Template Example

Create your custom template (e.g., `tmpl/mylogger.tmpl`):

```go
// Code generated by errors_generator; DO NOT EDIT.
// Generated at {{ .Date }}
// Version {{ .Version }}

package {{ .PackageName }}

// MarshalMyLogger is the implementation for MyLogger.
func (receiver *StructuredError) MarshalMyLogger() error {
    // Your custom marshaling logic here
    return nil
}
```

Generate code:

```bash
go run github.com/emiliogrv/errors/cmd/errors_generator \
    -input-dir ./tmpl \
    -output-dir ./generated \
    -formats mylogger
```

## Templates

### Core Templates

Core templates are **always generated** and provide the fundamental error handling functionality:

| Template    | Description                                                           |
|-------------|-----------------------------------------------------------------------|
| `attr.go`   | Type-safe attribute helpers (String, Int, Bool, Time, Duration, etc.) |
| `common.go` | Common utilities and depth control for marshaling                     |
| `error.go`  | Core `StructuredError` type and basic methods                         |
| `join.go`   | `Join` and `JoinIf` functions for combining errors                    |
| `json.go`   | JSON marshaling/unmarshaling support                                  |
| `map.go`    | Map representation for generic structured output                      |
| `string.go` | String formatting and `Error()` method implementation                 |
| `wrap.go`   | `Unwrap`, `Is`, and `As` methods for error wrapping                   |

### Logger Templates

Additional templates for specific logging framework integrations:

| Package       | Templates                            | Dependencies                 |
|---------------|--------------------------------------|------------------------------|
| `pkg/full`    | Core + Zap + Zerolog + Logrus + slog | All logger dependencies      |
| `pkg/zap`     | Core + Zap                           | `go.uber.org/zap`            |
| `pkg/zerolog` | Core + Zerolog                       | `github.com/rs/zerolog`      |
| `pkg/logrus`  | Core + Logrus                        | `github.com/sirupsen/logrus` |
| `pkg/slog`    | Core + slog                          | Standard library only        |
| `pkg/core`    | Core only                            | No external dependencies     |

### Template Overriding

You can **override any default template** by providing a template with the same name in your input directory:

```bash
# Override the default zap.tmpl with your custom version
go run github.com/emiliogrv/errors/cmd/errors_generator \
    -input-dir ./my-templates \
    -output-dir ./generated \
    -formats zap
```

If `my-templates/zap.tmpl` exists, it will replace the built-in Zap template.

## Generator Usage

The code generator supports various options:

```bash
go run github.com/emiliogrv/errors/cmd/errors_generator [options]

Options:
  -formats string
        Comma-separated list of formats to generate, or 'all' to generate all formats (default: core)
  -help
        Show this help message
  -input-dir string
        Path to user templates directory (optional)
  -output-dir string
        Output directory for generated files
  -export-dir string
        Export default templates to the specified directory and exit
  -package string
        Package name for generated code (default: errors) (default "errors")
  -test-gen string
        Test generation level: none, flex, strict (default: none) (default "none")
  -with-gen-header
        Include generated message in generated code (default: true) (default true)
```

### Examples

```bash
# Export default templates for customization
go run github.com/emiliogrv/errors/cmd/errors_generator -export-dir ./my-templates

# Generate core templates only
go run github.com/emiliogrv/errors/cmd/errors_generator -output-dir ./pkg/core

# Generate with Zap support
go run github.com/emiliogrv/errors/cmd/errors_generator \
    -output-dir ./pkg/zap \
    -formats zap

# Generate all available formats
go run github.com/emiliogrv/errors/cmd/errors_generator \
    -output-dir ./pkg/full \
    -formats all

# Generate with custom templates
go run github.com/emiliogrv/errors/cmd/errors_generator \
    -input-dir ./my-templates \
    -output-dir ./generated \
    -formats mylogger,zap

# Generate with tests
go run github.com/emiliogrv/errors/cmd/errors_generator \
    -output-dir ./pkg/full \
    -formats all \
    -test-gen strict
```

## API Reference

### Core Types

#### `StructuredError`

```go
type StructuredError struct {
    Message string // Primary error message
    Attrs   []Attr // Structured attributes
    Errors  []error  // Wrapped errors
    Tags    []string // Categorical labels
    Stack   []byte // Stack trace (optional)
}
```

#### `Attr`

```go
type Attr struct {
    Key   string
    Value any
    Type  Type
}
```

### Core Functions

- `New(message string) *StructuredError` - Create a new structured error
- `Join(errs ...error) error` - Join multiple errors (nil-safe)
- `JoinIf(errs ...error) error` - Join errors only if first is non-nil
- `Is(err, target error) bool` - Check error equality (alias to `errors.Is`)
- `As(err error, target any) bool` - Type assertion (alias to `errors.As`)
- `Unwrap(err error) error` - Unwrap single error (alias to `errors.Unwrap`)

### Attribute Helpers

Type-safe helpers for common types:

- `String(key, value string) Attr`
- `Int(key string, value int) Attr`
- `Int64(key string, value int64) Attr`
- `Uint64(key string, value uint64) Attr`
- `Float64(key string, value float64) Attr`
- `Bool(key string, value bool) Attr`
- `Time(key string, value time.Time) Attr`
- `Duration(key string, value time.Duration) Attr`
- `Any(key string, value any) Attr`
- `Object(key string, attrs ...Attr) Attr`

Each helper also has a plural version (e.g., `Ints`, `Strings`, `Bools`) for slices.

### Methods

#### `*StructuredError` Methods

- `WithAttrs(attrs ...Attr) *StructuredError` - Add attributes
- `WithErrors(errors ...error) *StructuredError` - Set wrapped errors
- `WithTags(tags ...string) *StructuredError` - Add tags
- `WithStack(stack []byte) *StructuredError` - Set stack trace
- `PrependErrors(errors ...error) *StructuredError` - Add errors at the beginning
- `AppendErrors(errors ...error) *StructuredError` - Add errors at the end
- `Error() string` - Implement error interface
- `Unwrap() []error` - Implement multi-unwrapper interface
- `MarshalJSON() ([]byte, error)` - JSON marshaling
- `UnmarshalJSON(data []byte) error` - JSON unmarshaling

### Configuration

```go
// Set maximum depth for error marshaling (default: 100)
errors.SetMaxDepthMarshal(depth int)

// Get current maximum depth
errors.MaxDepthMarshal() int
```

## API Stability

⚠️ **Pre-1.0 Notice**: The library API is essentially complete but may undergo changes before version 1.0. While we strive to maintain backward compatibility, breaking changes may occur in minor versions (0.x.y) until the 1.0 release.

Once version 1.0 is released, the API will follow semantic versioning strictly.

## Examples

Complete examples are available in the [examples/](examples/) directory:

- [full import](examples/full_import/) - Full import with all logger integrations
- [partial import](examples/partial_import/) - Partial import (Zap only)
- [custom tmpl](examples/custom_tmpl/) - Custom template generation

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.

## License

This project is licensed under the BSD 3-Clause License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

This library is designed to be a drop-in replacement for Go's standard `errors` package while providing enhanced functionality for modern application development.
