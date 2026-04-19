package processing

import (
	"sync"

	c "juraji.nl/chat-quest/model/characters"
	cs "juraji.nl/chat-quest/model/chat-sessions"
	p "juraji.nl/chat-quest/model/preferences"
	sc "juraji.nl/chat-quest/model/scenarios"
	sp "juraji.nl/chat-quest/model/species"
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
	SpeciesPresent() ([]TemplateSpecies, error)
	CurrentTimeOfDay() *cs.TimeOfDay
	CurrentTimeOfDayFmtEN() string
	ChatNotes() string
}

type chatInstructionVarsImpl struct {
	triggerMessage      *cs.ChatMessage
	currentMessageIndex int
	timeOfDay           *cs.TimeOfDay
	chatNotes           *string
	character           func() (TemplateCharacter, error)
	persona             func() (TemplateCharacter, error)
	otherParticipants   func() ([]TemplateCharacter, error)
	world               func() (string, error)
	scenario            func() (string, error)
	presentSpecies      func() ([]TemplateSpecies, error)
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
func (c *chatInstructionVarsImpl) SpeciesPresent() ([]TemplateSpecies, error) {
	return c.presentSpecies()
}
func (c *chatInstructionVarsImpl) CurrentTimeOfDay() *cs.TimeOfDay {
	return c.timeOfDay
}
func (c *chatInstructionVarsImpl) CurrentTimeOfDayFmtEN() string { return c.timeOfDay.HumanFmtEn() }
func (c *chatInstructionVarsImpl) ChatNotes() string {
	if c.chatNotes == nil {
		return ""
	}
	return *c.chatNotes
}

func NewChatInstructionVars(
	session *cs.ChatSession,
	prefs *p.Preferences,
	chatHistory []cs.ChatMessage,
	triggerMessage *cs.ChatMessage,
	sessionMessageCount int,
	currentCharacterId int,
) ChatInstructionVars {
	fullHistory := chatHistory
	if triggerMessage != nil {
		fullHistory = append(chatHistory, *triggerMessage)
	}

	return &chatInstructionVarsImpl{
		triggerMessage:      triggerMessage,
		currentMessageIndex: sessionMessageCount,
		timeOfDay:           session.CurrentTimeOfDay,
		chatNotes:           session.ChatNotes,
		character: sync.OnceValues(func() (TemplateCharacter, error) {
			character, err := c.CharacterById(currentCharacterId)
			if err != nil {
				return nil, err
			}
			return NewTemplateCharacter(character, prefs, session, fullHistory), nil
		}),
		persona: sync.OnceValues(func() (TemplateCharacter, error) {
			if session.PersonaID == nil {
				return nil, nil
			}
			character, err := c.CharacterById(*session.PersonaID)
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
				if participant.ID != currentCharacterId {
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
		presentSpecies: sync.OnceValues(func() ([]TemplateSpecies, error) {
			species, err := sp.GetSpeciesPresentInSession(session.ID)
			if err != nil || len(species) == 0 {
				return nil, err
			}

			var templateSpecies []TemplateSpecies
			for _, s := range species {
				templateSpecies = append(templateSpecies, NewTemplateSpecies(&s))
			}

			return templateSpecies, nil
		}),
	}
}
