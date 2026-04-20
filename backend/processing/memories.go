package processing

import (
	"context"
	"encoding/json"
	"slices"
	"sync"

	"go.uber.org/zap"
	"juraji.nl/chat-quest/core/log"
	p "juraji.nl/chat-quest/core/providers"
	cs "juraji.nl/chat-quest/model/chat-sessions"
	i "juraji.nl/chat-quest/model/instructions"
	m "juraji.nl/chat-quest/model/memories"
	pf "juraji.nl/chat-quest/model/preferences"
)

const memoriesResponseFormat = `{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "required": ["memories"],
  "properties": {
    "memories": {
      "type": "array",
      "items": {
        "type": "object",
        "required": ["characterId","content"],
        "properties": {
          "characterId": {"type": "number"},
          "content": {"type": "string"}
        }
      }
    }
  }
}`

// OpenAI requires an object type root.
type memoriesContainer struct {
	Memories []*m.Memory
}

var memoryGenerationMutex sync.Mutex

func UpdateBookmarkOnMemoryGenEnable(_ context.Context, e *cs.ChatSessionUpdatedBAEvent) {
	if e.Before.GenerateMemories == e.After.GenerateMemories {
		return
	}

	logger := log.Get().With(
		zap.Int("sessionID", e.SessionId))

	logger.Info("Memory generation enabled for session, moving bookmark to latest message...")

	var err error
	// The user has re-enabled memory generation, assuming it starts enabled.
	// Here we move the memory bookmark to the last message in the chat, so we ignore messages prior to this point.
	messages, err := cs.GetTailChatMessages(e.SessionId, 1)
	if err != nil {
		logger.Error("Error getting last chat message", zap.Error(err))
		return
	}
	if len(messages) == 0 {
		logger.Info("No messages found, skipping.")
		return
	}

	messageID := messages[0].ID
	logger = logger.With(zap.Int("messageID", messageID))

	err = m.SetMemoryBookmark(e.SessionId, messageID)
	if err != nil {
		logger.Error("Error setting bookmark", zap.Error(err))
		return
	}

	logger.Info("Updated bookmark successfully")
}

func GenerateMemoriesForMessageID(
	ctx context.Context,
	messageId int,
	includeNPreceding int,
) {
	// Lock while processing to avoid multiple messages invoking simultaneous generation
	// for the same message window.
	// If the lock is already active we cancel this invocation.
	lock := memoryGenerationMutex.TryLock()
	if !lock {
		return
	}
	defer memoryGenerationMutex.Unlock()

	logger := log.Get().With(
		zap.Int("sourceMessageId", messageId),
		zap.Int("includeNPreceding", includeNPreceding))

	// Cancellation
	ctx, cleanup := setupCancelBySystem(ctx, logger, "GenerateMemories")
	defer cleanup()

	logger.Info("Generating memories for specific message...")

	message, err := cs.GetMessageById(messageId)
	if err != nil {
		logger.Error("Error fetching message", zap.Error(err))
		return
	}

	sessionID := message.ChatSessionID
	logger = logger.With(
		zap.Int("sessionId", sessionID))

	session, err := cs.GetById(sessionID)
	if err != nil {
		logger.Error("Error fetching session", zap.Error(err))
		return
	}

	prefs, err := pf.GetPreferences(true)
	if err != nil {
		logger.Error("Error fetching preferences", zap.Error(err))
		return
	}

	if contextCheckPoint(ctx, logger) {
		return
	}

	precedingMessages, err := cs.GetMessagesInSessionBeforeId(message.ChatSessionID, messageId, includeNPreceding)
	if err != nil {
		logger.Error("Error fetching previous messages", zap.Error(err))
		return
	}

	// Create message window (preceding messages + current message)
	messageWindow := append(precedingMessages, *message)

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
	return
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
	// If the lock is already active we cancel this invocation.
	lock := memoryGenerationMutex.TryLock()
	if !lock {
		return
	}
	defer memoryGenerationMutex.Unlock()

	sessionID := message.ChatSessionID
	logger := log.Get().With(zap.Int("chatSessionId", sessionID))

	// Cancellation
	ctx, cleanup := setupCancelBySystem(ctx, logger, "GenerateMemories")
	defer cleanup()

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

	prefs, err := pf.GetPreferences(true)
	if err != nil {
		logger.Error("Error getting preferences", zap.Error(err))
		return
	}

	var messageWindow []cs.ChatMessage
	{
		bookmark, err := m.GetMemoryBookmark(sessionID)
		if err != nil {
			logger.Error("Error getting message bookmark", zap.Error(err))
			return
		}

		limit := prefs.MemoryWindowSize + prefs.MemoryTriggerAfter
		if bookmark == nil {
			// No bookmark set yet, just get the tail messages
			messageWindow, err = cs.GetTailChatMessages(sessionID, limit)
		} else {
			// Get messages AFTER the current bookmark, up to the maximum
			messageWindow, err = cs.GetMessagesInSessionAfterId(sessionID, *bookmark, limit)
		}
		if err != nil {
			logger.Error("Error getting chat messages", zap.Error(err))
			return
		}

		windowSize := len(messageWindow) - prefs.MemoryTriggerAfter
		if windowSize < prefs.MemoryWindowSize {
			logger.Info("Skipping memory generation because window size is too small",
				zap.Int("requiredWindowSize", prefs.MemoryWindowSize),
				zap.Int("windowSize", windowSize))
			return
		}

		messageWindow = messageWindow[:windowSize]
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

	// Update bookmark
	lastMessageId := messageWindow[len(messageWindow)-1].ID
	if err = m.SetMemoryBookmark(sessionID, lastMessageId); err != nil {
		logger.Error("Error setting message bookmark ID", zap.Error(err))
		return
	}

	logger.Info("Memory generation completed", zap.Int("newMemories", len(memories)))
}

func generateMemories(
	logger *zap.Logger,
	ctx context.Context,
	session *cs.ChatSession,
	prefs *pf.Preferences,
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

	// Log instruction contents
	logInstructionsToFile(logger, instruction, messageWindow)

	// Generate memories
	requestMessages := createChatRequestMessages(messageWindow, instruction)
	llmParameters := instruction.AsLlmParameters()
	llmParameters.ResponseFormat = new(memoriesResponseFormat)

	chatResponseChan := p.GenerateChatResponse(ctx, modelInstance, requestMessages, llmParameters)
	var memoryGenResponse string

responseLoop:
	for {
		select {
		case r, hasNext := <-chatResponseChan:
			if r.Error != nil {
				logger.Error("Error generating memories",
					zap.String("generated", memoryGenResponse),
					zap.Error(r.Error))
				return nil, false
			}

			memoryGenResponse = memoryGenResponse + r.Content
			if !hasNext {
				// Done
				break responseLoop
			}
		case <-ctx.Done():
			logger.Debug("Cancelled by context")
			return nil, false
		}
	}

	// Unmarshal memories
	var container memoriesContainer
	err = json.Unmarshal([]byte(memoryGenResponse), &container)
	if err != nil {
		logger.Error("Could not unmarshal memory response", zap.Error(err))
		return nil, false
	}

	memoryFilter := func(memory *m.Memory) bool { return len(memory.Content) == 0 }
	memories := slices.DeleteFunc(container.Memories, memoryFilter)
	logger.Debug("Memories generated successfully", zap.Int("memoryCount", len(memories)))
	return memories, true
}
