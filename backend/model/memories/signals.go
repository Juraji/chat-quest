package memories

import (
	"juraji.nl/chat-quest/core/sse"
	"juraji.nl/chat-quest/core/util/signals"
)

var MemoryCreatedSignal = signals.New[*Memory]()
var MemoryUpdatedSignal = signals.New[*Memory]()
var MemoryDeletedSignal = signals.New[int]()

func init() {
	sse.RegisterOnSSE("MemoryCreated", MemoryCreatedSignal)
	sse.RegisterOnSSE("MemoryUpdated", MemoryUpdatedSignal)
	sse.RegisterOnSSE("MemoryDeleted", MemoryDeletedSignal)
}
