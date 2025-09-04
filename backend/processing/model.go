package processing

import (
	c "juraji.nl/chat-quest/model/characters"
	m "juraji.nl/chat-quest/model/memories"
)

type characterTemplateVars struct {
	Character *c.Character
}

type responseInstructionVars struct {
	MessageIndex         int
	Message              string
	IsTriggeredByMessage bool

	// Responding character info
	Character        *c.Character
	DialogueExamples []string

	// Persona
	Persona *c.Character

	// Session details
	IsSingleCharacter   bool
	OtherParticipants   []c.Character
	WorldDescription    string
	ScenarioDescription string

	// Memories
	Memories []m.Memory
}

type memoryInstructionVars struct {
	Participants []c.Character
}
