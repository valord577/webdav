package logger

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// @author valor.

func getBaseEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:       "ts",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		FunctionKey:   zapcore.OmitKey,
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeCaller:  zapcore.ShortCallerEncoder,
	}
}

func getConsoleEncoder(tmfmt string, highlight bool) zapcore.Encoder {
	config := getBaseEncoderConfig()

	if highlight {
		config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		config.EncodeLevel = zapcore.CapitalLevelEncoder
	}
	config.EncodeTime = zapcore.TimeEncoderOfLayout(tmfmt)
	config.EncodeDuration = zapcore.StringDurationEncoder

	return zapcore.NewConsoleEncoder(config)
}

func getLogfileEncoder(tmfmt string) zapcore.Encoder {
	config := getBaseEncoderConfig()

	config.EncodeLevel = zapcore.LowercaseLevelEncoder
	config.EncodeTime = zapcore.TimeEncoderOfLayout(tmfmt)
	config.EncodeDuration = zapcore.SecondsDurationEncoder

	return zapcore.NewConsoleEncoder(config)
}

const consoleLevel = zap.DebugLevel
const logfileLevel = zap.InfoLevel

func getLevelEnabler(level string, def zapcore.Level) zapcore.LevelEnabler {

	switch level {
	case "debug":
		return zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		return zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		return zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		return zap.NewAtomicLevelAt(zap.ErrorLevel)
	case "panic":
		return zap.NewAtomicLevelAt(zap.PanicLevel)
	case "fatal":
		return zap.NewAtomicLevelAt(zap.FatalLevel)

	default:
		return zap.NewAtomicLevelAt(def)
	}
}

func getOptions() []zap.Option {
	options := make([]zap.Option, 0, 2)

	options = append(options, zap.AddCaller())
	options = append(options, zap.AddCallerSkip(2))
	return options
}

type consoleWriteSyncer struct{}

func (s *consoleWriteSyncer) Write(bs []byte) (int, error) {
	return os.Stderr.Write(bs)
}

func (s *consoleWriteSyncer) Sync() error {
	return nil
}

func getConsoleWriteSyncer() zapcore.WriteSyncer {
	return &consoleWriteSyncer{}
}

type logfileWriteSyncer struct {
	f    *os.File
	file string

	line       int
	maxLineNum int

	mu sync.Mutex
}

func (s *logfileWriteSyncer) Write(bs []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// On startup
	if s.f == nil {
		err := s.openNew()
		if err != nil {
			return 0, nil
		}
	}
	if s.line+1 > s.maxLineNum {
		err := s.openNew()
		if err != nil {
			return 0, nil
		}
	}

	n, err := s.f.Write(bs)
	s.line += 1
	return n, err
}

func (s *logfileWriteSyncer) Sync() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return logfileSync(s.f)
}

func (s *logfileWriteSyncer) filename() string {
	if s.file == "" {
		s.file = filepath.Join(os.TempDir(), "webdav/webdav.log")
	}
	return s.file
}

func (s *logfileWriteSyncer) openNew() error {
	filename := s.filename()
	dir := filepath.Dir(filename)

	err := os.MkdirAll(dir, 0744)
	if err != nil {
		return errors.New("can not mkdir for new logfile: " + err.Error())
	}

	_, err = os.Stat(filename)
	if err == nil {
		// File existed. So rename file.
		name := filepath.Base(filename)
		ext := filepath.Ext(name)
		prefix := name[:len(name)-len(ext)]

		backupName := fmt.Sprintf("%s-%d%s", prefix, time.Now().UTC().Unix(), ext)
		e := os.Rename(filename, filepath.Join(dir, backupName))
		if e != nil {
			return errors.New("can not rename old logfile: " + e.Error())
		}
	}

	// Open new file
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return errors.New("can not open new logfile: " + err.Error())
	}

	newf, err := logfileDup2(f, s.f)
	if err != nil {
		return err
	}

	s.f = newf
	s.line = 0
	return nil
}

func getLogfileWriteSyncer(file string, maxLineNum int) zapcore.WriteSyncer {
	return &logfileWriteSyncer{
		file: file,

		maxLineNum: maxLineNum,
	}
}
