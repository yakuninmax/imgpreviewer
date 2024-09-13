package logger

import (
	"log/slog"
	"os"
)

type Log struct {
	logger *slog.Logger
}

func New() *Log {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	return &Log{logger: logger}
}

func (l *Log) Info(message string) {
	l.logger.Info(message)
}

func (l *Log) Warn(message string) {
	l.logger.Warn(message)
}

func (l *Log) Error(message string) {
	l.logger.Error(message)
}

func (l *Log) Debug(message string) {
	l.logger.Debug(message)
}
