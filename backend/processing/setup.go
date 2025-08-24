package processing

import (
	"juraji.nl/chat-quest/core/util/signals"
	chatsessions "juraji.nl/chat-quest/model/chat-sessions"
	"juraji.nl/chat-quest/processing/chat_init"
	"juraji.nl/chat-quest/processing/chat_response"
	"juraji.nl/chat-quest/processing/memory_generation"
	"time"
)

func SetupProcessing() {
	// Chat init
	chatsessions.ChatSessionCreatedSignal.AddListener(
		"CreateChatSessionGreetings", chat_init.CreateChatSessionGreetings)

	// Chat response
	chatsessions.ChatMessageCreatedSignal.AddListener(
		"GenerateResponseForMessage", chat_response.GenerateResponseForMessage)
	chatsessions.ChatParticipantResponseRequestedSignal.AddListener(
		"GenerateResponseForParticipant", chat_response.GenerateResponseForParticipant)

	// Memory generation
	chatsessions.ChatMessageCreatedSignal.AddListener(
		"GenerateMemoriesOnNewMessage",
		signals.DebounceListener(1*time.Second, memory_generation.GenerateMemories))
	chatsessions.ChatMessageUpdatedSignal.AddListener(
		"GenerateMemoriesOnUpdatedMessage",
		signals.DebounceListener(2*time.Second, memory_generation.GenerateMemories))
}
