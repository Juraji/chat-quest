package characters

import (
	"juraji.nl/chat-quest/core/sse"
	"juraji.nl/chat-quest/core/util/signals"
)

var CharacterCreatedSignal = signals.New[*Character]()
var CharacterUpdatedSignal = signals.New[*Character]()
var CharacterDeletedSignal = signals.New[int]()

var CharacterTagAddedSignal = signals.New[[]int]()
var CharacterTagRemovedSignal = signals.New[[]int]()

func init() {
	sse.RegisterOnSSE("CharacterCreated", CharacterCreatedSignal)
	sse.RegisterOnSSE("CharacterUpdated", CharacterUpdatedSignal)
	sse.RegisterOnSSE("CharacterDeleted", CharacterDeletedSignal)
	sse.RegisterOnSSE("CharacterTagAdded", CharacterTagAddedSignal)
	sse.RegisterOnSSE("CharacterTagRemoved", CharacterTagRemovedSignal)
}
