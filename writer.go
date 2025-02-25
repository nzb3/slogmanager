// Package writer provides a customizable writer implementation for structured logging
// with support for both JSON and text formats using the slog package.
package slogmnager

import (
	"io"
	"log/slog"
)

// Writer represents a custom writer that implements io.Writer interface
// and includes configuration options for logging format and handling.
type Writer struct {
	Config *writerConfig
	io.Writer
}

// NewWriter creates a new Writer instance with specified options.
// It uses the functional options pattern for flexible configuration.
//
// Parameters:
//   - writer: An io.Writer interface implementation for actual writing
//   - opts: Variable number of Option functions for configuration
//
// Returns:
//   - *Writer: A configured Writer instance.
func NewWriter(writer io.Writer, opts ...Option) *Writer {
	writerCfg := &writerConfig{
		HandlerOpts: nil,
		UseJSON:     false,
	}

	for _, opt := range opts {
		opt(writerCfg)
	}

	return &Writer{
		Config: writerCfg,
		Writer: writer,
	}
}

// writerConfig holds the configuration options for the Writer.
type writerConfig struct {
	// HandlerOpts contains slog handler options like level, time format, etc.
	HandlerOpts *slog.HandlerOptions

	// UseJSON determines whether to use JSON format (true) or text format (false)
	UseJSON bool
}

// Option defines the function type for writer configuration options.
// It follows the functional options pattern for flexible and extensible configuration.
type Option func(*writerConfig)

// WithTextFormat returns an Option that sets the writer format to text.
// This is the default format if no format option is specified.
//
// Returns:
//   - Option: A function that sets UseJSON to false when applied.
func WithTextFormat() Option {
	return func(o *writerConfig) {
		o.UseJSON = false
	}
}

// WithJSONFormat returns an Option that sets the writer format to JSON.
// This is useful when structured logging output is needed.
//
// Returns:
//   - Option: A function that sets UseJSON to true when applied.
func WithJSONFormat() Option {
	return func(o *writerConfig) {
		o.UseJSON = true
	}
}

// WithSlogHandlerOptions returns an Option that sets custom slog handler options.
// These options can configure various aspects of log handling like level filtering,
// time formats, and source file information.
//
// Parameters:
//   - h: Pointer to slog.HandlerOptions containing desired handler configuration
//
// Returns:
//   - Option: A function that applies the provided handler options when used.
func WithSlogHandlerOptions(h *slog.HandlerOptions) Option {
	return func(o *writerConfig) {
		o.HandlerOpts = h
	}
}
