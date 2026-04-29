package processing

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"juraji.nl/chat-quest/core/log"
	"juraji.nl/chat-quest/core/util"
	"juraji.nl/chat-quest/model/characters"
	"juraji.nl/chat-quest/model/chat-sessions"
)

func GreetOnParticipantAdded(_ context.Context, participant *chat_sessions.ChatParticipant) error {
	if participant == nil || !participant.NewlyAdded {
		// Skip nil or not newly added
		return nil
	}

	sessionID := participant.ChatSessionID
	characterID := participant.CharacterID
	logger := log.Get().With(
		zap.Int("sessionId", sessionID),
		zap.Int("characterId", characterID))
	var err error

	defer func() {
		if err != nil {
			logger.Error("Error greeting as new participant", zap.Error(err))
		}
	}()

	greeting, err := characters.RandomGreetingByCharacterId(characterID)
	if err != nil {
		return errors.Wrap(err, "failed to fetch greeting")
	}
	if greeting == nil {
		logger.Debug("Skipping empty greeting")
		return nil
	}

	message := chat_sessions.NewChatMessage(false, false, &characterID, *greeting)

	if util.ContainsTemplateVars(message.Content) {
		char, err := characters.CharacterById(characterID)
		if err != nil {
			return errors.Wrap(err, "failed to fetch character")
		}

		vars := NewGreetingVars(sessionID, char)
		result, err := util.ParseAndApplyTextTemplate("Character Greeting", message.Content, vars)
		if err != nil {
			return errors.Wrap(err, "failed to parse and apply text template")
		}

		message.Content = result
	}

	err = chat_sessions.CreateChatMessage(sessionID, message)
	if err != nil {
		return errors.Wrap(err, "error creating chat message")
	}

	return nil
}
