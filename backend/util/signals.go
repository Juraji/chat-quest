package util

import (
	"context"
	"github.com/maniartech/signals"
)

func EmitOnSuccess[T any](signal signals.Signal[T], value T, cancelOnErr error) {
	if cancelOnErr == nil {
		signal.Emit(context.TODO(), value)
	}
}

func EmitAllOnSuccess[T any](signal signals.Signal[T], value []T, cancelOnErr error) {
	if cancelOnErr == nil {
		for _, value := range value {
			signal.Emit(context.TODO(), value)
		}
	}
}
