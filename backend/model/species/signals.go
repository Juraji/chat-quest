package species

import (
	"juraji.nl/chat-quest/core/sse"
	"juraji.nl/chat-quest/core/util/signals"
)

var SpeciesCreatedSignal = signals.New[*Species]()
var SpeciesUpdatedSignal = signals.New[*Species]()
var SpeciesDeletedSignal = signals.New[int]()

func init() {
	sse.RegisterOnSSE("SpeciesCreated", SpeciesCreatedSignal)
	sse.RegisterOnSSE("SpeciesUpdated", SpeciesUpdatedSignal)
	sse.RegisterOnSSE("SpeciesDeleted", SpeciesDeletedSignal)
}
