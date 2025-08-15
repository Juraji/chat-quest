package processing

import (
	"context"
	"go.uber.org/zap"
	"juraji.nl/chat-quest/core/log"
	"juraji.nl/chat-quest/model/characters"
	chatsessions "juraji.nl/chat-quest/model/chat-sessions"
)

func onNewChatSessionHandleCharacterGreetings(
	ctx context.Context,
	session *chatsessions.ChatSession,
) {
	if session == nil {
		log.Get().Debug("Received nil chat session, skipping greeting handling")
		return
	}

	if ctx.Err() != nil {
		log.Get().Debug("Cancelled by context")
		return
	}

	sessionLog := log.Get().With(zap.Int("chatSessionId", session.ID))

	isGroupChat, err := chatsessions.IsGroupSession(session.ID)
	if err != nil {
		sessionLog.Error("Error checking if group chat session is a group chat", zap.Error(err))
	}

	if isGroupChat {
		createGroupSessionGreetings(ctx, session.ID, sessionLog)
	} else {
		createSingleCharacterGreeting(ctx, session.ID, sessionLog)
	}
}

func createGroupSessionGreetings(ctx context.Context, sessionId int, sessionLog *zap.Logger) {
	maxAttempts := 5 // Prevent infinite loops if no participants have greetings
	attempt := 0

	for attempt < maxAttempts {
		if ctx.Err() != nil {
			log.Get().Debug("Cancelled by context")
			return
		}

		participantId, err := chatsessions.RandomParticipantId(sessionId)
		if err != nil {
			sessionLog.Error("Error getting random chat session participant in group session",
				zap.Int("attempt", attempt), zap.Error(err))
			return
		}

		greeting, err := characters.RandomGreetingByCharacterId(*participantId, true)
		if err != nil {
			sessionLog.Error("Error getting random group greeting",
				zap.Int("participantId", *participantId), zap.Int("attempt", attempt), zap.Error(err))
			return
		}

		if greeting == nil {
			sessionLog.Debug("No greeting found, re-rolling for participant",
				zap.Int("participantId", *participantId), zap.Int("attempt", attempt))
			attempt++
			continue
		}

		message := chatsessions.NewChatMessage(sessionId, false, participantId, *greeting)
		err = chatsessions.CreateChatMessage(sessionId, message)
		if err != nil {
			sessionLog.Error("Error creating chat message",
				zap.Int("participantId", *participantId), zap.String("message", *greeting), zap.Error(err))
			return
		}

		return // Exit after successfully posting a greeting
	}

	sessionLog.Debug("No greetings found for any participants in group session after maximum attempts",
		zap.Int("maxAttempts", maxAttempts))
}

func createSingleCharacterGreeting(ctx context.Context, sessionId int, sessionLog *zap.Logger) {
	if ctx.Err() != nil {
		log.Get().Debug("Cancelled by context")
		return
	}

	participantId, err := chatsessions.RandomParticipantId(sessionId)
	if err != nil {
		sessionLog.Error("Error getting chat session participant in session", zap.Error(err))
		return
	}

	if participantId == nil {
		sessionLog.Debug("No participants in this session.")
		return
	}

	greeting, err := characters.RandomGreetingByCharacterId(*participantId, false)
	if err != nil {
		sessionLog.Error("Error getting random greeting", zap.Error(err))
		return
	}

	if greeting == nil {
		sessionLog.Debug("No greeting found, skipping greeting")
		return
	}

	message := chatsessions.NewChatMessage(sessionId, false, participantId, *greeting)
	err = chatsessions.CreateChatMessage(sessionId, message)
	if err != nil {
		sessionLog.Error("Error creating chat message", zap.Error(err))
		return
	}
}
