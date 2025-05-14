package core

import (
	"log/slog"
	"os"
)

type (
	Logger interface {
		Info(msg string)
		Warn(msg string)
		Debug(msg string)
		Error(msg string)
		With(key string, val any) Logger
		WithError(err error) Logger
	}

	logger struct {
		*slog.Logger
		attrs []slog.Attr
	}
)

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

func (l *logger) WithError(err error) Logger {
	return &logger{
		Logger: l.Logger,
		attrs:  append(l.attrs, slog.String("error", err.Error())),
	}
}

// attrsToArgs は、loggerの属性を汎用的な引数スライスに変換します。
// これにより、ログ関数に渡す前に属性を準備します。
func (l *logger) attrsToArgs() []any {
	args := make([]any, len(l.attrs))
	for i, attr := range l.attrs {
		args[i] = attr
	}
	return args
}

func NewLogger() Logger {
	env := GetEnv("APP_ENV", "dev")

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
