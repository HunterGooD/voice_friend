package logger

import (
	"io"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

type logrusLogger struct {
	l *logrus.Logger
}

func NewTextLogrusLogger(w io.Writer, log_level string) *logrusLogger {
	log := logrus.New()
	switch strings.ToUpper(log_level) {
	case "INFO":
		log.SetLevel(logrus.InfoLevel)
	case "DEBUG":
		log.SetLevel(logrus.DebugLevel)
	}
	log.SetFormatter(&logrus.TextFormatter{})

	return &logrusLogger{log}
}

func NewJsonLogrusLogger(w io.Writer, log_level string) *logrusLogger {
	log := logrus.New()
	switch strings.ToUpper(log_level) {
	case "INFO":
		log.SetLevel(logrus.InfoLevel)
	case "DEBUG":
		log.SetLevel(logrus.DebugLevel)
	default:
		log.SetLevel(logrus.DebugLevel)
	}
	log.SetFormatter(&logrus.JSONFormatter{})

	return &logrusLogger{log}
}

func (log *logrusLogger) Info(message string, opt ...any) {
	params := parseLogrusOpt(opt...)
	log.l.WithFields(params).Info(message)
}

func (log *logrusLogger) Debug(message string, opt ...any) {
	params := parseLogrusOpt(opt...)
	log.l.WithFields(params).Debug(message)
}

func (log *logrusLogger) Warn(message string, opt ...any) {
	params := parseLogrusOpt(opt...)
	log.l.WithFields(params).Warn(message)
}

func (log *logrusLogger) Error(message string, opt ...any) {
	params := parseLogrusOpt(opt...)
	log.l.WithFields(params).Error(message)
}

func parseLogrusOpt(opt ...any) logrus.Fields {
	params := make(logrus.Fields)
	for k, v := range opt {
		switch val := v.(type) {
		case map[string]any:
			// TODO: if use map any maybe using this func
			// logrus.Fields(val)
			for key, value := range val {
				params[key] = value
			}
		default:
			params["param_"+strconv.Itoa(k)] = val
		}
	}
	return params
}
