package chat_response

import (
	"context"
	"go.uber.org/zap"
	"juraji.nl/chat-quest/core/log"
	p "juraji.nl/chat-quest/core/providers"
	chatsessions "juraji.nl/chat-quest/model/chat-sessions"
	"juraji.nl/chat-quest/model/instructions"
	"juraji.nl/chat-quest/model/worlds"
)

func GenerateResponseForParticipant(ctx context.Context, payload *chatsessions.ChatParticipant) {
	sessionId := payload.ChatSessionID
	characterId := payload.CharacterID

	logger := log.Get().With(
		zap.String("source", "GenerateResponseForParticipant"),
		zap.Int("chatSessionId", sessionId),
		zap.Int("characterId", characterId))

	if ctx.Err() != nil {
		logger.Debug("Cancelled by context")
		return
	}

	// Fetch instruction and connection preferences
	prefs, ok := worlds.GetChatPreferences()
	if !ok {
		return
	}
	if err := prefs.Validate(); err != nil {
		logger.Error("Error validating chat preferences", zap.Error(err))
		return
	}

	instruction, ok := instructions.InstructionById(*prefs.ChatInstructionID)
	if !ok {
		return
	}

	modelInstance, ok := p.GetLlmModelInstanceById(*prefs.ChatModelID)
	if !ok {
		return
	}

	// Fetch session details
	session, ok := chatsessions.GetById(sessionId)
	if !ok {
		return
	}

	// Fetch current chat history
	chatHistory, ok := chatsessions.GetChatMessages(sessionId)
	if !ok {
		return
	}

	// Process instruction templates
	templateVars, err := newInstructionTemplateVars(session, nil, chatHistory, characterId)
	if err != nil {
		logger.Error("Error collecting instruction template variables", zap.Error(err))
		return
	}

	instruction, err = instructions.ApplyInstructionTemplates(*instruction, templateVars)
	if err != nil {
		logger.Error("Error processing instruction templates", zap.Error(err))
		return
	}

	generateResponse(ctx, logger, sessionId, &characterId, instruction, modelInstance, chatHistory)
}

func GenerateResponseForMessage(
	ctx context.Context,
	message *chatsessions.ChatMessage,
) {
	if message == nil || !message.IsUser {
		// Ignore nil or non-user
		return
	}

	sessionId := message.ChatSessionID
	logger := log.Get().With(
		zap.String("source", "GenerateResponseForMessage"),
		zap.Int("chatSessionId", sessionId),
		zap.Int("triggerMessageId", message.ID))

	if ctx.Err() != nil {
		logger.Debug("Cancelled by context")
		return
	}

	// Fetch instruction and connection preferences
	prefs, ok := worlds.GetChatPreferences()
	if !ok {
		return
	}
	if err := prefs.Validate(); err != nil {
		return
	}

	instruction, ok := instructions.InstructionById(*prefs.ChatInstructionID)
	if !ok {
		return
	}

	modelInstance, ok := p.GetLlmModelInstanceById(*prefs.ChatModelID)
	if !ok {
		return
	}

	// Fetch session details
	session, ok := chatsessions.GetById(sessionId)
	if !ok {
		return
	}

	// Select participant to respond with
	characterId, ok := chatsessions.RandomParticipantId(sessionId)
	if !ok {
		return
	}
	if characterId == nil {
		logger.Warn("error selecting character to respond with: no participants found")
		return
	}

	// Fetch current chat history
	chatHistory, ok := chatsessions.GetChatMessagesPreceding(sessionId, message.ID)
	if !ok {
		return
	}

	logger = logger.With(zap.Int("selectedCharacterId", *characterId))

	// Process instruction templates
	templateVars, err := newInstructionTemplateVars(session, message, chatHistory, *characterId)
	if err != nil {
		logger.Error("Error collecting instruction template variables", zap.Error(err))
		return
	}

	instruction, err = instructions.ApplyInstructionTemplates(*instruction, templateVars)
	if err != nil {
		logger.Error("Error processing instruction templates", zap.Error(err))
		return
	}

	generateResponse(ctx, logger, sessionId, characterId, instruction, modelInstance, chatHistory)
}

func generateResponse(
	ctx context.Context,
	logger *zap.Logger,
	sessionId int,
	characterId *int,
	instruction *instructions.InstructionTemplate,
	modelInstance *p.LlmModelInstance,
	chatHistory []chatsessions.ChatMessage,
) {
	// Check for cancellation
	if ctx.Err() != nil {
		logger.Debug("Cancelled by context")
		return
	}

	// Build messages to send to LLM
	messages := createChatRequestMessages(instruction, chatHistory)

	// Do LLM and handle output
	chatResponseChan := p.GenerateChatResponse(modelInstance, messages, instruction.Temperature)

	// Create response message
	responseMessage := chatsessions.NewChatMessage(false, true, characterId, "")
	if ok := chatsessions.CreateChatMessage(sessionId, responseMessage); !ok {
		logger.Warn("Failed to create response chat message")
		return
	}
	defer func() {
		responseMessage.IsGenerating = false
		if ok := chatsessions.UpdateChatMessage(sessionId, responseMessage.ID, responseMessage); !ok {
			logger.Error("Failed to update response chat message upon finalization")
		}
	}()

	for {
		select {
		case response, hasNext := <-chatResponseChan:
			if !hasNext {
				return
			}
			if response.Error != nil {
				logger.Error("Error generating response", zap.Error(response.Error))
				return
			}

			responseMessage.Content = responseMessage.Content + response.Content
			if ok := chatsessions.UpdateChatMessage(sessionId, responseMessage.ID, responseMessage); !ok {
				logger.Error("Failed to update response chat message")
				return
			}
		case <-ctx.Done():
			logger.Debug("Cancelled by context")
			return
		}
	}
}

func createChatRequestMessages(
	instruction *instructions.InstructionTemplate,
	chatHistory []chatsessions.ChatMessage,
) []p.ChatRequestMessage {
	// Pre-allocate messages with history len + max number of messages added here
	messages := make([]p.ChatRequestMessage, 0, len(chatHistory)+4)

	// Add system and world setup messages
	messages = append(messages,
		p.ChatRequestMessage{Role: p.RoleSystem, Content: instruction.SystemPrompt},
		p.ChatRequestMessage{Role: p.RoleUser, Content: instruction.WorldSetup},
	)

	// Add initial assistant message if needed, to maintain role order
	if len(chatHistory) == 0 || chatHistory[0].IsUser {
		messages = append(messages, p.ChatRequestMessage{
			Role:    p.RoleAssistant,
			Content: "[OOC: Understood, I will from now on respond explicitly adhering to the instructions and information provided.]",
		})
	}

	// Add chat history
	for _, m := range chatHistory {
		if m.IsUser {
			messages = append(messages, p.ChatRequestMessage{Role: p.RoleUser, Content: m.Content})
		} else {
			messages = append(messages, p.ChatRequestMessage{Role: p.RoleAssistant, Content: m.Content})
		}
	}

	// Add user instruction message
	messages = append(messages, p.ChatRequestMessage{Role: p.RoleUser, Content: instruction.Instruction})

	return messages
}
