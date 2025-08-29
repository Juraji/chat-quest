package processing

import (
	cs "juraji.nl/chat-quest/model/chat-sessions"
	m "juraji.nl/chat-quest/model/memories"
	p "juraji.nl/chat-quest/model/preferences"
	"juraji.nl/chat-quest/processing/chat_init"
	"juraji.nl/chat-quest/processing/chat_responses"
	"juraji.nl/chat-quest/processing/memory_generation"
)

func SetupProcessing() {
	// Chat init
	// TODO: Run greeting when character is being added, instead of at session creation
	cs.ChatSessionCreatedSignal.AddListener(
		"ProcessingGreetings", chat_init.CreateChatSessionGreetings)

	// Chat response
	cs.ChatMessageCreatedSignal.AddListener(
		"GenerateResponse", chat_responses.GenerateResponseByMessageCreated)
	cs.ChatParticipantResponseRequestedSignal.AddListener(
		"GenerateResponse", chat_responses.GenerateResponseByParticipantTrigger)

	// Memory generation
	cs.ChatMessageCreatedSignal.AddListener(
		"GenerateMemories", memory_generation.GenerateMemories)
	cs.ChatMessageUpdatedSignal.AddListener(
		"GenerateMemories", memory_generation.GenerateMemories)
	m.MemoryCreatedSignal.AddListener(
		"GenerateMemoryEmbeddings", memory_generation.GenerateEmbeddings)
	m.MemoryUpdatedSignal.AddListener(
		"GenerateMemoryEmbeddings", memory_generation.GenerateEmbeddings)
	p.PreferencesUpdatedSignal.AddListener(
		"RegenerateMemoryEmbeddings", memory_generation.RegenerateEmbeddingsOnPrefsUpdate)
}
