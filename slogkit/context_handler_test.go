package slogkit_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"strings"
	"testing"

	"github.com/entrylens/godash/slogkit"
	"github.com/stretchr/testify/suite"
)

type reqIdCtxKey struct{}

type ContextHandlerSuite struct {
	suite.Suite
	buf *bytes.Buffer
}

func (s *ContextHandlerSuite) SetupTest() {
	s.buf = &bytes.Buffer{}
}

// NewContextHandler should create a JSON handler when UseJson is true
func (s *ContextHandlerSuite) TestNewContextHandler_JSONHandler() {
	handler := slogkit.NewContextHandler(slogkit.ContextHandlerOptions{
		UseJson: true,
		Writer:  s.buf,
		Level:   slog.LevelInfo,
	})

	logger := slog.New(handler)
	logger.Info("test message", "key", "value")

	s.Contains(s.buf.String(), "test message")
	s.Contains(s.buf.String(), "key")
	s.Contains(s.buf.String(), "value")
	// JSON handler should produce JSON output
	var jsonData map[string]interface{}
	err := json.Unmarshal(s.buf.Bytes(), &jsonData)
	s.NoError(err)
}

// NewContextHandler should create a text handler when UseJson is false
func (s *ContextHandlerSuite) TestNewContextHandler_TextHandler() {
	handler := slogkit.NewContextHandler(slogkit.ContextHandlerOptions{
		UseJson: false,
		Writer:  s.buf,
		Level:   slog.LevelInfo,
	})

	logger := slog.New(handler)
	logger.Info("test message", "key", "value")

	output := s.buf.String()
	s.Contains(output, "test message")
	s.Contains(output, "key")
	s.Contains(output, "value")
}

// NewContextHandler should add PID when WithPID is true
func (s *ContextHandlerSuite) TestNewContextHandler_WithPID() {
	handler := slogkit.NewContextHandler(slogkit.ContextHandlerOptions{
		UseJson: true,
		Writer:  s.buf,
		Level:   slog.LevelInfo,
		WithPID: true,
	})

	logger := slog.New(handler)
	logger.Info("test message")

	var jsonData map[string]interface{}
	err := json.Unmarshal(s.buf.Bytes(), &jsonData)
	s.NoError(err)
	s.Contains(jsonData, "pid")
}

// NewContextHandler should use custom PID key when PIDKey is provided
func (s *ContextHandlerSuite) TestNewContextHandler_CustomPIDKey() {
	handler := slogkit.NewContextHandler(slogkit.ContextHandlerOptions{
		UseJson: true,
		Writer:  s.buf,
		Level:   slog.LevelInfo,
		WithPID: true,
		PIDKey:  "process_id",
	})

	logger := slog.New(handler)
	logger.Info("test message")

	var jsonData map[string]interface{}
	err := json.Unmarshal(s.buf.Bytes(), &jsonData)
	s.NoError(err)
	s.Contains(jsonData, "process_id")
	s.NotContains(jsonData, "pid")
}

// NewContextHandler should add extra attributes when ExtraAttrs is provided
func (s *ContextHandlerSuite) TestNewContextHandler_ExtraAttrs() {
	handler := slogkit.NewContextHandler(slogkit.ContextHandlerOptions{
		UseJson: true,
		Writer:  s.buf,
		Level:   slog.LevelInfo,
		ExtraAttrs: []slog.Attr{
			slog.String("service", "test-service"),
			slog.Int("version", 1),
		},
	})

	logger := slog.New(handler)
	logger.Info("test message")

	var jsonData map[string]interface{}
	err := json.Unmarshal(s.buf.Bytes(), &jsonData)
	s.NoError(err)
	s.Equal("test-service", jsonData["service"])
	s.Equal(float64(1), jsonData["version"]) // JSON numbers are float64
}

