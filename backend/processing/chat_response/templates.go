package chat_response

import (
	"fmt"
	"juraji.nl/chat-quest/core/util"
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
	DialogueExamples []string

	// Session details
	OtherParticipants   []c.Character
	WorldDescription    string
	ScenarioDescription string
}

func newInstructionTemplateVars(
	session *s.ChatSession,
	triggerMessage *s.ChatMessage,
	chatHistory []s.ChatMessage,
	characterId int,
) (*instructionTemplateVars, error) {
	vars := instructionTemplateVars{
		MessageIndex: len(chatHistory),
		Message:      triggerMessage.Content,
	}
	errChan := make(chan error, 5)

	// Fetch main character
	go func() {
		char, err := c.CharacterById(characterId)
		if err != nil {
			errChan <- fmt.Errorf("error getting character by id: %v", err)
			return
		}
		err = applyCharacterTemplates(char)
		if err != nil {
			errChan <- fmt.Errorf("error applying character templates: %v", err)
			return
		}

		vars.Character = char
		errChan <- nil
	}()

	// Fetch dialogue examples
	go func() {
		de, err := c.DialogueExamplesByCharacterId(characterId)
		if err != nil {
			errChan <- fmt.Errorf("error getting dialogue examples by character: %v", err)
			return
		}

		vars.DialogueExamples = de
		errChan <- nil
	}()

	// Fetch other participants
	go func() {
		ops, err := s.GetParticipants(session.ID)
		if err != nil {
			errChan <- fmt.Errorf("error getting participants: %v", err)
			return
		}

		vars.OtherParticipants = ops
		errChan <- nil
	}()

	// Fetch world description
	go func() {
		w, err := worlds.WorldById(session.WorldID)
		if err != nil {
			errChan <- fmt.Errorf("error getting world: %v", err)
			return
		}

		if w.Description != nil {
			vars.WorldDescription = *w.Description
		}
		errChan <- nil
	}()

	// Fetch scenario description
	go func() {
		if session.ScenarioID == nil {
			errChan <- nil
			return
		}

		scenario, err := scenarios.ScenarioById(*session.ScenarioID)
		if err != nil {
			errChan <- fmt.Errorf("error getting scenario: %v", err)
			return
		}

		vars.ScenarioDescription = scenario.Description
		errChan <- nil
	}()

	for i := 0; i < 5; i++ {
		if err := <-errChan; err != nil {
			return nil, err
		}
	}

	return &vars, nil
}

func applyCharacterTemplates(char *c.Character) error {
	characterVars := struct{ *c.Character }{char}

	fieldsToProcess := []*string{
		char.Appearance,
		char.Personality,
		char.History,
	}

	for _, fieldPtr := range fieldsToProcess {
		if fieldPtr == nil || !util.HasTemplateVars(*fieldPtr) {
			continue
		}

		tpl, err := util.NewTextTemplate(char.Name, *fieldPtr)
		if err != nil {
			return fmt.Errorf("failed to create template for character field: %w", err)
		}

		*fieldPtr = util.WriteToString(tpl, characterVars)
	}

	return nil
}
