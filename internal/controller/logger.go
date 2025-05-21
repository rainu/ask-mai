package controller

import (
	"github.com/wailsapp/wails/v2/pkg/logger"
	"log/slog"
)

type logWrapper struct {
	log *slog.Logger
}

func (l logWrapper) Print(message string) {
	l.log.With("level", "print").Info(message)
}

func (l logWrapper) Trace(message string) {
	l.log.With("level", "trace").Debug(message)
}

func (l logWrapper) Debug(message string) {
	l.log.With("level", "debug").Debug(message)
}

func (l logWrapper) Info(message string) {
	l.log.With("level", "info").Info(message)
}

func (l logWrapper) Warning(message string) {
	l.log.With("level", "warning").Warn(message)
}

func (l logWrapper) Error(message string) {
	l.log.With("level", "error").Error(message)
}

func (l logWrapper) Fatal(message string) {
	l.log.With("level", "fatal").Error(message)
}

func newDefaultLogger() logger.Logger {
	return &logWrapper{
		log: slog.With("source", "wails"),
	}
}
