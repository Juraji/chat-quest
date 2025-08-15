package memories

import (
	"github.com/maniartech/signals"
	"juraji.nl/chat-quest/core/sse"
)

var MemoryCreatedSignal = signals.New[*Memory]()
var MemoryUpdatedSignal = signals.New[*Memory]()
var MemoryDeletedSignal = signals.New[int]()

var MemoryPreferencesUpdatedSignal = signals.New[*MemoryPreferences]()

func init() {
	sse.RegisterOnSSE("MemoryCreated", MemoryCreatedSignal)
	sse.RegisterOnSSE("MemoryUpdated", MemoryUpdatedSignal)
	sse.RegisterOnSSE("MemoryDeleted", MemoryDeletedSignal)
	sse.RegisterOnSSE("MemoryPreferencesUpdated", MemoryPreferencesUpdatedSignal)
}
