package slogkit

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime"
)

type ContextHandlerOptions struct {
	UseJson               bool
	Level                 slog.Level
	AddSource             bool
	SourceKey             string
	Writer                io.Writer
	WithPID               bool   // append pid to the log
	PIDKey                string // the pid key name in the log
	CallerSkip            int
	ExtraAttrs            []slog.Attr
	AppendAttrFromContext func(ctx context.Context) ([]slog.Attr, error)
}

type ContextHandler struct {
	slog.Handler
	AddSource             bool
	SourceKey             string
	CallerSkip            int
	AppendAttrFromContext func(ctx context.Context) ([]slog.Attr, error)
}

func NewContextHandler(options ContextHandlerOptions) *ContextHandler {

	slogOpts := slog.HandlerOptions{
		Level: options.Level,
		// disable default slog source tracking, we will use our own source tracking
		AddSource: false,
	}

	var handler slog.Handler

	var writer = options.Writer
	if writer == nil {
		writer = os.Stdout
	}

	if !options.UseJson {
		handler = slog.NewTextHandler(writer, &slogOpts)
	} else {
		handler = slog.NewJSONHandler(writer, &slogOpts)
	}

	// Add PID
	if options.WithPID {
		pidKey := "pid"
		if options.PIDKey != "" {
			pidKey = options.PIDKey
		}
		handler = handler.WithAttrs([]slog.Attr{slog.Int(pidKey, os.Getpid())})
	}

	// Add append attrs
	if len(options.ExtraAttrs) > 0 {
		handler = handler.WithAttrs(options.ExtraAttrs)
	}

	return &ContextHandler{
		Handler:               handler,
		AddSource:             options.AddSource,
		SourceKey:             options.SourceKey,
		CallerSkip:            options.CallerSkip,
		AppendAttrFromContext: options.AppendAttrFromContext,
	}
}

// Append different fields according to the context
func (h ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if h.AddSource {
		skip := h.CallerSkip
		if skip == 0 {
			skip = 3
		}

		_, file, line, ok := runtime.Caller(skip)
		if ok {
			source := fmt.Sprintf("%s:%d", file, line)

			sourceKey := "source"
			if h.SourceKey != "" {
				sourceKey = h.SourceKey
			}

			r.AddAttrs(slog.String(sourceKey, source))
		}
	}

	if ctx == context.Background() {
		return h.Handler.Handle(ctx, r)
	}

	if h.AppendAttrFromContext != nil {
		attrs, err := h.AppendAttrFromContext(ctx)
		if err != nil {
			slog.Error("failed to append attributes from context", slog.String("error", err.Error()))
		} else {
			r.AddAttrs(attrs...)
		}

	}

	return h.Handler.Handle(ctx, r)
}

// WithAttrs returns a new handler with the given attributes added to all records.
// This ensures the ContextHandler wrapper is preserved when .With() is called on the logger.
func (h ContextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return ContextHandler{
		Handler:               h.Handler.WithAttrs(attrs),
		AddSource:             h.AddSource,
		SourceKey:             h.SourceKey,
		CallerSkip:            h.CallerSkip,
		AppendAttrFromContext: h.AppendAttrFromContext,
	}
}

// WithGroup returns a new handler that starts a group with the given name.
// This ensures the ContextHandler wrapper is preserved when grouping is used.
func (h ContextHandler) WithGroup(name string) slog.Handler {
	return ContextHandler{
		Handler:               h.Handler.WithGroup(name),
		AddSource:             h.AddSource,
		SourceKey:             h.SourceKey,
		CallerSkip:            h.CallerSkip,
		AppendAttrFromContext: h.AppendAttrFromContext,
	}
}
