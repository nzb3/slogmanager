# slogmanager

---

**Concurrent-Safe Manager for Slog with Multiple Writers**

A Go library providing simple writers setup for Slog
## Features

- **Multiple Writers**: Configure different outputs (console, files, buffers)
- **Concurrent-Safe**: Thread-safe writer operations
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
	"log/slog"
)

func main() {
	manager := slogmanager.New()
	
	// Add console writer
	consoleWriter := slogmanager.NewWriter(&bytes.Buffer{}, slogmanager.WithTextFormat())
	manager.AddWriter("console", consoleWriter)
	
	slog.Info("System initialized", "status", "ready")
}
```


### JSON Formatting

```go
jsonWriter := slogmanager.NewWriter(
    &bytes.Buffer{},
    slogmanager.WithJSONFormat(),
)
manager.AddWriter("json-output", jsonWriter)

slog.Info("System initialized", "status", "ready")
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

slog.Info("System initialized", "status", "ready")
// Output:
// 2025/02/22 23:49:00 INFO request completed method=GET duration=150ms
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

