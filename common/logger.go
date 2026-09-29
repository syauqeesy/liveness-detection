package common

import (
	"log/slog"
	"os"
)

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

type logger struct {
	logger *slog.Logger
}

func NewLogger(service string, environment string) *logger {
	handler := slog.NewJSONHandler(
		os.Stdout,
		&slog.HandlerOptions{
			Level: slog.LevelInfo,
			ReplaceAttr: func(
				groups []string,
				attr slog.Attr,
			) slog.Attr {
				if attr.Key == slog.TimeKey {
					return slog.String(
						slog.TimeKey,
						attr.Value.Time().UTC().Format("2006-01-02T15:04:05.999999999+00:00"),
					)
				}
				return attr
			},
		},
	)

	return &logger{
		logger: slog.New(handler).With("service", service, "environment", environment),
	}
}

func (l *logger) Debug(msg string, args ...any) {
	l.logger.Debug(msg, args...)
}

func (l *logger) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

func (l *logger) Warn(msg string, args ...any) {
	l.logger.Warn(msg, args...)
}

func (l *logger) Error(msg string, args ...any) {
	l.logger.Error(msg, args...)
}
