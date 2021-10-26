package logger

import (
	"go.uber.org/zap"
)

// @author valor.

type logger struct {
	log   *zap.Logger
	suagr *zap.SugaredLogger
}

func (l *logger) sync() error {
	return l.log.Sync()
}

func (l *logger) debug(msg string, fields ...zap.Field) {
	l.log.Debug(msg, fields...)
}

func (l *logger) info(msg string, fields ...zap.Field) {
	l.log.Info(msg, fields...)
}

func (l *logger) warn(msg string, fields ...zap.Field) {
	l.log.Warn(msg, fields...)
}

func (l *logger) error(msg string, fields ...zap.Field) {
	l.log.Error(msg, fields...)
}

func (l *logger) panic(msg string, fields ...zap.Field) {
	l.log.Panic(msg, fields...)
}

func (l *logger) fatal(msg string, fields ...zap.Field) {
	l.log.Fatal(msg, fields...)
}

func (l *logger) debugf(template string, args ...interface{}) {
	l.suagr.Debugf(template, args...)
}

func (l *logger) infof(template string, args ...interface{}) {
	l.suagr.Infof(template, args...)
}

func (l *logger) warnf(template string, args ...interface{}) {
	l.suagr.Warnf(template, args...)
}

func (l *logger) errorf(template string, args ...interface{}) {
	l.suagr.Errorf(template, args...)
}

func (l *logger) panicf(template string, args ...interface{}) {
	l.suagr.Panicf(template, args...)
}

func (l *logger) fatalf(template string, args ...interface{}) {
	l.suagr.Fatalf(template, args...)
}
