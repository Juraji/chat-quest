package log

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"juraji.nl/chat-quest/core"
)

var loggerInstance *zap.Logger

func Get() *zap.Logger {
	if loggerInstance == nil {
		panic("logger not initialized")
	}

	return loggerInstance
}

func InitLogger(env core.Environment) {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		MessageKey:     "msg",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	level := zap.InfoLevel
	if env.DebugEnabled {
		level = zap.DebugLevel
	}

	var consoleCore zapcore.Core
	{
		consoleEncoderConfig := encoderConfig
		consoleEncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder := zapcore.NewConsoleEncoder(consoleEncoderConfig)
		writer := zapcore.AddSync(os.Stdout)
		consoleCore = zapcore.NewCore(encoder, writer, level)
	}

	var fileCore zapcore.Core
	{
		logDir := env.MkDataDir("log")

		currentTime := time.Now()
		currentFileName := fmt.Sprintf("chat-quest_%s.log", currentTime.Format("2006-01-02_15-04-05"))
		currentFilePath := filepath.Join(logDir, currentFileName)

		file, err := os.OpenFile(currentFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(fmt.Errorf("error opening log file: %w", err))
		}

		encoder := zapcore.NewConsoleEncoder(encoderConfig)
		writer := zapcore.AddSync(file)
		fileCore = zapcore.NewCore(encoder, writer, level)
	}

	loggerInstance = zap.New(
		zapcore.NewTee(consoleCore, fileCore),
		zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel), // Add stack trace for error logs
	)
}
