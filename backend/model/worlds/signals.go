package worlds

import (
	"juraji.nl/chat-quest/core/sse"
	"juraji.nl/chat-quest/core/util/signals"
)

var WorldCreatedSignal = signals.New[*World]()
var WorldUpdatedSignal = signals.New[*World]()
var WorldDeletedSignal = signals.New[int]()

var ChatPreferencesUpdatedSignal = signals.New[*ChatPreferences]()

func init() {
	sse.RegisterOnSSE("WorldCreated", WorldCreatedSignal)
	sse.RegisterOnSSE("WorldUpdated", WorldUpdatedSignal)
	sse.RegisterOnSSE("WorldDeleted", WorldDeletedSignal)
	sse.RegisterOnSSE("ChatPreferencesUpdated", ChatPreferencesUpdatedSignal)
}
