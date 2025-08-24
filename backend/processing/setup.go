package processing

import (
	chatsessions "juraji.nl/chat-quest/model/chat-sessions"
	"juraji.nl/chat-quest/model/memories"
	"juraji.nl/chat-quest/processing/chat_init"
	"juraji.nl/chat-quest/processing/chat_response"
	"juraji.nl/chat-quest/processing/memory_generation"
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
		"GenerateMemoriesOnNewMessage", memory_generation.GenerateMemories)
	chatsessions.ChatMessageUpdatedSignal.AddListener(
		"GenerateMemoriesOnUpdatedMessage", memory_generation.GenerateMemories)
	memories.MemoryCreatedSignal.AddListener(
		"EmbeddingsForNewMemory", memory_generation.GenerateEmbeddings)
	memories.MemoryUpdatedSignal.AddListener(
		"EmbeddingsForExistingMemory", memory_generation.GenerateEmbeddings)
}
