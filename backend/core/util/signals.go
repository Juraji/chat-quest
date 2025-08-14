package util

import (
	"context"
	"github.com/maniartech/signals"
)

func Emit[T any](
	signal signals.Signal[T],
	value T,
) {
	signal.Emit(context.Background(), value)
}

func EmitAll[T any](
	signal signals.Signal[T],
	values []T,
) {
	for _, value := range values {
		signal.Emit(context.Background(), value)
	}
}

func EmitAllNonNil[T any](
	signal signals.Signal[T],
	values []*T,
) {
	for _, value := range values {
		if value != nil {
			signal.Emit(context.Background(), *value)
		}
	}
}
