// Package slogmanager provides a centralized logging management system using slog.
// It supports multiple writers and concurrent access to logging facilities.
package slogmanager

import (
	"log/slog"
	"sync"

	slogmulti "github.com/samber/slog-multi"
)

// Manager is the central structure for managing multiple log writers and handlers.
// It provides thread-safe operations for managing logging configuration.
type Manager struct {
	mu      sync.RWMutex
	writers map[string]*Writer
	logger  *slog.Logger
}

// New creates and initializes a new Manager instance.
// It initializes an empty map of writers that can be populated later.
//
// Returns:
//   - *Manager: A new manager instance ready for use.
func New() *Manager {
	return &Manager{
		writers: make(map[string]*Writer),
		logger:  slog.Default(),
	}
}

// Logger returns the current logger instance from the manager.
// It provides thread-safe access to the logger.
//
// Returns:
//   - *slog.Logger: The current logger instance.
func (m *Manager) Logger() *slog.Logger {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.logger
}

// Writers returns the current writers map from the manager.
// It provides thread-safe access to the writers.
//
// Returns:
//   - map[string]*Writer: The current writers.
func (m *Manager) Writers() map[string]*Writer {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.writers
}

// AddWriter registers a new writer with the manager.
// It's thread-safe to add a writer and after updates logger
//
// Parameters:
//   - name: The name of writer to add to the manager
//   - writer: The writer instance to add to the manager.
func (m *Manager) AddWriter(name string, writer *Writer) {
	defer m.updateLogger()

	m.mu.Lock()
	defer m.mu.Unlock()

	m.writers[name] = writer
}

// RemoveWriter unregisters a writer from the manager.
// It's thread-safe and automatically updates the logger configuration
// after removing the writer
//
// Parameters:
//   - name: The writer's name to remove from the manager.
func (m *Manager) RemoveWriter(name string) {
	defer m.updateLogger()

	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.writers, name)
}

// updateLogger reconstructs the logger with all currently registered writers.
// It creates appropriate handlers for each writer and sets up a new logger
// instance with multiple outputs using slogmulti.Fanout.
func (m *Manager) updateLogger() {
	handlers := make([]slog.Handler, 0)
	for _, writer := range m.writers {
		handlers = append(handlers, createHandler(writer))
	}

	logger := slog.New(slogmulti.Fanout(handlers...))
	m.logger = logger

	slog.SetDefault(logger)
}

// createHandler generates an appropriate slog.Handler based on the writer's configuration.
// It supports both JSON and text output formats.
//
// Parameters:
//   - writer: The writer instance to create a handler for
//
// Returns:
//   - slog.Handler: A configured handler for the writer.
func createHandler(writer *Writer) slog.Handler {
	if writer.Config.UseJSON {
		return slog.NewJSONHandler(writer, writer.Config.HandlerOpts)
	}

	return slog.NewTextHandler(writer, writer.Config.HandlerOpts)
}
