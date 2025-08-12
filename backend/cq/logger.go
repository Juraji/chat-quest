package cq

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func SetupLogger(enableFileLogging bool) (*zap.Logger, error) {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		MessageKey:     "msg",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	consoleWriter := zapcore.AddSync(os.Stdout)

	var core zapcore.Core
	if enableFileLogging {
		file, err := os.OpenFile("app.logger", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}

		fileWriter := zapcore.AddSync(file)
		multiWriter := zapcore.NewMultiWriteSyncer(consoleWriter, fileWriter)

		core = zapcore.NewCore(
			encoder,
			multiWriter,
			zap.DebugLevel,
		)
	} else {
		core = zapcore.NewCore(
			encoder,
			consoleWriter,
			zap.DebugLevel,
		)
	}

	logger := zap.New(
		core,
		zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel), // Add stack trace for error logs
	)

	return logger, nil
}
