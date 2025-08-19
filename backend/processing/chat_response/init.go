package chat_response

import (
	"context"
	"go.uber.org/zap"
	"juraji.nl/chat-quest/core/log"
	p "juraji.nl/chat-quest/core/providers"
	"juraji.nl/chat-quest/core/util"
	chatsessions "juraji.nl/chat-quest/model/chat-sessions"
	"juraji.nl/chat-quest/model/instructions"
	"juraji.nl/chat-quest/model/worlds"
	"strings"
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
	prefs, err := worlds.GetChatPreferences()
	if err != nil {
		logger.Error("Error fetching chat preferences", zap.Error(err))
		return
	}
	if err = prefs.Validate(); err != nil {
		logger.Error("Error validating chat preferences", zap.Error(err))
		return
	}

	instruction, err := instructions.InstructionById(*prefs.ChatInstructionID)
	if err != nil {
		logger.Error("Error fetching chat instruction", zap.Error(err))
		return
	}

	llmModel, err := p.LlmModelById(*prefs.ChatModelID)
	if err != nil {
		logger.Error("Error fetching preferred chat model", zap.Error(err))
		return
	}

	connectionProfile, err := p.ConnectionProfileById(llmModel.ConnectionProfileId)
	if err != nil {
		logger.Error("Error fetching preferred connection profile", zap.Error(err))
		return
	}

	// Fetch session details
	session, err := chatsessions.GetById(sessionId)
	if err != nil {
		logger.Error("Failed to get session", zap.Error(err))
		return
	}

	// Fetch current chat history
	chatHistory, err := chatsessions.GetChatMessages(sessionId)
	if err != nil {
		logger.Error("Error fetching chat history", zap.Error(err))
		return
	}

	// Process instruction templates
	ok := processInstructionTemplates(logger, instruction, session, chatHistory, characterId, nil)
	if !ok {
		return
	}

	generateResponse(ctx, logger, sessionId, &characterId, instruction, llmModel, connectionProfile, chatHistory)
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
	prefs, err := worlds.GetChatPreferences()
	if err != nil {
		logger.Error("Error fetching chat preferences", zap.Error(err))
		return
	}
	if err = prefs.Validate(); err != nil {
		logger.Error("Error validating chat preferences", zap.Error(err))
		return
	}

	instruction, err := instructions.InstructionById(*prefs.ChatInstructionID)
	if err != nil {
		logger.Error("Error fetching chat instruction", zap.Error(err))
		return
	}

	llmModel, err := p.LlmModelById(*prefs.ChatModelID)
	if err != nil {
		logger.Error("Error fetching preferred chat model", zap.Error(err))
		return
	}

	connectionProfile, err := p.ConnectionProfileById(llmModel.ConnectionProfileId)
	if err != nil {
		logger.Error("Error fetching preferred connection profile", zap.Error(err))
		return
	}

	// Fetch session details
	session, err := chatsessions.GetById(sessionId)
	if err != nil {
		logger.Error("Failed to get session", zap.Error(err))
		return
	}

	// Select participant to respond with
	characterId, err := chatsessions.RandomParticipantId(sessionId)
	if err != nil {
		logger.Error("Error selecting character to respond with", zap.Error(err))
		return
	} else if characterId == nil {
		logger.Warn("error selecting character to respond with: no participants found")
		return
	}

	// Fetch current chat history
	chatHistory, err := chatsessions.GetChatMessagesPreceding(sessionId, message.ID)
	if err != nil {
		logger.Error("Error fetching chat history", zap.Error(err))
		return
	}

	logger = logger.With(zap.Int("selectedCharacterId", *characterId))

	// Process instruction templates
	ok := processInstructionTemplates(logger, instruction, session, chatHistory, *characterId, message)
	if !ok {
		return
	}

	generateResponse(ctx, logger, sessionId, characterId, instruction, llmModel, connectionProfile, chatHistory)
}

func generateResponse(
	ctx context.Context,
	logger *zap.Logger,
	sessionId int,
	characterId *int,
	instruction *instructions.InstructionTemplate,
	llmModel *p.LlmModel,
	connectionProfile *p.ConnectionProfile,
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
	chatResponseChan := connectionProfile.GenerateChatResponse(messages, *llmModel, instruction.Temperature)

	// Create response message
	responseMessage := chatsessions.NewChatMessage(false, false, true, characterId, "")
	if err := chatsessions.CreateChatMessage(sessionId, responseMessage); err != nil {
		logger.Warn("Failed to create response chat message", zap.Error(err))
		return
	}

	for {
		select {
		case response, ok := <-chatResponseChan:
			if !ok {
				responseMessage.IsGenerating = false
				if err := chatsessions.UpdateChatMessage(sessionId, responseMessage.ID, responseMessage); err != nil {
					logger.Error("Failed to update response chat message upon finalization", zap.Error(err))
				}
				return
			}
			if response.Error != nil {
				logger.Error("Error generating response", zap.Error(response.Error))
				return
			}

			responseMessage.Content = responseMessage.Content + response.Content
			if err := chatsessions.UpdateChatMessage(sessionId, responseMessage.ID, responseMessage); err != nil {
				logger.Error("Failed to update response chat message", zap.Error(err))
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

func processInstructionTemplates(
	logger *zap.Logger,
	instruction *instructions.InstructionTemplate,
	session *chatsessions.ChatSession,
	chatHistory []chatsessions.ChatMessage,
	characterId int,
	triggerMessage *chatsessions.ChatMessage,
) bool {
	logger = logger.With(zap.Int("instructionId", instruction.ID))

	templateVars, err := newInstructionTemplateVars(session, triggerMessage, chatHistory, characterId)
	if err != nil {
		logger.Error("Error generating template variables", zap.Error(err))
		return false
	}

	fields := []*string{
		&instruction.SystemPrompt,
		&instruction.WorldSetup,
		&instruction.Instruction,
	}

	okChan := make(chan bool, len(fields))

	for _, fieldPtr := range fields {
		go func() {
			if util.HasTemplateVars(*fieldPtr) {
				tpl, err := util.NewTextTemplate("Template", *fieldPtr)
				if err != nil {
					logger.Error("Failed parsing system instruction template", zap.Error(err))
					okChan <- false
					return
				}

				*fieldPtr = strings.TrimSpace(util.WriteToString(tpl, templateVars))
			}
			okChan <- true
		}()
	}

	for i := 0; i < len(fields); i++ {
		res := <-okChan
		if !res {
			return false
		}
	}
	return true
}
