package log

import (
	"time"

	zap "go.uber.org/zap/zapcore"
	"juraji.nl/chat-quest/core/util/signals"
)

type signalLogMessageLevel string

const (
	signalLogMessageDebug signalLogMessageLevel = "DEBUG"
	signalLogMessageInfo  signalLogMessageLevel = "INFO"
	signalLogMessageWarn  signalLogMessageLevel = "WARN"
	signalLogMessageError signalLogMessageLevel = "ERROR"
)

type signalLogMessage struct {
	Level   signalLogMessageLevel `json:"level"`
	Time    time.Time             `json:"time"`
	Message string                `json:"message"`
	Fields  map[string]any        `json:"fields"`
}

type signalCoreImpl struct {
	level  zap.Level
	enc    *zap.MapObjectEncoder
	signal *signals.Signal[signalLogMessage]
}

func newSignalCore(level zap.Level, signal *signals.Signal[signalLogMessage]) *signalCoreImpl {
	return &signalCoreImpl{
		level:  level,
		enc:    zap.NewMapObjectEncoder(),
		signal: signal,
	}
}

func (c signalCoreImpl) Enabled(l zap.Level) bool { return c.level.Enabled(l) }
func (c signalCoreImpl) With(fields []zap.Field) zap.Core {
	return signalCoreImpl{
		level:  c.level,
		enc:    cloneMapObjectEncoder(c.enc, fields),
		signal: c.signal,
	}
}
func (c signalCoreImpl) Check(entry zap.Entry, ce *zap.CheckedEntry) *zap.CheckedEntry {
	if c.Enabled(entry.Level) {
		return ce.AddCore(entry, c)
	}
	return ce
}
func (c signalCoreImpl) Write(ent zap.Entry, fields []zap.Field) error {
	var level signalLogMessageLevel
	{
		switch ent.Level {
		case zap.DebugLevel:
			level = signalLogMessageDebug
		case zap.InfoLevel:
			level = signalLogMessageInfo
		case zap.WarnLevel:
			level = signalLogMessageWarn
		default:
			// Any error above WARN is error-like in Zap
			level = signalLogMessageError
		}
	}

	newEnc := cloneMapObjectEncoder(c.enc, fields)

	msg := signalLogMessage{
		Level:   level,
		Time:    ent.Time,
		Message: ent.Message,
		Fields:  newEnc.Fields,
	}

	c.signal.EmitBG(msg)
	return nil
}
func (c signalCoreImpl) Sync() error { return nil }

func cloneMapObjectEncoder(enc *zap.MapObjectEncoder, fields []zap.Field) *zap.MapObjectEncoder {
	newEnc := zap.NewMapObjectEncoder()
	for k, v := range enc.Fields {
		newEnc.Fields[k] = v
	}

	for _, field := range fields {
		field.AddTo(newEnc)
	}

	return newEnc
}

var LogMessagesSignal = signals.New[signalLogMessage]()
