package logger

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// @author valor.

var l *logger
var once sync.Once

func InitConsole(tmfmt string, highlight bool, level string) {
	once.Do(func() {
		encoder := getConsoleEncoder(tmfmt, highlight)
		writeSyncer := getConsoleWriteSyncer()
		enabler := getLevelEnabler(level, consoleLevel)
		core := zapcore.NewCore(encoder, writeSyncer, enabler)

		options := getOptions()
		log := zap.New(core, options...)
		suagr := log.Sugar()

		l = &logger{
			log: log,
			suagr: suagr,
		}
	})
}

func InitLogfile(tmfmt string, logfile string, maxLineNum int, level string) {
	once.Do(func() {
		encoder := getLogfileEncoder(tmfmt)
		writeSyncer := getLogfileWriteSyncer(logfile, maxLineNum)
		enabler := getLevelEnabler(level, logfileLevel)
		core := zapcore.NewCore(encoder, writeSyncer, enabler)

		options := getOptions()
		log := zap.New(core, options...)
		suagr := log.Sugar()

		l = &logger{
			log: log,
			suagr: suagr,
		}
	})
}

func Sync() error {
	var err error = nil

	if l != nil {
		err = l.sync()
	}
	return err
}

func Debug(msg string, fields ...zap.Field) {
	if l != nil {
		l.debug(msg, fields...)
	}
}

func Info(msg string, fields ...zap.Field) {
	if l != nil {
		l.info(msg, fields...)
	}
}

func Warn(msg string, fields ...zap.Field) {
	if l != nil {
		l.warn(msg, fields...)
	}
}

func Error(msg string, fields ...zap.Field) {
	if l != nil {
		l.error(msg, fields...)
	}
}

func Panic(msg string, fields ...zap.Field) {
	if l != nil {
		l.panic(msg, fields...)
	}
}

func Fatal(msg string, fields ...zap.Field) {
	if l != nil {
		l.fatal(msg, fields...)
	}
}

func Debugf(template string, args ...interface{}) {
	if l != nil {
		l.debugf(template, args...)
	}
}

func Infof(template string, args ...interface{}) {
	if l != nil {
		l.infof(template, args...)
	}
}

func Warnf(template string, args ...interface{}) {
	if l != nil {
		l.warnf(template, args...)
	}
}

func Errorf(template string, args ...interface{}) {
	if l != nil {
		l.errorf(template, args...)
	}
}

func Panicf(template string, args ...interface{}) {
	if l != nil {
		l.panicf(template, args...)
	}
}

func Fatalf(template string, args ...interface{}) {
	if l != nil {
		l.fatalf(template, args...)
	}
}
