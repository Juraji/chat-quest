package sse

import (
	"context"

	"juraji.nl/chat-quest/core/util/signals"
)

type Message struct {
	Source  string `json:"source"`
	Payload any    `json:"payload"`
}

var SseCombinedSignal = signals.New[Message]()

func RegisterOnSSE[T any](name string, s *signals.Signal[T]) {
	s.AddListener(name, func(ctx context.Context, t T) {
		SseCombinedSignal.Emit(ctx, Message{Source: name, Payload: t})
	})
}
