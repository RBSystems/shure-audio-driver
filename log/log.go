package log

import (
	"fmt"
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// L is a plain zap logger
var L *zap.Logger

// Config is the logger config used for L
var Config zap.Config

func StartLogger() {
	var err error
	Config = zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "@",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	L, err = Config.Build()
	if err != nil {
		log.Fatalf("failed to initialize zap logger: %v", err)
	}

	L.Info("Zap logger started")

	_ = L.Sync()
}

func SetLogLevel(level string) error {
	switch level {
	case "Debug":
		fmt.Printf("\nSetting log level to *debug*\n\n")
		Config.Level.SetLevel(zap.DebugLevel)
	case "Info":
		fmt.Printf("\nSetting log level to *info*\n\n")
		Config.Level.SetLevel(zap.InfoLevel)
	case "Warn":
		fmt.Printf("\nSetting log level to *warn*\n\n")
		Config.Level.SetLevel(zap.WarnLevel)
	case "Error":
		fmt.Printf("\nSetting log level to *error*\n\n")
		Config.Level.SetLevel(zap.ErrorLevel)
	case "Panic":
		fmt.Printf("\nSetting log level to *panic*\n\n")
		Config.Level.SetLevel(zap.PanicLevel)
	default:
		return fmt.Errorf("invalid log level: must be [1-4]")
	}
	return nil
}
