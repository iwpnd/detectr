package logger

import (
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var defaultConfig = zap.Config{
	Encoding:         "json",
	Level:            zap.NewAtomicLevelAt(zapcore.DebugLevel),
	OutputPaths:      []string{"stdout"},
	ErrorOutputPaths: []string{"stderr"},
	EncoderConfig: zapcore.EncoderConfig{
		MessageKey: "message",

		LevelKey:    "level",
		EncodeLevel: zapcore.CapitalLevelEncoder,

		TimeKey:    "time",
		EncodeTime: zapcore.ISO8601TimeEncoder,
	},
}

func getLogLevelFromString(level string) (zapcore.Level, error) {
	l := strings.ToLower(level)

	switch l {
	case "debug":
		return zapcore.DebugLevel, nil
	case "info":
		return zapcore.InfoLevel, nil
	case "warn":
		return zapcore.WarnLevel, nil
	case "error":
		return zapcore.ErrorLevel, nil
	default:
		return 0, &ErrInvalidLogLevel{Level: level}
	}
}

func SetLogLevel(level string) error {
	l, err := getLogLevelFromString(level)
	if err != nil {
		return err
	}

	defaultConfig.Level = zap.NewAtomicLevelAt(l)
	return nil
}

func New() (*zap.Logger, error) {
	logger, err := defaultConfig.Build()

	if err != nil {
		return nil, err
	}

	return logger, nil
}
