package log

import (
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestNewLogger(t *testing.T) {
	tests := []struct {
		name     string
		debug    bool
		expected zapcore.Level
	}{
		{
			name:     "info level when debug is false",
			debug:    false,
			expected: zap.InfoLevel,
		},
		{
			name:     "debug level when debug is true",
			debug:    true,
			expected: zap.DebugLevel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := NewLogger(tt.debug)

			if err != nil {
				t.Errorf("NewLogger() error = %v, want nil", err)
				return
			}

			if logger == nil {
				t.Error("NewLogger() returned nil logger")
				return
			}

			// Check if the logger level is set correctly
			if !logger.Core().Enabled(tt.expected) {
				t.Errorf("NewLogger() level not enabled for %v", tt.expected)
			}
		})
	}
}

func TestNewLoggerConfig(t *testing.T) {
	logger, err := NewLogger(true)

	if err != nil {
		t.Errorf("NewLogger() error = %v, want nil", err)
		return
	}

	defer logger.Sync()

	// Test that we can actually log with the created logger
	logger.Info("test info message")
	logger.Debug("test debug message")

	// If we get here without panicking, the logger is working correctly
}

func TestLoggerSingleton(t *testing.T) {
	// Test that multiple calls to NewLogger return valid loggers
	logger1, err1 := NewLogger(false)
	if err1 != nil {
		t.Fatalf("First NewLogger() call failed: %v", err1)
	}

	logger2, err2 := NewLogger(true)
	if err2 != nil {
		t.Fatalf("Second NewLogger() call failed: %v", err2)
	}

	if logger1 == nil || logger2 == nil {
		t.Error("NewLogger() should not return nil")
	}
}
