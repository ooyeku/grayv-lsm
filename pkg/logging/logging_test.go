package logging

import (
	"bytes"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNewColorfulLogger(t *testing.T) {
	logger := NewColorfulLogger()

	assert.NotNil(t, logger)
	assert.Equal(t, logrus.InfoLevel, logger.GetLevel())
	assert.Equal(t, logrus.InfoLevel, logger.Logger.GetLevel())
}

func TestColorfulLogger_SetLevel(t *testing.T) {
	logger := NewColorfulLogger()
	logger.SetLevel(logrus.DebugLevel)

	assert.Equal(t, logrus.DebugLevel, logger.GetLevel())
	assert.Equal(t, logrus.DebugLevel, logger.Logger.GetLevel())
}

func TestColorfulLogger_SetOutput(t *testing.T) {
	logger := NewColorfulLogger()
	var buf bytes.Buffer
	logger.SetOutput(&buf)

	logger.Info("test message")
	assert.Contains(t, buf.String(), "test message")
}

func TestColorfulLogger_Debug(t *testing.T) {
	logger := NewColorfulLogger()
	var buf bytes.Buffer
	logger.SetOutput(&buf)
	logger.SetLevel(logrus.DebugLevel)

	logger.Debug("debug message")
	assert.Contains(t, buf.String(), "debug message")
}

func TestColorfulLogger_Info(t *testing.T) {
	logger := NewColorfulLogger()
	var buf bytes.Buffer
	logger.SetOutput(&buf)

	logger.Info("info message")
	assert.Contains(t, buf.String(), "info message")
}

func TestColorfulLogger_Warn(t *testing.T) {
	logger := NewColorfulLogger()
	var buf bytes.Buffer
	logger.SetOutput(&buf)

	logger.Warn("warn message")
	assert.Contains(t, buf.String(), "warn message")
}

func TestColorfulLogger_Error(t *testing.T) {
	logger := NewColorfulLogger()
	var buf bytes.Buffer
	logger.SetOutput(&buf)

	logger.Error("error message")
	assert.Contains(t, buf.String(), "error message")
}

func TestColorfulLogger_Panic(t *testing.T) {
	logger := NewColorfulLogger()
	var buf bytes.Buffer
	logger.SetOutput(&buf)

	// To avoid stopping the test, we use a deferred function to recover from the panic
	defer func() {
		if r := recover(); r != nil {
			assert.Contains(t, buf.String(), "panic message")
		}
	}()

	logger.Panic("panic message")
}
