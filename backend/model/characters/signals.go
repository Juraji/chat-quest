package characters

import (
	"github.com/maniartech/signals"
	"juraji.nl/chat-quest/core/sse"
)

var CharacterCreatedSignal = signals.New[*Character]()
var CharacterUpdatedSignal = signals.New[*Character]()
var CharacterDeletedSignal = signals.New[int]()

func init() {
	sse.RegisterOnSSE("CharacterCreated", CharacterCreatedSignal)
	sse.RegisterOnSSE("CharacterUpdated", CharacterUpdatedSignal)
	sse.RegisterOnSSE("CharacterDeleted", CharacterDeletedSignal)
}
