package util

import (
	"context"
	"github.com/maniartech/signals"
)

func EmitOnSuccess[T any](signal signals.Signal[T], value T, cancelOnErr error) {
	if cancelOnErr != nil {
		return
	}

	signal.Emit(context.TODO(), value)
}

func EmitAllNonNilOnSuccess[T any](signal signals.Signal[T], values []*T, cancelOnErr error) {
	if cancelOnErr != nil {
		return
	}

	for _, value := range values {
		if value != nil {
			signal.Emit(context.TODO(), *value)
		}
	}
}

func EmitAllOnSuccess[T any](signal signals.Signal[T], values []T, cancelOnErr error) {
	if cancelOnErr != nil {
		return
	}

	for _, value := range values {
		signal.Emit(context.TODO(), value)
	}
}
