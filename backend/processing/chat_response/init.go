package chat_response

import (
	"context"
	"github.com/pkg/errors"
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
		logger.Error("Error fetching chat instruction",
			zap.Intp("instructionId", prefs.ChatInstructionID), zap.Error(err))
		return
	}

	modelInstance, err := p.GetLlmModelInstanceById(*prefs.ChatModelID)
	if err != nil {
		logger.Error("Error fetching preferred chat model instance",
			zap.Intp("modelId", prefs.ChatModelID), zap.Error(err))
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
	err = processInstructionTemplates(instruction, session, chatHistory, characterId, nil)
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
		logger.Error("Error fetching chat instruction",
			zap.Intp("instructionId", prefs.ChatInstructionID), zap.Error(err))
		return
	}

	modelInstance, err := p.GetLlmModelInstanceById(*prefs.ChatModelID)
	if err != nil {
		logger.Error("Error fetching preferred chat model instance",
			zap.Intp("modelId", prefs.ChatModelID), zap.Error(err))
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
	err = processInstructionTemplates(instruction, session, chatHistory, *characterId, message)
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
	chatResponseChan := modelInstance.GenerateChatResponse(messages, instruction.Temperature)

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
	instruction *instructions.InstructionTemplate,
	session *chatsessions.ChatSession,
	chatHistory []chatsessions.ChatMessage,
	characterId int,
	triggerMessage *chatsessions.ChatMessage,
) error {
	templateVars, err := newInstructionTemplateVars(session, triggerMessage, chatHistory, characterId)
	if err != nil {
		return err
	}

	fields := []*string{
		&instruction.SystemPrompt,
		&instruction.WorldSetup,
		&instruction.Instruction,
	}

	errChan := make(chan error, len(fields))

	for _, fieldPtr := range fields {
		go func() {
			if util.HasTemplateVars(*fieldPtr) {
				tpl, err := util.NewTextTemplate("Template", *fieldPtr)
				if err != nil {
					errChan <- errors.Wrap(err, "Error creating template for instruction template")
					return
				}

				*fieldPtr = strings.TrimSpace(util.WriteToString(tpl, templateVars))
			}
			errChan <- nil
		}()
	}

	for i := 0; i < len(fields); i++ {
		err = <-errChan
		if err != nil {
			return errors.Wrap(err, "Error processing instruction template")
		}
	}
	return nil
}
