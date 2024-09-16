package logger

import (
	"log/slog"
	"os"
)

type Logger struct {
	logger *slog.Logger
}

func New() *Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	return &Logger{logger: logger}
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
