package logger

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// @author valor.

var l *logger
var once sync.Once

// InitConsole initializes the logger that logs to the console.
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
			log:   log,
			suagr: suagr,
		}
	})
}

// InitLogfile initializes the logger that logs to the logfile.
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
			log:   log,
			suagr: suagr,
		}
	})
}

// Sync calls *zap.Logger.Sync
func Sync() error {
	var err error = nil

	if l != nil {
		err = l.sync()
	}
	return err
}

// Debug calls *zap.Logger.Debug
func Debug(msg string, fields ...zap.Field) {
	if l != nil {
		l.debug(msg, fields...)
	}
}

// Info calls *zap.Logger.Info
func Info(msg string, fields ...zap.Field) {
	if l != nil {
		l.info(msg, fields...)
	}
}

// Warn calls *zap.Logger.Warn
func Warn(msg string, fields ...zap.Field) {
	if l != nil {
		l.warn(msg, fields...)
	}
}

// Error calls *zap.Logger.Error
func Error(msg string, fields ...zap.Field) {
	if l != nil {
		l.error(msg, fields...)
	}
}

// Panic calls *zap.Logger.Panic
func Panic(msg string, fields ...zap.Field) {
	if l != nil {
		l.panic(msg, fields...)
	}
}

// Fatal calls *zap.Logger.Fatal
func Fatal(msg string, fields ...zap.Field) {
	if l != nil {
		l.fatal(msg, fields...)
	}
}

// Debugf calls *zap.SugaredLogger.Debugf
func Debugf(template string, args ...interface{}) {
	if l != nil {
		l.debugf(template, args...)
	}
}

// Infof calls *zap.SugaredLogger.Infof
func Infof(template string, args ...interface{}) {
	if l != nil {
		l.infof(template, args...)
	}
}

// Warnf calls *zap.SugaredLogger.Warnf
func Warnf(template string, args ...interface{}) {
	if l != nil {
		l.warnf(template, args...)
	}
}

// Errorf calls *zap.SugaredLogger.Errorf
func Errorf(template string, args ...interface{}) {
	if l != nil {
		l.errorf(template, args...)
	}
}

// Panicf calls *zap.SugaredLogger.Panicf
func Panicf(template string, args ...interface{}) {
	if l != nil {
		l.panicf(template, args...)
	}
}

// Fatalf calls *zap.SugaredLogger.Fatalf
func Fatalf(template string, args ...interface{}) {
	if l != nil {
		l.fatalf(template, args...)
	}
}
