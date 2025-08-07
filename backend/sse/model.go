package sse

import (
	"context"
	"github.com/maniartech/signals"
)

type messageBody struct {
	Source  string `json:"source"`
	Payload any    `json:"payload"`
}

type sourceSignal struct {
	sourceName string
	signal     anySignalInterface
}

type anySignalInterface interface {
	AddListener(func(context.Context, any), string)
	RemoveListener(string)
}

type signalWrapper[T any] struct {
	signal signals.Signal[T]
}

func (sw *signalWrapper[T]) AddListener(listener func(context.Context, any), key string) {
	// Convert the any listener to a T-specific listener
	typedListener := func(ctx context.Context, payload T) {
		listener(ctx, payload)
	}
	sw.signal.AddListener(typedListener, key)
}

func (sw *signalWrapper[T]) RemoveListener(key string) {
	sw.signal.RemoveListener(key)
}

func newSourceSignal[T any](name string, s signals.Signal[T]) sourceSignal {
	return sourceSignal{
		sourceName: name,
		signal:     &signalWrapper[T]{signal: s},
	}
}
