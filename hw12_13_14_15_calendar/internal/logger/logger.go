package logger

import (
	"fmt"
	
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type logger struct {
	*zap.Logger
}

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
}

func New(level int8, filePath string) (*logger, error) {
	cfg := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(zapcore.Level(level)),
		OutputPaths:      []string{filePath},
		ErrorOutputPaths: []string{filePath},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",

			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,

			TimeKey:    "time",
			EncodeTime: zapcore.ISO8601TimeEncoder,

			CallerKey:    "caller",
			EncodeCaller: zapcore.FullCallerEncoder,
		},
	}

	logg, err := cfg.Build()
	if err != nil {
		return nil, fmt.Errorf("failed build logger: %w", err)
	}

	return &logger{logg}, nil
}

func (l logger) Debug(args ...interface{}) {
	l.Debug(args...)
}

func (l logger) Info(args ...interface{}) {
	l.Info(args...)
}

func (l logger) Warn(args ...interface{}) {
	l.Warn(args...)
}

func (l logger) Error(args ...interface{}) {
	l.Error(args...)
}

func (l logger) Fatal(args ...interface{}) {
	l.Fatal(args...)
}