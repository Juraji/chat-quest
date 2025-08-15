package sse

import (
	"context"
	"github.com/maniartech/signals"
)

type message struct {
	Source  string `json:"source"`
	Payload any    `json:"payload"`
}

var sseCombinedSignal = signals.New[message]()

func RegisterOnSSE[T any](name string, s signals.Signal[T]) {
	s.AddListener(func(ctx context.Context, t T) {
		sseCombinedSignal.Emit(ctx, message{Source: name, Payload: t})
	})
}
