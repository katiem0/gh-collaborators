package log

import (
	"go.uber.org/zap"
)

func NewLogger(debug bool) (*zap.Logger, error) {
	level := zap.InfoLevel

	if debug {
		level = zap.DebugLevel
	}

	loggerConfig := zap.Config{
		Level:            zap.NewAtomicLevelAt(level),
		Encoding:         "console",
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
		// Disable stack traces for cleaner output
		DisableStacktrace: !debug, // Only show stack traces in debug mode
	}

	// Customize encoder config for better readability
	if debug {
		loggerConfig.EncoderConfig.StacktraceKey = "stacktrace"
	} else {
		loggerConfig.EncoderConfig.StacktraceKey = ""
	}

	return loggerConfig.Build()
}
