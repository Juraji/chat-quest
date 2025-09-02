package chat_init

import (
	"context"
	"go.uber.org/zap"
	"juraji.nl/chat-quest/core/log"
	"juraji.nl/chat-quest/model/characters"
	sessions "juraji.nl/chat-quest/model/chat-sessions"
)

func CreateChatSessionGreetings(
	ctx context.Context,
	session *sessions.ChatSession,
) {
	if session == nil {
		log.Get().Debug("Received nil chat session, skipping greeting handling")
		return
	}

	if ctx.Err() != nil {
		log.Get().Debug("Cancelled by context")
		return
	}

	sessionID := session.ID
	logger := log.Get().With(zap.Int("chatSessionId", sessionID))

	isGroupChat, err := sessions.IsGroupSession(sessionID)
	if err != nil {
		logger.Error("Error checking group chat status", zap.Error(err))
		return
	}

	participants, err := sessions.GetAllParticipantsAsCharacters(sessionID)
	if err != nil {
		logger.Error("Error getting participants", zap.Error(err))
		return
	}

	for _, participant := range participants {
		greeting, err := characters.RandomGreetingByCharacterId(participant.ID, *isGroupChat)
		if err != nil {
			logger.Warn("Failed to fetch greeting",
				zap.Int("participantId", participant.ID), zap.Error(err))
			continue
		}
		if greeting == nil {
			logger.Debug("Skipping empty greeting",
				zap.Int("participantId", participant.ID))
			continue
		}

		message := sessions.NewChatMessage(false, false, &participant.ID, *greeting)
		err = sessions.CreateChatMessage(sessionID, message)
		if err != nil {
			logger.Error("Error creating chat message", zap.Error(err))
		}
	}
}
