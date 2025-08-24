package signals

import (
	"context"
	"time"
)

func DebounceListener[T any](d time.Duration, listener SignalListener[T]) SignalListener[T] {
	var timer *time.Timer
	return func(ctx context.Context, payload T) {
		if timer != nil {
			timer.Reset(d)
		} else {
			timer = time.AfterFunc(d, func() {
				listener(ctx, payload)
				timer = nil
			})
		}
	}
}
