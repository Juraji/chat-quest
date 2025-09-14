package processing

import (
	"context"

	"go.uber.org/zap"
	"juraji.nl/chat-quest/core/log"
	cs "juraji.nl/chat-quest/model/chat-sessions"
	p "juraji.nl/chat-quest/model/preferences"
)

func AutoArchiveMessages(
	ctx context.Context,
	message *cs.ChatMessage,
) {
	if message.IsGenerating {
		// Skip messages that are still being generated
		return
	}

	sessionID := message.ChatSessionID
	logger := log.Get().With(zap.Int("chatSessionId", sessionID))

	if contextCheckPoint(ctx, logger) {
		return
	}

	session, err := cs.GetById(sessionID)
	if err != nil {
		logger.Error("Error getting session", zap.Error(err))
		return
	}

	if session.GenerateMemories || !session.AutoArchiveMessages {
		// Auto archival is turned of or handled by memory gen.
		return
	}

	logger.Info("Archiving messages...")

	prefs, err := p.GetPreferences(true)
	if err != nil {
		logger.Error("Error getting preferences", zap.Error(err))
		return
	}

	messageWindow, err := getArchivalMessageWindow(logger, prefs, sessionID)
	if err != nil {
		logger.Error("Error getting message window", zap.Error(err))
		return
	}
	if messageWindow == nil {
		return
	}
	if contextCheckPoint(ctx, logger) {
		return
	}

	// Update message archival states
	for _, chatMessage := range messageWindow {
		err = cs.SetMessageArchived(sessionID, chatMessage.ID)
		if err != nil {
			logger.Error("Error setting message archived bit", zap.Error(err))
		}
	}

	logger.Info("Auto archival completed", zap.Int("archivedMessages", len(messageWindow)))
}

func getArchivalMessageWindow(
	logger *zap.Logger,
	prefs *p.Preferences,
	sessionID int,
) ([]cs.ChatMessage, error) {
	messages, err := cs.GetUnarchivedChatMessages(sessionID)
	if err != nil {
		return nil, err
	}

	triggerAfter := prefs.MemoryTriggerAfter
	windowSize := len(messages) - triggerAfter

	// Only proceed if we have enough messages to create a valid window
	if windowSize < 1 {
		logger.Info("Message window not yet full, skipping archival",
			zap.Int("windowSize", windowSize))
		return nil, nil
	}

	return messages[:windowSize], nil
}
