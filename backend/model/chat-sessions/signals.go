package chat_sessions

import (
	"juraji.nl/chat-quest/core/sse"
	"juraji.nl/chat-quest/core/util/signals"
)

var ChatSessionCreatedSignal = signals.New[*ChatSession]()
var ChatSessionUpdatedSignal = signals.New[*ChatSession]()
var ChatSessionDeletedSignal = signals.New[int]()

var ChatMessageCreatedSignal = signals.New[*ChatMessage]()
var ChatMessageUpdatedSignal = signals.New[*ChatMessage]()
var ChatMessageDeletedSignal = signals.New[int]()

var ChatParticipantAddedSignal = signals.New[*ChatParticipant]()
var ChatParticipantRemovedSignal = signals.New[*ChatParticipant]()
var ChatParticipantResponseRequestedSignal = signals.New[*ChatParticipant]()

func init() {
	sse.RegisterOnSSE("ChatSessionCreated", ChatSessionCreatedSignal)
	sse.RegisterOnSSE("ChatSessionUpdated", ChatSessionUpdatedSignal)
	sse.RegisterOnSSE("ChatSessionDeleted", ChatSessionDeletedSignal)
	sse.RegisterOnSSE("ChatMessageCreated", ChatMessageCreatedSignal)
	sse.RegisterOnSSE("ChatMessageUpdated", ChatMessageUpdatedSignal)
	sse.RegisterOnSSE("ChatMessageDeleted", ChatMessageDeletedSignal)
	sse.RegisterOnSSE("ChatParticipantAdded", ChatParticipantAddedSignal)
	sse.RegisterOnSSE("ChatParticipantRemoved", ChatParticipantRemovedSignal)
}
