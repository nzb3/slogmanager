# slogmanager

---

**Concurrent-Safe Slog Manager with Multiple Writers**

A Go library providing enhanced `slog` logging management with support for multiple output formats and concurrent access.

## Features

- **Multiple Writers**: Configure different outputs (console, files, buffers)
- **Concurrent-Safe**: Thread-safe operations using RWMutex
- **Format Support**: JSON and text logging formats
- **Dynamic Configuration**: Add/remove writers at runtime
- **Handler Customization**: Create custom handlers for specialized logging


## Installation

```bash
go get github.com/nzb3/slogmanager
```


## Usage

### Basic Example

```go
package main

import (
	"bytes"
	"github.com/nzb3/slogmanager"
)

func main() {
	manager := slogmanager.New()
	
	// Add console writer
	consoleWriter := slogmanager.NewWriter(&bytes.Buffer{}, slogmanager.WithTextFormat())
	manager.AddWriter("console", consoleWriter)
	
	logger := manager.Logger()
	logger.Info("System initialized", "status", "ready")
}
```


### JSON Formatting

```go
jsonWriter := slogmanager.NewWriter(
    &bytes.Buffer{},
    slogmanager.WithJSONFormat(),
)
manager.AddWriter("json-output", jsonWriter)

// Output:
// {"time":"2025-02-22T23:49:00Z","level":"INFO","msg":"Data processed","items":42}
```


### Text Formatting

```go
textWriter := slogmanager.NewWriter(
    &bytes.Buffer{},
    slogmanager.WithTextFormat(),
)
manager.AddWriter("text-output", textWriter)

// Output:
// 2025/02/22 23:49:00 INFO request completed method=GET duration=150ms
```


##  Logging

```go
func main() {
	manager := slogmanager.New()
	manager.AddWriter("concurrent", slogmanager.NewWriter(&bytes.Buffer{}))
	
	slog.Debug("Worker started", "id", 10)
}
```


## Configuration Options

| Option | Description |
| :-- | :-- |
| `WithJSONFormat()` | JSON-structured logs |
| `WithTextFormat()` | Human-readable text format |
| `WithLevel(level slog.Level)` | Set minimum log level |
| `WithHandlerOpts(opts *slog.HandlerOptions)` | Custom handler options |

## License

MIT License - See [LICENSE](LICENSE) for details.

<div style="text-align: center">⁂</div>

[^1]: https://github.com/nzb3/slogmanager

