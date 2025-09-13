package processing

import (
	cs "juraji.nl/chat-quest/model/chat-sessions"
	m "juraji.nl/chat-quest/model/memories"
	p "juraji.nl/chat-quest/model/preferences"
)

func SetupProcessing() {
	// Chat greetings
	cs.ChatParticipantAddedSignal.AddListener(
		"GreetOnParticipantAdded", GreetOnParticipantAdded)

	// Chat response
	cs.ChatMessageCreatedSignal.AddListener(
		"GenerateResponse", GenerateResponseByMessageCreated)
	cs.ChatParticipantResponseRequestedSignal.AddListener(
		"GenerateResponse", GenerateResponseByParticipantTrigger)

	// Memory generation
	cs.ChatMessageCreatedSignal.AddListener(
		"GenerateMemories", GenerateMemories)
	cs.ChatMessageUpdatedSignal.AddListener(
		"GenerateMemories", GenerateMemories)
	m.MemoryGenerationForMessageRequestedSignal.AddListener(
		"GenerateMemoriesForMessageID", GenerateMemoriesForMessageID)
	m.MemoryCreatedSignal.AddListener(
		"GenerateMemoryEmbeddings", GenerateEmbeddings)
	m.MemoryUpdatedSignal.AddListener(
		"GenerateMemoryEmbeddings", GenerateEmbeddings)
	p.PreferencesUpdatedSignal.AddListener(
		"RegenerateMemoryEmbeddings", RegenerateEmbeddingsOnPrefsUpdate)
}
