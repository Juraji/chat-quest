package chat_sessions

import (
	"github.com/maniartech/signals"
	"juraji.nl/chat-quest/sse"
)

var ChatSessionCreatedSignal = signals.New[*ChatSession]()
var ChatSessionUpdatedSignal = signals.New[*ChatSession]()
var ChatSessionDeletedSignal = signals.New[int64]()

var ChatMessageCreatedSignal = signals.New[*ChatMessage]()
var ChatMessageUpdatedSignal = signals.New[*ChatMessage]()
var ChatMessageDeletedSignal = signals.New[int64]()

func init() {
	sse.RegisterSseSourceSignal("ChatSessionCreated", ChatSessionCreatedSignal)
	sse.RegisterSseSourceSignal("ChatSessionUpdated", ChatSessionUpdatedSignal)
	sse.RegisterSseSourceSignal("ChatSessionDeleted", ChatSessionDeletedSignal)
	sse.RegisterSseSourceSignal("ChatMessageCreated", ChatMessageCreatedSignal)
	sse.RegisterSseSourceSignal("ChatMessageUpdated", ChatMessageUpdatedSignal)
	sse.RegisterSseSourceSignal("ChatMessageDeleted", ChatMessageDeletedSignal)
}
