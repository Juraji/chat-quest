package memories

import (
	"github.com/maniartech/signals"
	"juraji.nl/chat-quest/sse"
)

var MemoryCreatedSignal = signals.New[*Memory]()
var MemoryUpdatedSignal = signals.New[*Memory]()
var MemoryDeletedSignal = signals.New[int]()

var MemoryPreferencesUpdatedSignal = signals.New[*MemoryPreferences]()

func init() {
	sse.RegisterSseSourceSignal("MemoryCreated", MemoryCreatedSignal)
	sse.RegisterSseSourceSignal("MemoryUpdated", MemoryUpdatedSignal)
	sse.RegisterSseSourceSignal("MemoryDeleted", MemoryDeletedSignal)
	sse.RegisterSseSourceSignal("MemoryPreferencesUpdated", MemoryPreferencesUpdatedSignal)
}
