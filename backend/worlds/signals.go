package worlds

import (
	"github.com/maniartech/signals"
	"juraji.nl/chat-quest/sse"
)

var WorldCreatedSignal = signals.New[*World]()
var WorldUpdatedSignal = signals.New[*World]()
var WorldDeletedSignal = signals.New[int]()

var ChatPreferencesUpdatedSignal = signals.New[*ChatPreferences]()

func init() {
	sse.RegisterSseSourceSignal("WorldCreated", WorldCreatedSignal)
	sse.RegisterSseSourceSignal("WorldUpdated", WorldUpdatedSignal)
	sse.RegisterSseSourceSignal("WorldDeleted", WorldDeletedSignal)
}
