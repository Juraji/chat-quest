package processing

import (
	"sync"
	"time"

	c "juraji.nl/chat-quest/model/characters"
	cs "juraji.nl/chat-quest/model/chat-sessions"
)

type MemoryInstructionVars interface {
	Participants() ([]TemplateCharacter, error)
	Persona() (SparseTemplateCharacter, error)
}

type memoryInstructionVarsImpl struct {
	participants func() ([]TemplateCharacter, error)
	persona      func() (SparseTemplateCharacter, error)
}

func (m memoryInstructionVarsImpl) Participants() ([]TemplateCharacter, error) {
	return m.participants()
}
func (m memoryInstructionVarsImpl) Persona() (SparseTemplateCharacter, error) {
	return m.persona()
}

func NewMemoryInstructionVars(session *cs.ChatSession, before time.Time) MemoryInstructionVars {
	return &memoryInstructionVarsImpl{
		participants: sync.OnceValues(func() ([]TemplateCharacter, error) {
			allParticipants, err := cs.GetAllParticipantsAsCharactersBefore(session.ID, before)
			if err != nil {
				return nil, err
			}
			templateVars := make([]TemplateCharacter, len(allParticipants))
			for i, participant := range allParticipants {
				templateVars[i] = NewTemplateCharacter(&participant, nil, nil, nil)
			}
			return templateVars, nil
		}),
		persona: sync.OnceValues(func() (SparseTemplateCharacter, error) {
			if session.PersonaID == nil {
				return nil, nil
			}
			character, err := c.CharacterById(*session.PersonaID)
			if err != nil {
				return nil, err
			}
			return NewSparseTemplateCharacter(character), nil
		}),
	}
}
