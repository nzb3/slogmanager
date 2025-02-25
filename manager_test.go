package slogmanager

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ManagerTestSuite struct {
	suite.Suite
	manager *Manager
}

func TestManagerTestSuite(t *testing.T) {
	suite.Run(t, new(ManagerTestSuite))
}

func (s *ManagerTestSuite) SetupTest() {
	s.manager = New()
}

func (s *ManagerTestSuite) TestNew() {
	assert.NotNil(s.T(), s.manager, "Manager should not be nil")
	assert.NotNil(s.T(), s.manager.writers, "Writers map should be initialized")
	assert.Empty(s.T(), s.manager.writers, "Writers map should be empty initially")
	assert.NotNil(s.T(), s.manager.logger, "Default logger should be initialized")
}

func (s *ManagerTestSuite) TestAddWriter() {
	buf := &bytes.Buffer{}
	w := NewWriter(buf)

	s.manager.AddWriter("buf", w)

	assert.Len(s.T(), s.manager.writers, 1, "Should have one writer")
	assert.Equal(s.T(), w, s.manager.writers["buf"], "Writer should be stored with correct key")
}

func (s *ManagerTestSuite) TestRemoveWriter() {
	buf1 := &bytes.Buffer{}
	buf2 := &bytes.Buffer{}
	w1 := NewWriter(buf1, WithJSONFormat())
	w2 := NewWriter(buf2, WithTextFormat())

	s.manager.AddWriter("json", w1)
	s.manager.AddWriter("text", w2)
	s.manager.RemoveWriter("json")

	writers := s.manager.Writers()
	assert.Len(s.T(), writers, 1, "Should have one remaining writer")
	assert.Contains(s.T(), writers, "text", "Remaining writer should be 'text'")
	assert.NotContains(s.T(), writers, "json", "'json' writer should be removed")
}

func (s *ManagerTestSuite) TestGetLogger() {
	logger := s.manager.Logger()
	assert.NotNil(s.T(), logger, "Logger should not be nil")
}

func (s *ManagerTestSuite) TestLoggerOutputJSON() {
	jsonBuf := &bytes.Buffer{}
	jsonWriter := NewWriter(jsonBuf, WithJSONFormat())
	s.manager.AddWriter("json", jsonWriter)
	logger := s.manager.Logger()

	logger.Info("test message", "key", "value")

	var jsonOutput map[string]interface{}
	require.NoError(s.T(), json.NewDecoder(jsonBuf).Decode(&jsonOutput), "Should decode JSON output")

	assert.Equal(s.T(), "test message", jsonOutput["msg"], "Message should match")
	assert.Equal(s.T(), "value", jsonOutput["key"], "Key value should match")
}

func (s *ManagerTestSuite) TestLoggerOutputText() {
	textBuf := &bytes.Buffer{}
	textWriter := NewWriter(textBuf, WithTextFormat())
	s.manager.AddWriter("text", textWriter)
	logger := s.manager.Logger()

	logger.Info("test message", "key", "value")

	textOutput := textBuf.String()
	assert.Contains(s.T(), textOutput, "test message", "Should contain message")
	assert.Contains(s.T(), textOutput, "key=value", "Should contain key-value pair")
}

func (s *ManagerTestSuite) TestCreateHandler() {
	buf := &bytes.Buffer{}

	s.T().Run("JSON handler", func(t *testing.T) {
		jsonWriter := NewWriter(buf, WithJSONFormat())
		handler := createHandler(jsonWriter)
		_, ok := handler.(*slog.JSONHandler)
		assert.True(t, ok, "Should create JSON handler")
	})

	s.T().Run("Text handler", func(t *testing.T) {
		textWriter := NewWriter(buf, WithTextFormat())
		handler := createHandler(textWriter)
		_, ok := handler.(*slog.TextHandler)
		assert.True(t, ok, "Should create Text handler")
	})
}

func (s *ManagerTestSuite) TestConcurrentAccess() {
	buf := &bytes.Buffer{}
	w := NewWriter(buf)

	for range 10 {
		go func() {
			s.manager.AddWriter("buf", w)
			s.manager.RemoveWriter("buf")
		}()
	}

	logger := s.manager.Logger()
	for i := range 10 {
		go func(i int) {
			logger.Info("test message", "iteration", i)
		}(i)
	}
}
