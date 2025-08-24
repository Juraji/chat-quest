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
	sessionLog := log.Get().With(zap.Int("chatSessionId", sessionID))

	isGroupChat, ok := sessions.IsGroupSession(sessionID)
	if !ok {
		return
	}

	participants, ok := sessions.GetParticipants(sessionID)
	if !ok {
		return
	}

	for _, participant := range participants {
		greeting, ok := characters.RandomGreetingByCharacterId(participant.ID, isGroupChat)
		if !ok {
			continue
		}
		if greeting == nil {
			sessionLog.Debug("Skipping empty greeting", zap.Int("participantId", participant.ID))
			continue
		}

		message := sessions.NewChatMessage(false, false, &participant.ID, *greeting)
		sessions.CreateChatMessage(sessionID, message)
	}
}
