package preferences

import (
	"juraji.nl/chat-quest/core/sse"
	"juraji.nl/chat-quest/core/util/signals"
)

var PreferencesUpdatedSignal = signals.New[*Preferences]()

func init() {
	sse.RegisterOnSSE("PreferencesUpdated", PreferencesUpdatedSignal)
}
