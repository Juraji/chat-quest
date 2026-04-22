package log

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
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

	// Console logging
	var consoleCore zapcore.Core
	{
		consoleEncoderConfig := encoderConfig
		consoleEncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder := zapcore.NewConsoleEncoder(consoleEncoderConfig)
		writer := zapcore.AddSync(os.Stdout)
		consoleCore = zapcore.NewCore(encoder, writer, level)
	}

	// File logging
	var fileCore zapcore.Core
	{
		logDir := env.MkDataDir("log")

		currentTime := time.Now()
		logFileName := fmt.Sprintf("chat-quest_%s.log", currentTime.Format("2006-01-02"))
		currentFilePath := filepath.Join(logDir, logFileName)

		currentLogFile, err := os.OpenFile(currentFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(fmt.Errorf("error opening log file: %w", err))
		}

		encoder := zapcore.NewConsoleEncoder(encoderConfig)
		writer := zapcore.AddSync(currentLogFile)
		fileCore = zapcore.NewCore(encoder, writer, level)

		go func() {
			// Remove old files.
			logFiles, err := os.ReadDir(logDir)
			if err != nil {
				panic(fmt.Errorf("error listing existing log files: %w", err))
			}
			slices.Reverse(logFiles)

			for idx, file := range logFiles {
				if idx < env.KeepNLogFiles {
					// Skip n files to keep
					continue
				}
				err := os.Remove(filepath.Join(logDir, file.Name()))
				if err != nil {
					panic(fmt.Errorf("failed removing old log file '%s' %w", file.Name(), err))
				}
			}
		}()
	}

	// Signal logging
	signalCore := newSignalCore(level, LogMessagesSignal)

	// Combine cores and set logger instance
	loggerInstance = zap.New(
		zapcore.NewTee(consoleCore, fileCore, signalCore),
		zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel), // Add stack trace for error logs
	)
}