// Handle should add source information when AddSource is true
func (s *ContextHandlerSuite) TestHandle_WithSource() {
	// Test JSON handler with source
	handler := slogkit.NewContextHandler(slogkit.ContextHandlerOptions{
		UseJson:   true,
		Writer:    s.buf,
		Level:     slog.LevelInfo,
		AddSource: true,
	})

	logger := slog.New(handler)
	logger.Info("test message")

	var jsonData map[string]interface{}
	err := json.Unmarshal(s.buf.Bytes(), &jsonData)
	s.NoError(err)
	s.Contains(jsonData, "source")
	source, ok := jsonData["source"].(string)
	s.True(ok)
	s.Contains(source, ".go")
	s.Contains(source, ":")

	// Test text handler with source
	s.buf.Reset()
	handler = slogkit.NewContextHandler(slogkit.ContextHandlerOptions{
		UseJson:   false,
		Writer:    s.buf,
		Level:     slog.LevelInfo,
		AddSource: true,
	})

	logger = slog.New(handler)
	logger.Info("test message")

	output := s.buf.String()
	s.Contains(output, "test message")
	s.Contains(output, "source")
	s.True(strings.Contains(output, ".go") || strings.Contains(output, "context_handler_test.go"))
}

// Handle should use custom source key when SourceKey is provided
func (s *ContextHandlerSuite) TestHandle_CustomSourceKey() {
	handler := slogkit.NewContextHandler(slogkit.ContextHandlerOptions{
		UseJson:   true,
		Writer:    s.buf,
		Level:     slog.LevelInfo,
		AddSource: true,
		SourceKey: "file_location",
	})

	logger := slog.New(handler)
	logger.Info("test message")

	var jsonData map[string]interface{}
	err := json.Unmarshal(s.buf.Bytes(), &jsonData)
	s.NoError(err)
	s.Contains(jsonData, "file_location")
	s.NotContains(jsonData, "source")
}

// Handle should not add attributes from context when context is Background
func (s *ContextHandlerSuite) TestHandle_BackgroundContext() {
	handler := slogkit.NewContextHandler(slogkit.ContextHandlerOptions{
		UseJson: true,
		Writer:  s.buf,
		Level:   slog.LevelInfo,
		AppendAttrFromContext: func(ctx context.Context) ([]slog.Attr, error) {
			return []slog.Attr{slog.String("from_context", "value")}, nil
		},
	})

	logger := slog.New(handler)
	logger.Info("test message")

	var jsonData map[string]interface{}
	err := json.Unmarshal(s.buf.Bytes(), &jsonData)
	s.NoError(err)
	s.NotContains(jsonData, "from_context")
}

// Handle should add attributes from context when AppendAttrFromContext is provided
func (s *ContextHandlerSuite) TestHandle_WithAppendAttrFromContext() {
	handler := slogkit.NewContextHandler(slogkit.ContextHandlerOptions{
		UseJson: true,
		Writer:  s.buf,
		Level:   slog.LevelInfo,
		AppendAttrFromContext: func(ctx context.Context) ([]slog.Attr, error) {
			return []slog.Attr{
				slog.String("request_id", "12345"),
				slog.String("user_id", "user-1"),
			}, nil
		},
	})

	logger := slog.New(handler)
	ctx := context.WithValue(context.Background(), reqIdCtxKey{}, "12345")
	logger.Log(ctx, slog.LevelInfo, "test message")

	var jsonData map[string]interface{}
	err := json.Unmarshal(s.buf.Bytes(), &jsonData)
	s.NoError(err)
	s.Equal("12345", jsonData["request_id"])
	s.Equal("user-1", jsonData["user_id"])
}

// Handle should return error when AppendAttrFromContext returns error
func (s *ContextHandlerSuite) TestHandle_AppendAttrFromContextError() {
	testErr := errors.New("context error")
	handler := slogkit.NewContextHandler(slogkit.ContextHandlerOptions{
		UseJson: true,
		Writer:  s.buf,
		Level:   slog.LevelInfo,
		AppendAttrFromContext: func(ctx context.Context) ([]slog.Attr, error) {
			return []slog.Attr{
				slog.String("request_id", "12345"),
				slog.String("user_id", "user-1"),
			}, testErr
		},
	})

	logger := slog.New(handler)
	ctx := context.WithValue(context.Background(), reqIdCtxKey{}, "12345")
	logger.Log(ctx, slog.LevelInfo, "test message")

	var jsonData map[string]interface{}
	err := json.Unmarshal(s.buf.Bytes(), &jsonData)
	s.NoError(err)
	s.NotContains(jsonData, "request_id")
}

