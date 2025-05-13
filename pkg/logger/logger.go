package logger

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	instance *Logger
	once     sync.Once
)

// Logger wraps zap.SugaredLogger for structured logging
type Logger struct {
	*zap.SugaredLogger
}

// GetLogger returns the singleton logger instance
func GetLogger() *Logger {
	once.Do(func() {
		encoderConfig := zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "message",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}

		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

		core := zapcore.NewCore(
			consoleEncoder,
			zapcore.Lock(os.Stdout),
			zapcore.DebugLevel,
		)

		baseLogger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
		instance = &Logger{baseLogger.Sugar()}
	})
	return instance
}

// WithField creates a logger with an additional field
func (log *Logger) WithField(key string, value any) *Logger {
	return &Logger{log.SugaredLogger.With(key, value)}
}

// WithFields creates a logger with multiple additional fields
func (log *Logger) WithFields(fields map[string]any) *Logger {
	return &Logger{log.SugaredLogger.With(fields)}
}

// Sync flushes any buffered log entries
func (log *Logger) Sync() {
	_ = log.SugaredLogger.Sync()
}