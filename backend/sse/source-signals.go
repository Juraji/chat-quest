package sse

import (
	"juraji.nl/chat-quest/characters"
	"juraji.nl/chat-quest/chat-sessions"
	"juraji.nl/chat-quest/instructions"
	"juraji.nl/chat-quest/memories"
	"juraji.nl/chat-quest/providers"
	"juraji.nl/chat-quest/scenarios"
	"juraji.nl/chat-quest/worlds"
)

var sseSourceSignals = []SourceSignal{
	// Characters
	NewSourceSignal("CharacterCreated", characters.CharacterCreatedSignal),
	NewSourceSignal("CharacterUpdated", characters.CharacterUpdatedSignal),
	NewSourceSignal("CharacterDeleted", characters.CharacterDeletedSignal),
	// Chat Sessions
	NewSourceSignal("ChatSessionCreated", chat_sessions.ChatSessionCreatedSignal),
	NewSourceSignal("ChatSessionUpdated", chat_sessions.ChatSessionUpdatedSignal),
	NewSourceSignal("ChatSessionDeleted", chat_sessions.ChatSessionDeletedSignal),
	NewSourceSignal("ChatMessageCreated", chat_sessions.ChatMessageCreatedSignal),
	NewSourceSignal("ChatMessageUpdated", chat_sessions.ChatMessageUpdatedSignal),
	NewSourceSignal("ChatMessageDeleted", chat_sessions.ChatMessageDeletedSignal),
	// Instructions
	NewSourceSignal("InstructionTemplateCreated", instructions.InstructionTemplateCreatedSignal),
	NewSourceSignal("InstructionTemplateUpdated", instructions.InstructionTemplateUpdatedSignal),
	NewSourceSignal("InstructionTemplateDeleted", instructions.InstructionTemplateDeletedSignal),
	// Memories
	NewSourceSignal("MemoryCreated", memories.MemoryCreatedSignal),
	NewSourceSignal("MemoryUpdated", memories.MemoryUpdatedSignal),
	NewSourceSignal("MemoryDeleted", memories.MemoryDeletedSignal),
	NewSourceSignal("MemoryPreferencesUpdated", memories.MemoryPreferencesUpdatedSignal),
	//Providers
	NewSourceSignal("ConnectionProfileCreated", providers.ConnectionProfileCreatedSignal),
	NewSourceSignal("ConnectionProfileUpdated", providers.ConnectionProfileUpdatedSignal),
	NewSourceSignal("ConnectionProfileDeleted", providers.ConnectionProfileDeletedSignal),
	NewSourceSignal("LlmModelCreated", providers.LlmModelCreatedSignal),
	NewSourceSignal("LlmModelUpdated", providers.LlmModelUpdatedSignal),
	NewSourceSignal("LlmModelDeleted", providers.LlmModelDeletedSignal),
	// Scenarios
	NewSourceSignal("ScenarioCreated", scenarios.ScenarioCreatedSignal),
	NewSourceSignal("ScenarioUpdated", scenarios.ScenarioUpdatedSignal),
	NewSourceSignal("ScenarioDeleted", scenarios.ScenarioDeletedSignal),
	// Worlds
	NewSourceSignal("WorldCreated", worlds.WorldCreatedSignal),
	NewSourceSignal("WorldUpdated", worlds.WorldUpdatedSignal),
	NewSourceSignal("WorldDeleted", worlds.WorldDeletedSignal),
	NewSourceSignal("ChatPreferencesUpdated", worlds.ChatPreferencesUpdatedSignal),
}
