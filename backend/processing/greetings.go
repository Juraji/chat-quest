package processing

import (
	"context"

	"go.uber.org/zap"
	"juraji.nl/chat-quest/core/log"
	"juraji.nl/chat-quest/model/characters"
	"juraji.nl/chat-quest/model/chat-sessions"
)

func GreetOnParticipantAdded(ctx context.Context, participant *chat_sessions.ChatParticipant) {
	if participant == nil || !participant.NewlyAdded {
		// Skip nil or not newly added
		return
	}

	sessionID := participant.ChatSessionID
	characterID := participant.CharacterID
	logger := log.Get().With(
		zap.Int("sessionId", sessionID),
		zap.Int("characterId", characterID))

	contextCheckPoint(ctx, logger)

	greeting, err := characters.RandomGreetingByCharacterId(characterID)
	if err != nil {
		logger.Error("Failed to fetch greeting", zap.Error(err))
		return
	}
	if greeting == nil {
		logger.Debug("Skipping empty greeting")
		return
	}

	message := chat_sessions.NewChatMessage(false, false, &characterID, *greeting)
	err = chat_sessions.CreateChatMessage(sessionID, message)
	if err != nil {
		logger.Error("Error creating chat message", zap.Error(err))
	}
}
