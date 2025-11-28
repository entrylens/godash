# Godash

[![Go Report Card](https://goreportcard.com/badge/github.com/entrylens/godash)](https://goreportcard.com/report/github.com/entrylens/godash)
[![Go Version](https://img.shields.io/github/go-mod/go-version/entrylens/godash)](https://go.dev/)
[![License](https://img.shields.io/github/license/entrylens/godash)](LICENSE)

A comprehensive Go utility library providing essential tools for modern Go applications. Godash offers type-safe, performant utilities for common programming tasks with a focus on functional programming patterns and developer experience.

## 🛠️ Installation

```bash
go get github.com/entrylens/godash@latest
```

## 🚀 Features

- **Type Safe**: All utilities use Go generics for compile-time type safety
- **Performance Focused**: Optimized for Go's performance characteristics
- **Comprehensive Testing**: Extensive test coverage for all packages
- **Modern Go**: Built for Go 1.23+ with latest language features
- **Production Ready**: Used in production environments with robust error handling

## 📦 Packages

### [sliceskit](./sliceskit/) - Functional Slice Utilities

Comprehensive slice manipulation utilities inspired by functional programming patterns.

```go
import "github.com/entrylens/godash/sliceskit"

numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

// Filter even numbers and map to squares
squares := sliceskit.Map(
    sliceskit.Filter(numbers, func(n int) bool { return n%2 == 0 }),
    func(n int) int { return n * n },
)
// Result: [4, 16, 36, 64, 100]

// Find first number greater than 5
found, exists := sliceskit.Find(numbers, func(n int) bool { return n > 5 })
// Result: 6, true

// Check if any number is even
hasEven := sliceskit.Any(numbers, func(n int) bool { return n%2 == 0 })
// Result: true
```

**Key Functions:**

- `Map`, `Filter`, `Reduce` - Core functional operations
- `Any`, `Every` - Predicate checking
- `Find`, `FindPtr` - Element searching
- `Chunk` - Slice partitioning
- Error handling variants for robust applications

[📖 Full Documentation](./sliceskit/README.md)

### [jsonkit](./jsonkit/) - JSON & Protobuf Utilities

Streamlined JSON and Protocol Buffer handling for HTTP applications.

```go
import "github.com/entrylens/godash/jsonkit"

// HTTP Request/Response handling
type User struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

// Bind JSON request body to struct
var user User
err := jsonkit.BindRequestBody(r, &user)

// Send JSON response
err = jsonkit.JSONResponse(w, http.StatusOK, user)

// Protobuf support
message := &pb.User{Name: "Alice", Age: 30}
err = jsonkit.ProtoJSONResponse(w, http.StatusOK, message)
```

**Key Features:**

- HTTP request/response JSON binding
- Protocol Buffer JSON serialization
- Strict JSON validation (no unknown fields)
- Error handling for malformed JSON

[📖 Full Documentation](./jsonkit/README.md)

### [slogkit](./slogkit/) - Structured Logging Utilities

Context-aware structured logging extensions for Go's `log/slog` package.

```go
import (
	"context"
	"log/slog"
	"github.com/entrylens/godash/slogkit"
)

handler := slogkit.NewContextHandler(slogkit.ContextHandlerOptions{
	UseJson:   true,
	Level:     slog.LevelInfo,
	AddSource: true,
	AppendAttrFromContext: func(ctx context.Context) ([]slog.Attr, error) {
		if requestID, ok := ctx.Value("request_id").(string); ok {
			return []slog.Attr{slog.String("request_id", requestID)}, nil
		}
		return nil, nil
	},
})

logger := slog.New(handler)
ctx := context.WithValue(context.Background(), "request_id", "req-123")
logger.InfoContext(ctx, "Processing request")
```

**Key Features:**

- Context-aware attribute extraction
- JSON and text output formats
- Source file tracking
- Process ID logging

[📖 Full Documentation](./slogkit/README.md)

## 📋 Requirements

- Go 1.23 or higher
- For protobuf features: `google.golang.org/protobuf`

## 🧪 Testing

Run the full test suite:

```bash
make test
```

Or test specific packages:

```bash
go test ./sliceskit/...
go test ./jsonkit/...
go test ./slogkit/...
```

## 🔧 Development

### Prerequisites

- Go 1.23+
- Make (for build automation)
- golangci-lint (for code quality)

### Available Commands

```bash
make check      # Run all checks (format, lint, test)
make fmt        # Format code
make lint       # Run linter
make test       # Run tests
make tidy       # Tidy dependencies
make regen_proto # Regenerate protobuf code
```

### Code Quality

The project uses:

- **golangci-lint** for static analysis
- **testify** for testing utilities
- **Go modules** for dependency management
- **Git hooks** (via lefthook) for pre-commit checks

## 📈 Performance

All utilities are designed with performance in mind:

- Zero allocations where possible
- Efficient memory usage
- Optimized for Go's slice operations
- Minimal overhead over standard library

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines

- Follow Go best practices and idioms
- Add comprehensive tests for new features
- Update documentation for API changes
- Ensure all tests pass before submitting PR
- Use meaningful commit messages

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by functional programming patterns from various languages
- Built on Go's excellent standard library
- Community feedback and contributions

---

Made with ❤️ for the Go community
