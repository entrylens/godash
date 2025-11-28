# slogkit

`slogkit` package provides a context-aware `slog.Handler` implementation that extends Go's structured logging with dynamic context-based attributes, source tracking, and flexible configuration options.

## Features

- Context-aware attribute extraction from Go context
- JSON and text output formats
- Source file and line number tracking
- Process ID logging
- Static and dynamic attributes
- Full `slog.Handler` interface compliance

## Installation

```bash
go get github.com/entrylens/godash/slogkit
```

## Usage

```go
package main

import (
	"context"
	"log/slog"

	"github.com/entrylens/godash/slogkit"
)

func main() {
	handler := slogkit.NewContextHandler(slogkit.ContextHandlerOptions{
		UseJson:   true,
		Level:     slog.LevelInfo,
		AddSource: true,
		WithPID:   true,
		AppendAttrFromContext: func(ctx context.Context) ([]slog.Attr, error) {
			if requestID, ok := ctx.Value("request_id").(string); ok {
				return []slog.Attr{slog.String("request_id", requestID)}, nil
			}
			return nil, nil
		},
	})

	logger := slog.New(handler)
	slog.SetDefault(logger)

	slog.Info("Application started")

	ctx := context.WithValue(context.Background(), "request_id", "req-12345")
	slog.InfoContext(ctx, "Processing request")
}
```

## ContextHandlerOptions

```go
type ContextHandlerOptions struct {
	UseJson               bool                                  // JSON (true) or text (false) format
	Level                 slog.Level                            // Minimum log level
	AddSource             bool                                  // Enable source file tracking
	SourceKey             string                                // Custom source key (default: "source")
	Writer                io.Writer                             // Output writer (default: os.Stdout)
	WithPID               bool                                  // Include process ID
	PIDKey                string                                // Custom PID key (default: "pid")
	CallerSkip            int                                   // Stack frames to skip (default: 3)
	ExtraAttrs            []slog.Attr                           // Static attributes for all logs
	AppendAttrFromContext func(ctx context.Context) ([]slog.Attr, error) // Extract attributes from context
}
```

## Examples

### Basic JSON Handler

```go
handler := slogkit.NewContextHandler(slogkit.ContextHandlerOptions{
	UseJson: true,
	Level:   slog.LevelInfo,
})

logger := slog.New(handler)
logger.Info("User logged in", "user_id", 123)
```

### With Source Tracking

```go
handler := slogkit.NewContextHandler(slogkit.ContextHandlerOptions{
	UseJson:   true,
	AddSource: true,
	SourceKey: "file_location",
})

logger := slog.New(handler)
logger.Info("Processing request")
```

### With Static Attributes

```go
handler := slogkit.NewContextHandler(slogkit.ContextHandlerOptions{
	UseJson: true,
	ExtraAttrs: []slog.Attr{
		slog.String("service", "api-server"),
		slog.String("version", "1.0.0"),
	},
})

logger := slog.New(handler)
logger.Info("Request received")
```

### Context-Aware Logging

```go
handler := slogkit.NewContextHandler(slogkit.ContextHandlerOptions{
	UseJson: true,
	AppendAttrFromContext: func(ctx context.Context) ([]slog.Attr, error) {
		attrs := []slog.Attr{}
		if requestID, ok := ctx.Value("request_id").(string); ok {
			attrs = append(attrs, slog.String("request_id", requestID))
		}
		if userID, ok := ctx.Value("user_id").(int); ok {
			attrs = append(attrs, slog.Int("user_id", userID))
		}
		return attrs, nil
	},
})

logger := slog.New(handler)
ctx := context.WithValue(context.Background(), "request_id", "req-123")
logger.InfoContext(ctx, "Processing request")
```

## Testing

```bash
go test ./slogkit/...
```
