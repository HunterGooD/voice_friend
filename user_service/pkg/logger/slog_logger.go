package logger

import (
	"io"
	"log/slog"
	"strings"
)

type slogLogger struct {
	l *slog.Logger
}

func NewTextSlogLogger(w io.Writer, log_level string) *slogLogger {
	var opts *slog.HandlerOptions
	switch strings.ToUpper(log_level) {
	case "INFO":
		opts = &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}
	case "DEBUG":
		opts = &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}
	}
	l := slog.New(slog.NewTextHandler(w, opts))
	return &slogLogger{l}
}

func NewJsonSlogLogger(w io.Writer, log_level string) *slogLogger {
	var opts *slog.HandlerOptions
	switch strings.ToUpper(log_level) {
	case "INFO":
		opts = &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}
	case "DEBUG":
		opts = &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}
	}
	l := slog.New(slog.NewJSONHandler(w, opts))
	return &slogLogger{l}
}

func (log *slogLogger) Info(message string, opt ...any) {
	params := parseSlogOpt(opt...)
	log.l.Info(message, params...)
}

func (log *slogLogger) Debug(message string, opt ...any) {
	params := parseSlogOpt(opt...)
	log.l.Debug(message, params...)
}

func (log *slogLogger) Warn(message string, opt ...any) {
	params := parseSlogOpt(opt...)
	log.l.Warn(message, params...)
}

func (log *slogLogger) Error(message string, opt ...any) {
	params := parseSlogOpt(opt...)
	log.l.Error(message, params...)
}

func parseSlogOpt(opt ...any) []any {
	params := make([]any, 0)
	for _, v := range opt {
		switch val := v.(type) {
		case map[string]any:
			params = append(params, mapSlogParse(val)...)
		default:
			params = append(params, val)
		}
	}
	return params
}

func mapSlogParse(fields map[string]any) []any {
	res := make([]any, 0)
	for k, v := range fields {
		res = append(res, slog.Any(k, v))
	}
	return res
}
