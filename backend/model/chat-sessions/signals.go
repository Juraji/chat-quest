package chat_sessions

import (
	"github.com/maniartech/signals"
	"juraji.nl/chat-quest/core/sse"
)

var ChatSessionCreatedSignal = signals.New[*ChatSession]()
var ChatSessionUpdatedSignal = signals.New[*ChatSession]()
var ChatSessionDeletedSignal = signals.New[int]()

var ChatMessageCreatedSignal = signals.New[*ChatMessage]()
var ChatMessageUpdatedSignal = signals.New[*ChatMessage]()
var ChatMessageDeletedSignal = signals.New[int]()

var ChatParticipantAddedSignal = signals.New[*ChatParticipant]()
var ChatParticipantRemovedSignal = signals.New[*ChatParticipant]()

func init() {
	sse.RegisterSseSourceSignal("ChatSessionCreated", ChatSessionCreatedSignal)
	sse.RegisterSseSourceSignal("ChatSessionUpdated", ChatSessionUpdatedSignal)
	sse.RegisterSseSourceSignal("ChatSessionDeleted", ChatSessionDeletedSignal)
	sse.RegisterSseSourceSignal("ChatMessageCreated", ChatMessageCreatedSignal)
	sse.RegisterSseSourceSignal("ChatMessageUpdated", ChatMessageUpdatedSignal)
	sse.RegisterSseSourceSignal("ChatMessageDeleted", ChatMessageDeletedSignal)
	sse.RegisterSseSourceSignal("ChatParticipantAdded", ChatParticipantAddedSignal)
	sse.RegisterSseSourceSignal("ChatParticipantRemoved", ChatParticipantRemovedSignal)
}
