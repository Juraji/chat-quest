package processing

import (
	"sync"
	"time"

	c "juraji.nl/chat-quest/model/characters"
	cs "juraji.nl/chat-quest/model/chat-sessions"
	sc "juraji.nl/chat-quest/model/scenarios"
)

type MemoryInstructionVars interface {
	Participants() ([]TemplateCharacter, error)
	Persona() (SparseTemplateCharacter, error)
	Scenario() (string, error)
	ChatNotes() string
	CurrentTimeOfDay() *cs.TimeOfDay
	CurrentTimeOfDayFmtEN() string
}

type memoryInstructionVarsImpl struct {
	participants func() ([]TemplateCharacter, error)
	persona      func() (SparseTemplateCharacter, error)
	scenario     func() (string, error)
	chatNotes    *string
	timeOfDay    *cs.TimeOfDay
}

func (m memoryInstructionVarsImpl) Participants() ([]TemplateCharacter, error) {
	return m.participants()
}
func (m memoryInstructionVarsImpl) Persona() (SparseTemplateCharacter, error) { return m.persona() }
func (m memoryInstructionVarsImpl) Scenario() (string, error)                 { return m.scenario() }
func (m memoryInstructionVarsImpl) ChatNotes() string {
	if m.chatNotes == nil {
		return ""
	}
	return *m.chatNotes
}
func (m memoryInstructionVarsImpl) CurrentTimeOfDay() *cs.TimeOfDay { return m.timeOfDay }
func (m memoryInstructionVarsImpl) CurrentTimeOfDayFmtEN() string   { return m.timeOfDay.HumanFmtEn() }

func NewMemoryInstructionVars(session *cs.ChatSession, before time.Time) MemoryInstructionVars {
	return &memoryInstructionVarsImpl{
		participants: sync.OnceValues(func() ([]TemplateCharacter, error) {
			allParticipants, err := cs.GetAllParticipantsAsCharactersBefore(session.ID, before)
			if err != nil {
				return nil, err
			}
			templateVars := make([]TemplateCharacter, len(allParticipants))
			for i, participant := range allParticipants {
				templateVars[i] = NewTemplateCharacter(&participant, nil, nil)
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
		scenario: sync.OnceValues(func() (string, error) {
			if session.ScenarioID == nil {
				return "", nil
			}
			scenario, err := sc.ScenarioById(*session.ScenarioID)
			if err != nil {
				return "", err
			}

			return scenario.Description, nil
		}),
		chatNotes: session.ChatNotes,
		timeOfDay: session.CurrentTimeOfDay,
	}
}
