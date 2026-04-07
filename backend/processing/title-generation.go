package processing

import (
	"context"
	"strings"
	"sync"

	"go.uber.org/zap"
	"juraji.nl/chat-quest/core/log"
	p "juraji.nl/chat-quest/core/providers"
	cs "juraji.nl/chat-quest/model/chat-sessions"
	i "juraji.nl/chat-quest/model/instructions"
	pf "juraji.nl/chat-quest/model/preferences"
)

var titleResponseFormat = `{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "string",
  "minLength": 50,
  "maxLength": 100
}`

var titleGenerationMutex sync.Mutex

func GenerateTitleForSession(
	ctx context.Context,
	request *cs.ChatSessionTitleGenerateRequest,
) {
	var err error

	// Lock while processing to avoid multiple messages invoking simultaneous generation.
	// If the lock is already active we cancel this invocation.
	lock := titleGenerationMutex.TryLock()
	if !lock {
		return
	}
	defer titleGenerationMutex.Unlock()

	sessionID := request.ChatSessionID
	logger := log.Get().With(
		zap.Int("sessionId", sessionID))

	// Cancellation
	ctx, cleanup := setupCancelBySystem(ctx, logger, "GenerateTitle")
	defer cleanup()

	if contextCheckPoint(ctx, logger) {
		return
	}
	logger.Info("Generating title for session....")

	session, err := cs.GetById(sessionID)
	if err != nil {
		logger.Error("Error getting session", zap.Error(err))
		return
	}

	prefs, err := pf.GetPreferences(true)
	if err != nil {
		logger.Error("Error getting preferences", zap.Error(err))
		return
	}

	if contextCheckPoint(ctx, logger) {
		return
	}

	modelInstance, err := p.GetLlmModelInstanceById(*prefs.TitleGenerationModelId)
	if err != nil {
		logger.Error("Could not fetch memory model", zap.Error(err))
		return
	}

	// Get message window (as per preferences)
	messageWindow, err := cs.GetMessagesInSession(sessionID, prefs.TitleGenerationMessageWindow)
	if err != nil {
		logger.Error("Error getting messages in session", zap.Error(err))
	}
	if len(messageWindow) == 0 {
		logger.Warn("No messages in session")
		return
	}

	sessionMessageCount, err := cs.GetChatSessionMessageCount(session.ID)
	if err != nil {
		logger.Error("Error fetching chat session messages count", zap.Error(err))
		return
	}

	// Build instruction
	templateVars := NewChatInstructionVars(session, prefs, messageWindow, nil, sessionMessageCount, 0)
	instruction, err := i.InstructionById(*prefs.TitleGenerationInstructionId)
	if err != nil {
		logger.Error("Could not fetch memory instruction", zap.Error(err))
		return
	}

	if err = instruction.ApplyTemplates(templateVars); err != nil {
		logger.Error("Error applying instruction templates", zap.Error(err))
		return
	}

	logInstructionsToFile(logger, instruction, messageWindow)

	// Call model
	requestMessages := createChatRequestMessages(messageWindow, instruction)
	llmParameters := instruction.AsLlmParameters()
	llmParameters.ResponseFormat = &titleResponseFormat

	chatResponseChan := p.GenerateChatResponse(ctx, modelInstance, requestMessages, llmParameters)
	var titleGenResponse string

responseLoop:
	for {
		select {
		case r, hasNext := <-chatResponseChan:
			if r.Error != nil {
				logger.Error("Error generating title for session",
					zap.String("generated", titleGenResponse),
					zap.Error(r.Error))
				return
			}

			titleGenResponse = titleGenResponse + r.Content
			if !hasNext {
				// Done
				break responseLoop
			}
		case <-ctx.Done():
			logger.Debug("Canceled by context")
			return
		}
	}

	session.Name = strings.Trim(titleGenResponse, "\"")
	err = cs.Update(session.WorldID, sessionID, session)
	if err != nil {
		logger.Error("Error updating session", zap.Error(err))
		return
	}
	logger.Debug("Successfully generated title for session", zap.String("session", session.Name))
}
