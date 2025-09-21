package processing

import (
	"sync"

	c "juraji.nl/chat-quest/model/characters"
	cs "juraji.nl/chat-quest/model/chat-sessions"
	p "juraji.nl/chat-quest/model/preferences"
	sc "juraji.nl/chat-quest/model/scenarios"
	w "juraji.nl/chat-quest/model/worlds"
)

type ChatInstructionVars interface {
	IsTriggeredByMessage() bool
	IsFirstMessage() bool
	CurrentMessageIndex() int
	MessageText() string

	Character() (TemplateCharacter, error)
	Persona() (TemplateCharacter, error)

	OtherParticipants() ([]TemplateCharacter, error)

	World() (string, error)
	Scenario() (string, error)
	CurrentTimeOfDay() (cs.TimeOfDay, error)
}

type chatInstructionVarsImpl struct {
	triggerMessage      *cs.ChatMessage
	currentMessageIndex int
	timeOfDay           *cs.TimeOfDay
	character           func() (TemplateCharacter, error)
	persona             func() (TemplateCharacter, error)
	otherParticipants   func() ([]TemplateCharacter, error)
	world               func() (string, error)
	scenario            func() (string, error)
}

func (c *chatInstructionVarsImpl) IsTriggeredByMessage() bool {
	return c.triggerMessage != nil
}
func (c *chatInstructionVarsImpl) CurrentMessageIndex() int {
	return c.currentMessageIndex
}
func (c *chatInstructionVarsImpl) IsFirstMessage() bool {
	return c.currentMessageIndex == 0
}
func (c *chatInstructionVarsImpl) MessageText() string {
	if c.triggerMessage == nil {
		return ""
	}
	return c.triggerMessage.Content
}
func (c *chatInstructionVarsImpl) Character() (TemplateCharacter, error) {
	return c.character()
}
func (c *chatInstructionVarsImpl) Persona() (TemplateCharacter, error) {
	return c.persona()
}
func (c *chatInstructionVarsImpl) OtherParticipants() ([]TemplateCharacter, error) {
	return c.otherParticipants()
}
func (c *chatInstructionVarsImpl) World() (string, error) {
	return c.world()
}
func (c *chatInstructionVarsImpl) Scenario() (string, error) {
	return c.scenario()
}
func (c *chatInstructionVarsImpl) CurrentTimeOfDay() (cs.TimeOfDay, error) {
	if c.timeOfDay == nil {
		return "", nil
	}
	return *c.timeOfDay, nil
}

func NewChatInstructionVars(
	chatHistory []cs.ChatMessage,
	session *cs.ChatSession,
	prefs *p.Preferences,
	triggerMessage *cs.ChatMessage,
	characterId int,
) ChatInstructionVars {
	fullHistory := chatHistory
	if triggerMessage != nil {
		fullHistory = append(chatHistory, *triggerMessage)
	}

	return &chatInstructionVarsImpl{
		triggerMessage:      triggerMessage,
		currentMessageIndex: len(fullHistory),
		timeOfDay:           session.CurrentTimeOfDay,
		character: sync.OnceValues(func() (TemplateCharacter, error) {
			character, err := c.CharacterById(characterId)
			if err != nil {
				return nil, err
			}
			return NewTemplateCharacter(character, prefs, session, fullHistory), nil
		}),
		persona: sync.OnceValues(func() (TemplateCharacter, error) {
			world, err := w.WorldById(session.WorldID)
			if err != nil {
				return nil, err
			}
			if world.PersonaID == nil {
				return nil, nil
			}
			character, err := c.CharacterById(*world.PersonaID)
			if err != nil {
				return nil, err
			}
			return NewTemplateCharacter(character, prefs, session, fullHistory), nil
		}),
		otherParticipants: sync.OnceValues(func() ([]TemplateCharacter, error) {
			allParticipants, err := cs.GetAllParticipantsAsCharacters(session.ID)
			if err != nil {
				return nil, err
			}
			templateVars := make([]TemplateCharacter, 0, len(allParticipants))
			for _, participant := range allParticipants {
				if participant.ID != characterId {
					templateVars = append(templateVars, NewTemplateCharacter(&participant, prefs, session, fullHistory))
				}
			}
			return templateVars, nil
		}),
		world: sync.OnceValues(func() (string, error) {
			world, err := w.WorldById(session.WorldID)
			if err != nil {
				return "", err
			}
			if world.Description == nil {
				return "", nil
			}
			return *world.Description, nil
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
	}
}
