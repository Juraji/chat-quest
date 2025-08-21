package chat_response

import (
	"github.com/pkg/errors"
	"juraji.nl/chat-quest/core/util"
	c "juraji.nl/chat-quest/model/characters"
	s "juraji.nl/chat-quest/model/chat-sessions"
	"juraji.nl/chat-quest/model/scenarios"
	"juraji.nl/chat-quest/model/worlds"
)

type instructionTemplateVars struct {
	MessageIndex         int
	Message              string
	IsTriggeredByMessage bool

	// Responding character info
	Character        *c.Character
	DialogueExamples []string

	// Session details
	IsSingleCharacter   bool
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
	}

	if triggerMessage != nil {
		vars.IsTriggeredByMessage = true
		vars.Message = triggerMessage.Content
	} else {
		vars.IsTriggeredByMessage = false
	}

	errChan := make(chan error, 5)

	// Fetch main character
	go func() {
		char, ok := c.CharacterById(characterId)
		if !ok {
			errChan <- errors.New("error getting character by id")
			return
		}
		err := applyCharacterTemplates(char)
		if err != nil {
			errChan <- errors.Wrap(err, "error applying character templates")
			return
		}

		vars.Character = char
		errChan <- nil
	}()

	// Fetch dialogue examples
	go func() {
		de, ok := c.DialogueExamplesByCharacterId(characterId)
		if !ok {
			errChan <- errors.New("error getting dialogue examples by character")
			return
		}

		vars.DialogueExamples = de
		errChan <- nil
	}()

	// Fetch other participants
	go func() {
		participants, ok := s.GetParticipants(session.ID)
		if !ok {
			errChan <- errors.New("error getting participants")
			return
		}

		// Filter out the character with id characterId
		theOther := make([]c.Character, 0, len(participants))
		for _, op := range participants {
			if op.ID != characterId {
				theOther = append(theOther, op)
			}
		}

		vars.IsSingleCharacter = len(theOther) == 0
		vars.OtherParticipants = theOther
		errChan <- nil
	}()

	// Fetch world description
	go func() {
		w, ok := worlds.WorldById(session.WorldID)
		if !ok {
			errChan <- errors.New("error getting world")
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

		scenario, ok := scenarios.ScenarioById(*session.ScenarioID)
		if !ok {
			errChan <- errors.New("error getting scenario")
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
			return errors.Wrap(err, "failed to create template for character field")
		}

		*fieldPtr = util.WriteToString(tpl, characterVars)
	}

	return nil
}
