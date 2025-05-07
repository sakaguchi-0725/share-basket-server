package logger

import (
	"log/slog"
	"os"
)

type logger struct {
	*slog.Logger
	attrs []slog.Attr
}

func (l *logger) Debug(msg string) {
	l.Logger.Debug(msg, l.attrsToArgs()...)
}

func (l *logger) Error(msg string) {
	l.Logger.Error(msg, l.attrsToArgs()...)
}

func (l *logger) Info(msg string) {
	l.Logger.Info(msg, l.attrsToArgs()...)
}

func (l *logger) Warn(msg string) {
	l.Logger.Warn(msg, l.attrsToArgs()...)
}

func (l *logger) With(key string, val any) Logger {
	return &logger{
		Logger: l.Logger,
		attrs:  append(l.attrs, slog.Any(key, val)),
	}
}

func (l *logger) attrsToArgs() []any {
	args := make([]any, len(l.attrs))
	for i, attr := range l.attrs {
		args[i] = attr
	}
	return args
}

func New(env string) Logger {
	var handler slog.Handler

	if env == "prod" {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	}

	return &logger{
		slog.New(handler),
		[]slog.Attr{},
	}
}
