package chat_response

import (
	"fmt"
	t "juraji.nl/chat-quest/core/util/template_utils"
	c "juraji.nl/chat-quest/model/characters"
	s "juraji.nl/chat-quest/model/chat-sessions"
	"juraji.nl/chat-quest/model/scenarios"
	"juraji.nl/chat-quest/model/worlds"
)

type instructionTemplateVars struct {
	MessageIndex int
	Message      string

	// Responding character info
	Character        *c.Character
	DialogueExamples *t.LazyTemplateSlice[string]

	// Session details
	OtherParticipants   *t.LazyTemplateSlice[c.Character]
	WorldDescription    *t.LazyTemplateVar[string]
	ScenarioDescription *t.LazyTemplateVar[string]
}

func newInstructionTemplateVars(
	session *s.ChatSession,
	triggerMessage *s.ChatMessage,
	chatHistory []s.ChatMessage,
	characterId int,
) (*instructionTemplateVars, error) {
	// We handle the main character first, as it is the only thing that can error "now".
	var character *c.Character
	{
		char, err := c.CharacterById(characterId)
		if err != nil {
			return nil, err
		}
		err = applyCharacterTemplates(char)
		if err != nil {
			return nil, err
		}
		character = char
	}

	dialogueExamplesVar := t.NewLazyTemplateSlice(func() ([]string, error) {
		return c.DialogueExamplesByCharacterId(characterId)
	})

	otherParticipantsVar := t.NewLazyTemplateSlice(func() ([]c.Character, error) {
		participants, err := s.GetParticipants(session.ID)
		if err != nil {
			return nil, err
		}

		for _, participant := range participants {
			err := applyCharacterTemplates(&participant)
			if err != nil {
				return nil, err
			}
		}

		return participants, nil
	})

	worldDescriptionVar := t.NewLazyTemplateVar(func() (string, error) {
		w, err := worlds.WorldById(session.WorldID)
		if err != nil {
			return "", err
		}

		if w.Description == nil {
			return "", nil
		} else {
			return *w.Description, nil
		}
	})

	scenarioDescriptionVar := t.NewLazyTemplateVar(func() (string, error) {
		if session.ScenarioID == nil {
			return "", nil
		}
		scenario, err := scenarios.ScenarioById(*session.ScenarioID)
		if err != nil {
			return "", err
		}

		return scenario.Description, nil
	})

	return &instructionTemplateVars{
		MessageIndex:        len(chatHistory),
		Message:             triggerMessage.Content,
		Character:           character,
		DialogueExamples:    dialogueExamplesVar,
		OtherParticipants:   otherParticipantsVar,
		WorldDescription:    worldDescriptionVar,
		ScenarioDescription: scenarioDescriptionVar,
	}, nil
}

func applyCharacterTemplates(char *c.Character) error {
	characterVars := struct{ *c.Character }{char}

	fieldsToProcess := []*string{
		char.Appearance,
		char.Personality,
		char.History,
	}

	for _, fieldPtr := range fieldsToProcess {
		if fieldPtr == nil || !t.HasTemplateVars(*fieldPtr) {
			continue
		}

		tpl, err := t.NewTemplate("", *fieldPtr, nil)
		if err != nil {
			return fmt.Errorf("failed to create template for character field: %w", err)
		}

		*fieldPtr = t.WriteToString(tpl, characterVars)
	}

	return nil
}
