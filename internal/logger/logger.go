package logger

import (
	"errors"
	"log/slog"
	"os"
	"strings"
)

const (
	logLevelEnvName = "IMPR_LOG_LEVEL"
)

var errInvalidLogLevel = errors.New(`invalid log level, expected "info", "warn", "error" or "debug"`)

type Logger struct {
	logger *slog.Logger
}

// New logger.
func New() (*Logger, error) {
	// Get log level.
	env := strings.ToLower(os.Getenv(logLevelEnvName))

	level := new(slog.LevelVar)

	switch env {
	case "":
		level.Set(slog.LevelError)
	case "info":
		level.Set(slog.LevelInfo)
	case "warn":
		level.Set(slog.LevelWarn)
	case "error":
		level.Set(slog.LevelError)
	case "debug":
		level.Set(slog.LevelDebug)
	default:
		return nil, errInvalidLogLevel
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))

	return &Logger{logger: logger}, nil
}

func (l *Logger) Info(message string) {
	l.logger.Info(message)
}

func (l *Logger) Warn(message string) {
	l.logger.Warn(message)
}

func (l *Logger) Error(message string) {
	l.logger.Error(message)
}

func (l *Logger) Debug(message string) {
	l.logger.Debug(message)
}