// WithAttrs should preserve ContextHandler wrapper and add static attributes
func (s *ContextHandlerSuite) TestWithAttrs_PreservesWrapper() {
	handler := slogkit.NewContextHandler(slogkit.ContextHandlerOptions{
		UseJson: true,
		Writer:  s.buf,
		Level:   slog.LevelInfo,
	})

	newHandler := handler.WithAttrs([]slog.Attr{slog.String("static", "attr")})

	logger := slog.New(newHandler)
	logger.Info("test message")

	var jsonData map[string]interface{}
	err := json.Unmarshal(s.buf.Bytes(), &jsonData)
	s.NoError(err)
	s.Equal("attr", jsonData["static"])
	s.Equal("test message", jsonData["msg"])
}

// WithGroup should preserve ContextHandler wrapper and group attributes
func (s *ContextHandlerSuite) TestWithGroup_PreservesWrapper() {
	handler := slogkit.NewContextHandler(slogkit.ContextHandlerOptions{
		UseJson: true,
		Writer:  s.buf,
		Level:   slog.LevelInfo,
	})

	newHandler := handler.WithGroup("group")

	logger := slog.New(newHandler)
	logger.LogAttrs(context.Background(), slog.LevelInfo, "test message", slog.String("nested", "value"))

	var jsonData map[string]interface{}
	err := json.Unmarshal(s.buf.Bytes(), &jsonData)
	s.NoError(err)
	// Check that group structure is preserved
	group, ok := jsonData["group"].(map[string]interface{})
	s.True(ok)
	s.Equal("value", group["nested"])
	s.Equal("test message", jsonData["msg"])
}

// Handle should respect log level
func (s *ContextHandlerSuite) TestHandle_RespectsLogLevel() {
	handler := slogkit.NewContextHandler(slogkit.ContextHandlerOptions{
		UseJson: false,
		Writer:  s.buf,
		Level:   slog.LevelWarn,
	})

	logger := slog.New(handler)
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message")

	output := s.buf.String()
	s.NotContains(output, "info message")
	s.Contains(output, "warn message")
	s.Contains(output, "error message")
}

// Handle should combine all features together
func (s *ContextHandlerSuite) TestHandle_AllFeatures() {
	handler := slogkit.NewContextHandler(slogkit.ContextHandlerOptions{
		UseJson:   true,
		Writer:    s.buf,
		Level:     slog.LevelInfo,
		WithPID:   true,
		PIDKey:    "process_id",
		AddSource: true,
		SourceKey: "file_location",
		ExtraAttrs: []slog.Attr{
			slog.String("service", "test"),
		},
		AppendAttrFromContext: func(ctx context.Context) ([]slog.Attr, error) {
			return []slog.Attr{slog.String("request_id", "req-123")}, nil
		},
	})

	logger := slog.New(handler)
	ctx := context.WithValue(context.Background(), reqIdCtxKey{}, "req-123")
	logger.LogAttrs(ctx, slog.LevelInfo, "test message", slog.String("key", "value"))

	var jsonData map[string]interface{}
	err := json.Unmarshal(s.buf.Bytes(), &jsonData)
	s.NoError(err)
	s.Contains(jsonData, "process_id")
	s.Contains(jsonData, "file_location")
	s.Equal("test", jsonData["service"])
	s.Equal("req-123", jsonData["request_id"])
	s.Equal("test message", jsonData["msg"])
	s.Equal("value", jsonData["key"])
}

func (s *ContextHandlerSuite) TestHandle_WithAttrs() {
	handler := slogkit.NewContextHandler(slogkit.ContextHandlerOptions{
		UseJson:   true,
		Writer:    s.buf,
		AddSource: true,
		Level:     slog.LevelInfo,
		ExtraAttrs: []slog.Attr{
			slog.String("service", "test"),
		},
	})

	logger := slog.New(handler).With(slog.String("static", "attr"))
	logger.Info("test message")

	var jsonData map[string]interface{}
	err := json.Unmarshal(s.buf.Bytes(), &jsonData)
	s.NoError(err)
	s.Contains(jsonData, "static")
	s.Contains(jsonData, "service")
	s.Contains(jsonData, "source")
}

func TestContextHandlerSuite(t *testing.T) {
	suite.Run(t, new(ContextHandlerSuite))
}
