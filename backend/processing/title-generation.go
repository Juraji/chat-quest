package processing

import (
	"context"
	"strings"
	"sync"

	"github.com/pkg/errors"
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

func GenerateTitle(
	ctx context.Context,
	sessionID int,
) error {
	var err error

	// Lock while processing to avoid multiple messages invoking simultaneous generation.
	// If the lock is already active we cancel this invocation.
	lock := titleGenerationMutex.TryLock()
	if !lock {
		return errors.New("previous title generation in progress")
	}
	defer titleGenerationMutex.Unlock()

	logger := log.Get().With(
		zap.Int("sessionId", sessionID))

	// Cancellation
	ctx, cleanup := setupCancelBySystem(ctx, logger, "GenerateTitle")
	defer cleanup()

	if contextCheckPoint(ctx, logger) {
		return nil
	}
	logger.Info("Generating title for session....")

	session, err := cs.GetById(sessionID)
	if err != nil {
		logger.Error("Error getting session", zap.Error(err))
		return errors.Wrap(err, "error getting session")
	}

	prefs, err := pf.GetPreferences(true)
	if err != nil {
		logger.Error("Error getting preferences", zap.Error(err))
		return errors.Wrap(err, "error getting preferences")
	}

	if contextCheckPoint(ctx, logger) {
		return nil
	}

	modelInstance, err := p.GetLlmModelInstanceById(*prefs.TitleGenerationModelId)
	if err != nil {
		logger.Error("Could not fetch memory model", zap.Error(err))
		return errors.Wrap(err, "could not fetch memory model")
	}

	// Get message window (as per preferences)
	messageWindow, err := cs.GetTailChatMessages(sessionID, prefs.TitleGenerationMessageWindow)
	if err != nil {
		logger.Error("Error getting messages in session", zap.Error(err))
	}
	if len(messageWindow) == 0 {
		logger.Warn("No messages in session")
		return nil
	}

	sessionMessageCount, err := cs.GetChatSessionMessageCount(session.ID)
	if err != nil {
		logger.Error("Error fetching chat session messages count", zap.Error(err))
		return errors.Wrap(err, "error fetching chat session messages count")
	}

	// Build instruction
	templateVars := NewChatInstructionVars(session, prefs, nil, sessionMessageCount, 0)
	instruction, err := i.InstructionById(*prefs.TitleGenerationInstructionId)
	if err != nil {
		logger.Error("Could not fetch memory instruction", zap.Error(err))
		return errors.Wrap(err, "could not fetch memory instruction")
	}

	if err = instruction.ApplyTemplates(templateVars); err != nil {
		logger.Error("Error applying instruction templates", zap.Error(err))
		return errors.Wrap(err, "error applying instruction templates")
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
				logger.Error("Error in response",
					zap.String("generated", titleGenResponse),
					zap.Error(r.Error))
				return errors.Wrap(r.Error, "error in response")
			}

			titleGenResponse = titleGenResponse + r.Content
			if !hasNext {
				// Done
				break responseLoop
			}
		case <-ctx.Done():
			logger.Debug("Canceled by context")
			return nil
		}
	}

	session.Name = strings.Trim(titleGenResponse, "\"\n ")
	err = cs.Update(session.WorldID, sessionID, session)
	if err != nil {
		logger.Error("Error updating session", zap.Error(err))
		return errors.Wrap(err, "error updating session")
	}

	logger.Debug("Successfully generated title for session", zap.String("session", session.Name))
	return nil
}
