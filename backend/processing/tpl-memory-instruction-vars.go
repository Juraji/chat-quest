package processing

import (
	"sync"
	"time"

	cs "juraji.nl/chat-quest/model/chat-sessions"
)

type MemoryInstructionVars interface {
	Participants() ([]TemplateCharacter, error)
}

type memoryInstructionVarsImpl struct {
	participants func() ([]TemplateCharacter, error)
}

func (m memoryInstructionVarsImpl) Participants() ([]TemplateCharacter, error) {
	return m.participants()
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
	}
}
