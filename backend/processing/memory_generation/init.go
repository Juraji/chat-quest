package memory_generation

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"juraji.nl/chat-quest/core/log"
	p "juraji.nl/chat-quest/core/providers"
	"juraji.nl/chat-quest/model/characters"
	chatsessions "juraji.nl/chat-quest/model/chat-sessions"
	"juraji.nl/chat-quest/model/instructions"
	"juraji.nl/chat-quest/model/memories"
)

type instructionTemplateVars struct {
	Participants []characters.Character
}

func GenerateMemories(
	ctx context.Context,
	message *chatsessions.ChatMessage,
) {
	if message.IsUser || message.IsGenerating {
		// Skip messages by the user or messages that are still being generated
		return
	}

	sessionID := message.ChatSessionID
	logger := log.Get().With(zap.Int("chatSessionId", sessionID))
	if ctx.Err() != nil {
		logger.Debug("Cancelled by context")
		return
	}

	preferences, ok := memories.GetMemoryPreferences()
	if !ok {
		return
	}
	if err := preferences.Validate(); err != nil {
		logger.Error("Error validating memory preferences", zap.Error(err))
		return
	}

	memorizableMessages, ok := getMemorizableMessages(logger, preferences, sessionID)
	if !ok {
		return
	}

	_, ok = generateRawMemoriesResponse(logger, ctx, sessionID, preferences, memorizableMessages)
	if !ok {
		return
	}

	//embeddingModel, ok := p.GetLlmModelInstanceById(*preferences.EmbeddingModelID)
	//if !ok {
	//	return
	//}
}

func generateRawMemoriesResponse(
	logger *zap.Logger,
	ctx context.Context,
	sessionID int,
	preferences *memories.MemoryPreferences,
	messages []chatsessions.ChatMessage,
) ([]memories.Memory, bool) {
	instruction, ok := instructions.InstructionById(*preferences.MemoriesInstructionID)
	if !ok || instruction == nil {
		logger.Debug("Could not fetch memory instruction")
		return nil, false
	}
	modelInstance, ok := p.GetLlmModelInstanceById(*preferences.MemoriesModelID)
	if !ok || modelInstance == nil {
		logger.Debug("Could not fetch memory model")
		return nil, false
	}

	participants, ok := chatsessions.GetParticipants(sessionID)
	if !ok {
		logger.Debug("Could not fetch participants")
		return nil, false
	}

	templateVars := instructionTemplateVars{
		Participants: participants,
	}

	instruction, err := instructions.ApplyInstructionTemplates(*instruction, templateVars)
	if err != nil {
		logger.Error("Error applying instruction templates", zap.Error(err))
		return nil, false
	}

	requestMessages := createChatRequestMessages(messages, instruction)

	// Do LLM
	llmResponse, ok := callLlm(logger, ctx, modelInstance, requestMessages, instruction.Temperature)

	logger.Sugar().Debugf("LLM Response: %v", llmResponse)
	// TODO: Parse memories output
	return nil, false
}

func callLlm(
	logger *zap.Logger,
	ctx context.Context,
	instance *p.LlmModelInstance,
	messages []p.ChatRequestMessage,
	temperature *float32,
) (string, bool) {
	chatResponseChan := p.GenerateChatResponse(instance, messages, temperature)
	var memoryGenResponse string

	for {
		select {
		case r, hasNext := <-chatResponseChan:
			if !hasNext {
				// Done
				return memoryGenResponse, true
			}

			if r.Error != nil {
				logger.Error("Error generating memories", zap.Error(r.Error))
				return "", false
			}

			memoryGenResponse = memoryGenResponse + r.Content
		case <-ctx.Done():
			logger.Debug("Cancelled by context")
			return "", false
		}
	}
}

func createChatRequestMessages(
	chatMessages []chatsessions.ChatMessage,
	instruction *instructions.InstructionTemplate,
) []p.ChatRequestMessage {
	// Pre-allocate messages with history len + max number of messages added here
	messages := make([]p.ChatRequestMessage, 0, len(chatMessages)+3)

	// Add system and world setup messages
	messages = append(messages,
		p.ChatRequestMessage{Role: p.RoleSystem, Content: instruction.SystemPrompt},
		p.ChatRequestMessage{Role: p.RoleUser, Content: instruction.WorldSetup},
	)

	// Add chat history
	for _, m := range chatMessages {
		if m.IsUser {
			messages = append(messages, p.ChatRequestMessage{Role: p.RoleUser, Content: m.Content})
		} else {
			content := fmt.Sprintf("<ByCharacterId>%v</ByCharacterId>\n\n%s", *m.CharacterID, m.Content)
			messages = append(messages, p.ChatRequestMessage{Role: p.RoleAssistant, Content: content})
		}
	}

	// Add user instruction message
	messages = append(messages, p.ChatRequestMessage{Role: p.RoleUser, Content: instruction.Instruction})
	return messages
}

func getMemorizableMessages(
	logger *zap.Logger,
	preferences *memories.MemoryPreferences,
	sessionID int,
) ([]chatsessions.ChatMessage, bool) {
	unprocessed, ok := memories.GetUnprocessedMessagesForSession(sessionID)
	if !ok {
		logger.Warn("Failed to get unprocessed messages")
		return nil, false
	}

	triggerAfter := preferences.MemoryTriggerAfter
	requiredWindowSize := preferences.MemoryWindowSize
	unprocessedMessagesLen := len(unprocessed)
	truncateBy := unprocessedMessagesLen - triggerAfter

	var window []chatsessions.ChatMessage
	if truncateBy > 0 {
		window = unprocessed[:truncateBy]
	}
	windowSize := len(window)

	logger = logger.With(
		zap.Int("triggerAfter", triggerAfter),
		zap.Int("requiredWindowSize", requiredWindowSize),
		zap.Int("totalUnprocessedMessages", unprocessedMessagesLen),
		zap.Int("windowSize", windowSize))

	if windowSize >= requiredWindowSize {
		logger.Debug("Found messages to memorize")
		return window, true
	} else {
		logger.Debug("Not enough messages to memorize")
		return nil, false
	}
}
