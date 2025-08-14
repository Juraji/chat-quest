package sse

import (
	"github.com/maniartech/signals"
)

var sseSourceSignals []sourceSignal

func RegisterSseSourceSignal[T any](name string, s signals.Signal[T]) {
	source := sourceSignal{
		sourceName: name,
		signal:     &signalWrapper[T]{signal: s},
	}

	sseSourceSignals = append(sseSourceSignals, source)
	//log.Get().Printf("Registered SSE signal for event '%s' with type [%v]", name, reflect.TypeOf(*new(T)))
}
