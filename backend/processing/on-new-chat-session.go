package processing

import (
	"context"
	"go.uber.org/zap"
	"juraji.nl/chat-quest/core/log"
	"juraji.nl/chat-quest/model/characters"
	chatsessions "juraji.nl/chat-quest/model/chat-sessions"
)

func init() {
	chatsessions.ChatSessionCreatedSignal.AddListener(onNewChatSessionHandleCharacterGreetings)
}

func onNewChatSessionHandleCharacterGreetings(
	_ context.Context,
	session *chatsessions.ChatSession,
) {
	if session == nil {
		log.Get().Debug("Received nil chat session, skipping greeting handling")
		return
	}

	sessionLog := log.Get().With(zap.Int("chatSessionId", session.ID))

	isGroupChat, err := chatsessions.IsGroupChatSession(session.ID)
	if err != nil {
		sessionLog.Error("Error checking if group chat session is a group chat", zap.Error(err))
	}

	if isGroupChat {
		createGroupSessionGreetings(session.ID, sessionLog)
	} else {
		createSingleCharacterGreeting(session.ID, sessionLog)
	}
}

func createGroupSessionGreetings(sessionId int, sessionLog *zap.Logger) {
	maxAttempts := 5 // Prevent infinite loops if no participants have greetings
	attempt := 0

	for attempt < maxAttempts {
		participantId, err := chatsessions.RandomChatSessionParticipantId(sessionId)
		if err != nil {
			sessionLog.Error("Error getting random chat session participant id in group session",
				zap.Int("attempt", attempt), zap.Error(err))
			return
		}

		greeting, err := characters.RandomGreetingByCharacterId(*participantId, true)
		if err != nil {
			sessionLog.Error("Error getting random group greeting",
				zap.Int("participant_id", *participantId), zap.Int("attempt", attempt), zap.Error(err))
			return
		}

		if greeting == nil {
			sessionLog.Debug("No greeting found for this participant id in group session, trying another random participant",
				zap.Int("participant_id", *participantId), zap.Int("attempt", attempt))
			attempt++
			continue
		}

		message := chatsessions.NewChatMessage(sessionId, false, participantId, *greeting)
		err = chatsessions.CreateChatMessage(sessionId, message)
		if err != nil {
			sessionLog.Error("Error creating chat message",
				zap.Int("participant_id", *participantId), zap.String("message", *greeting), zap.Error(err))
			return
		}

		return // Exit after successfully posting a greeting
	}

	sessionLog.Debug("No greetings found for any participants in group session after maximum attempts",
		zap.Int("maxAttempts", maxAttempts))
}

func createSingleCharacterGreeting(sessionId int, sessionLog *zap.Logger) {
	participantId, err := chatsessions.RandomChatSessionParticipantId(sessionId)
	if err != nil {
		sessionLog.Error("Error getting chat session participant id in session", zap.Error(err))
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
		sessionLog.Debug("No greeting found for participant id in session, skipping greeting")
		return
	}

	message := chatsessions.NewChatMessage(sessionId, false, participantId, *greeting)
	err = chatsessions.CreateChatMessage(sessionId, message)
	if err != nil {
		sessionLog.Error("Error creating chat message", zap.Error(err))
		return
	}
}
