package util

import (
	"github.com/maniartech/signals"
	"juraji.nl/chat-quest/core"
)

func EmitAll[T any](
	cq *core.ChatQuestContext,
	signal signals.Signal[T],
	values []T,
) {
	for _, value := range values {
		signal.Emit(cq.Context(), value)
	}
}

func EmitAllNonNil[T any](
	cq *core.ChatQuestContext,
	signal signals.Signal[T],
	values []*T,
) {
	for _, value := range values {
		if value != nil {
			signal.Emit(cq.Context(), *value)
		}
	}
}
