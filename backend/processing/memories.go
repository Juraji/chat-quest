package processing

import (
	"context"
	"encoding/json"
	"strings"
	"sync"

	"go.uber.org/zap"
	"juraji.nl/chat-quest/core/log"
	p "juraji.nl/chat-quest/core/providers"
	cs "juraji.nl/chat-quest/model/chat-sessions"
	i "juraji.nl/chat-quest/model/instructions"
	m "juraji.nl/chat-quest/model/memories"
	"juraji.nl/chat-quest/model/preferences"
)

var generationMutex sync.Mutex

func GenerateMemoriesForMessageID(
	ctx context.Context,
	messageId int,
) {
	// Lock while processing to avoid multiple messages invoking simultaneous generation
	// for the same message window.
	generationMutex.Lock()
	defer generationMutex.Unlock()

	logger := log.Get().With(
		zap.Int("sourceMessageId", messageId))

	if contextCheckPoint(ctx, logger) {
		return
	}

	logger.Info("Generating memories for specific message...")

	message, err := cs.GetMessageById(messageId)
	if err != nil {
		logger.Error("Error fetching message", zap.Error(err))
	}

	sessionID := message.ChatSessionID
	logger = logger.With(
		zap.Int("sessionId", sessionID))

	session, err := cs.GetById(sessionID)
	if err != nil {
		logger.Error("Error getting session", zap.Error(err))
		return
	}

	prefs, err := preferences.GetPreferences(true)
	if err != nil {
		logger.Error("Error getting preferences", zap.Error(err))
		return
	}

	if contextCheckPoint(ctx, logger) {
		return
	}

	messageWindow := []cs.ChatMessage{*message}

	memories, ok := generateMemories(logger, ctx, session, prefs, messageWindow)
	if !ok {
		return
	}

	if contextCheckPoint(ctx, logger) {
		return
	}

	// We're done, save memories
	for _, memory := range memories {
		err = m.CreateMemory(session.WorldID, memory)
		if err != nil {
			logger.Error("Error creating memory", zap.Error(err))
			return
		}
	}

	logger.Info("Memory generation completed", zap.Int("newMemories", len(memories)))
}

func GenerateMemories(
	ctx context.Context,
	message *cs.ChatMessage,
) {
	if message.IsGenerating {
		// Skip messages that are still being generated
		return
	}

	// Lock while processing to avoid multiple messages invoking simultaneous generation
	// for the same message window.
	generationMutex.Lock()
	defer generationMutex.Unlock()

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
	if !session.GenerateMemories {
		// Memories are disabled for this session
		return
	}

	logger.Info("Generating memories...")

	prefs, err := preferences.GetPreferences(true)
	if err != nil {
		logger.Error("Error getting preferences", zap.Error(err))
		return
	}

	messageWindow, err := getMessageWindow(logger, prefs, sessionID)
	if err != nil {
		logger.Error("Error getting message window", zap.Error(err))
		return
	}
	if messageWindow == nil {
		return
	}

	memories, ok := generateMemories(logger, ctx, session, prefs, messageWindow)
	if !ok {
		return
	}
	if contextCheckPoint(ctx, logger) {
		return
	}

	// We're done, save memories
	for _, memory := range memories {
		err = m.CreateMemory(session.WorldID, memory)
		if err != nil {
			logger.Error("Error creating memory", zap.Error(err))
			return
		}
	}

	// Update message processed states
	for _, chatMessage := range messageWindow {
		err = cs.SetMessageArchived(sessionID, chatMessage.ID)
		if err != nil {
			logger.Error("Error setting message archived bit", zap.Error(err))
		}
	}

	logger.Info("Memory generation completed", zap.Int("newMemories", len(memories)))
}

func generateMemories(
	logger *zap.Logger,
	ctx context.Context,
	session *cs.ChatSession,
	prefs *preferences.Preferences,
	messageWindow []cs.ChatMessage,
) ([]*m.Memory, bool) {
	instruction, err := i.InstructionById(*prefs.MemoriesInstructionId)
	if err != nil {
		logger.Error("Could not fetch memory instruction", zap.Error(err))
		return nil, false
	}
	modelInstance, err := p.GetLlmModelInstanceById(*prefs.MemoriesModelId)
	if err != nil {
		logger.Error("Could not fetch memory model", zap.Error(err))
		return nil, false
	}

	lastTimestampInWindow := *messageWindow[len(messageWindow)-1].CreatedAt
	templateVars := NewMemoryInstructionVars(session, lastTimestampInWindow)
	if err = instruction.ApplyTemplates(templateVars); err != nil {
		logger.Error("Error applying instruction templates", zap.Error(err))
		return nil, false
	}

	// Generate memories
	requestMessages := createChatRequestMessages(messageWindow, instruction)
	chatResponseChan := p.GenerateChatResponse(modelInstance, requestMessages, instruction.AsLlmParameters())
	var memoryGenResponse string

responseLoop:
	for {
		select {
		case r, hasNext := <-chatResponseChan:
			if !hasNext {
				// Done
				break responseLoop
			}

			if r.Error != nil {
				logger.Error("Error generating memories",
					zap.String("generated", memoryGenResponse),
					zap.Error(r.Error))
				return nil, false
			}

			memoryGenResponse = memoryGenResponse + r.Content
		case <-ctx.Done():
			logger.Debug("Cancelled by context")
			return nil, false
		}
	}

	// We expect a JSON markdown block, extract it
	markDownStartSeq := "```json"
	markDownEndSeq := "```"
	markdownStart := strings.Index(memoryGenResponse, markDownStartSeq)
	if markdownStart == -1 {
		logger.Error("Could not find markdown start in response")
		return nil, false
	}
	memoryGenResponse = memoryGenResponse[markdownStart+len(markDownStartSeq):]
	markdownEnd := strings.Index(memoryGenResponse, markDownEndSeq)
	if markdownEnd == -1 {
		logger.Error("Could not find markdown end in response")
		return nil, false
	}
	memoryGenResponse = memoryGenResponse[:markdownEnd]

	// Unmarshal memories
	var memories []*m.Memory
	err = json.Unmarshal([]byte(memoryGenResponse), &memories)
	if err != nil {
		logger.Error("Could not unmarshal memory response", zap.Error(err))
		return nil, false
	}

	logger.Debug("Memories generated successfully", zap.Int("memoryCount", len(memories)))
	return memories, true
}

func getMessageWindow(logger *zap.Logger, prefs *preferences.Preferences, sessionID int) ([]cs.ChatMessage, error) {
	messages, err := cs.GetUnarchivedChatMessages(sessionID)
	if err != nil {
		return nil, err
	}

	triggerAfter := prefs.MemoryTriggerAfter
	requiredWindowSize := prefs.MemoryWindowSize

	windowSize := len(messages) - triggerAfter

	// Only proceed if we have enough messages to create a valid window
	if windowSize < requiredWindowSize {
		logger.Info("Message window not yet full, skipping generation",
			zap.Int("requiredWindowSize", requiredWindowSize),
			zap.Int("windowSize", windowSize))
		return nil, nil
	}

	return messages[:windowSize], nil
}
