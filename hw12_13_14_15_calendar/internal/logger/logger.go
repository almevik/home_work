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
		},
	}

	logg, err := cfg.Build()
	if err != nil {
		return nil, fmt.Errorf("failed build logger: %w", err)
	}

	return &logger{logg}, nil
}

func (l logger) Debug(args ...interface{}) {
	l.Logger.Debug(fmt.Sprintf("%v", args))
}

func (l logger) Info(args ...interface{}) {
	l.Logger.Info(fmt.Sprintf("%v", args))
}

func (l logger) Warn(args ...interface{}) {
	l.Logger.Warn(fmt.Sprintf("%v", args))
}

func (l logger) Error(args ...interface{}) {
	l.Logger.Error(fmt.Sprintf("%v", args))
}

func (l logger) Fatal(args ...interface{}) {
	l.Logger.Fatal(fmt.Sprintf("%v", args))
}
