package processing

import (
	"context"
	"strconv"
	"strings"
	"unicode"

	"go.uber.org/zap"
	"juraji.nl/chat-quest/core/log"
	prov "juraji.nl/chat-quest/core/providers"
	"juraji.nl/chat-quest/core/util"
	cs "juraji.nl/chat-quest/model/chat-sessions"
	inst "juraji.nl/chat-quest/model/instructions"
	p "juraji.nl/chat-quest/model/preferences"
)

func GenerateResponseByMessageCreated(ctx context.Context, triggerMessage *cs.ChatMessage) {
	if triggerMessage == nil || !triggerMessage.IsUser {
		// Ignore null and non-user
		return
	}

	sessionId := triggerMessage.ChatSessionID
	logger := log.Get().With(
		zap.String("source", "MessageCreated"),
		zap.Int("chatSessionId", sessionId),
		zap.Int("triggerMsgId", triggerMessage.ID))

	// Cancellation
	ctx, cleanup := setupCancelBySystem(ctx, logger, "GenerateResponse")
	defer cleanup()

	// Fetch Session
	session, err := cs.GetById(sessionId)
	if err != nil {
		logger.Error("Error fetching session", zap.Error(err))
		return
	}

	if session.PauseAutomaticResponses {
		return
	}

	if contextCheckPoint(ctx, logger) {
		return
	}

	// Select participant to respond with
	responderId, err := cs.RandomParticipantId(sessionId)
	if err != nil {
		logger.Error("Error getting random responder", zap.Error(err))
		return
	}
	if responderId == nil {
		logger.Error("No participants to reply with, skipping generation")
		return
	}

	// Fetch chat history
	chatHistory, err := cs.GetUnarchivedChatMessages(session.ID)
	if err != nil {
		logger.Error("Error fetching chat history", zap.Error(err))
		return
	}
	// Remove the last message from the history, if it is the equal to the trigger message.
	// Which it most likely is, but let's be sure
	lastMsgIndex := len(chatHistory) - 1
	if len(chatHistory) != 0 && chatHistory[lastMsgIndex].ID == triggerMessage.ID {
		chatHistory = chatHistory[:lastMsgIndex]
	}

	logger = logger.With(
		zap.Intp("responderId", responderId))

	generateResponse(ctx, logger, session, chatHistory, triggerMessage, *responderId)
}

func GenerateResponseByParticipantTrigger(ctx context.Context, participant *cs.ChatParticipant) {
	if participant == nil {
		// Ignore null
		return
	}

	sessionId := participant.ChatSessionID
	responderId := participant.CharacterID
	logger := log.Get().With(
		zap.String("source", "ParticipantTrigger"),
		zap.Int("chatSessionId", sessionId),
		zap.Int("responderId", responderId))

	// Cancellation
	ctx, cleanup := setupCancelBySystem(ctx, logger, "GenerateResponse")
	defer cleanup()

	if contextCheckPoint(ctx, logger) {
		return
	}

	// Fetch Session
	session, err := cs.GetById(sessionId)
	if err != nil {
		logger.Error("Error fetching session", zap.Error(err))
		return
	}

	// Fetch chat history
	chatHistory, err := cs.GetUnarchivedChatMessages(session.ID)
	if err != nil {
		logger.Error("Error fetching chat history", zap.Error(err))
		return
	}

	generateResponse(ctx, logger, session, chatHistory, nil, responderId)
}

func generateResponse(
	ctx context.Context,
	logger *zap.Logger,
	session *cs.ChatSession,
	chatHistory []cs.ChatMessage,
	triggerMessage *cs.ChatMessage,
	responderId int,
) {
	if session.ChatModelId == nil {
		logger.Warn("Chat model id is required on session")
		return
	}
	if session.ChatInstructionId == nil {
		logger.Warn("Chat instruction id is required on session")
		return
	}

	// Fetch preferences
	prefs, err := p.GetPreferences(true)
	if err != nil {
		logger.Error("Error getting preferences", zap.Error(err))
		return
	}

	sessionMessageCount, err := cs.GetChatSessionMessageCount(session.ID)
	if err != nil {
		logger.Error("Error fetching chat session messages count", zap.Error(err))
		return
	}

	// Create instructions
	instructionVars := NewChatInstructionVars(session, prefs, chatHistory, triggerMessage, sessionMessageCount, responderId)
	instruction, err := inst.InstructionById(*session.ChatInstructionId)
	if err != nil {
		logger.Error("Error fetching chat instruction", zap.Error(err))
		return
	}
	if err = instruction.ApplyTemplates(instructionVars); err != nil {
		logger.Error("Error applying instruction templates", zap.Error(err))
		return
	}

	includedHistory := util.SliceLastNElements(chatHistory, prefs.MaxMessagesInContext)

	// Log instruction contents
	logInstructionsToFile(logger, instruction, includedHistory)

	// Build request messages
	requestMessages := createChatRequestMessages(includedHistory, instruction)

	if contextCheckPoint(ctx, logger) {
		return
	}

	// Get chat model instance
	chatModelInst, err := prov.GetLlmModelInstanceById(*session.ChatModelId)
	if err != nil {
		logger.Error("Error fetching chat model instance", zap.Error(err))
		return
	}

	// Create message stack
	const (
		InContent = iota
		PrefixDetected
		InReasoning
		InCharTransition
		CancelPrefix
	)

	var currentState = InContent
	var charTransitionSeen = false
	var prefixBuffer strings.Builder
	var contentBuffer strings.Builder
	var reasoningBuffer strings.Builder
	var messageStack []*cs.ChatMessage
	var currentMessage *cs.ChatMessage

	addMessageToStack := func() {
		currentMessage = cs.NewChatMessage(false, true, &responderId, "")
		if err := cs.CreateChatMessage(session.ID, currentMessage); err != nil {
			logger.Error("Failed to create response chat message", zap.Error(err))
			return
		}
		messageStack = append(messageStack, currentMessage)
	}

	defer func() {
		for _, message := range messageStack {
			message.IsGenerating = false
			message.Content = strings.TrimSpace(message.Content)
			if err := cs.UpdateChatMessage(session.ID, message.ID, message); err != nil {
				logger.Error("Failed to update response chat message upon finalization",
					zap.Int("messageId", message.ID), zap.Error(err))
			}
		}
	}()

	// Create initial response message
	addMessageToStack()

	chatResponseChan := prov.GenerateChatResponse(ctx, chatModelInst, requestMessages, instruction.AsLlmParameters())

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

			for _, token := range strings.Split(response.Content, "") {
				switch currentState {
				case InContent:
					if token == PrefixInit {
						currentState = PrefixDetected
						prefixBuffer.WriteString(token)
					} else {
						// Output the token directly as it's not part of a prefix
						contentBuffer.WriteString(token)
					}
				case PrefixDetected:
					// Accumulate prefix tokens.
					if token == "\n" {
						currentState = CancelPrefix
						continue
					}

					prefixBuffer.WriteString(token)
					currentPrefix := prefixBuffer.String()

					// Figure out if we are in a known prefix (reasoning or Char transition)
					if len(currentPrefix) == len(ReasoningPrefix) && strings.EqualFold(currentPrefix, ReasoningPrefix) {
						currentState = InReasoning
						continue
					}
					if len(currentPrefix) == len(CharTransitionPrefix) && strings.EqualFold(currentPrefix, CharTransitionPrefix) {
						currentState = InCharTransition
						continue
					}

				case InReasoning:
					reasoningBuffer.WriteString(token)
					currentReasoning := reasoningBuffer.String()

					if util.HasSuffixCaseInsensitive(currentReasoning, ReasoningSuffix) {
						currentReasoning = strings.TrimPrefix(currentReasoning, ReasoningPrefix)
						currentReasoning = strings.TrimSuffix(currentReasoning, ReasoningSuffix)
						currentReasoning = strings.TrimSpace(currentReasoning)
						reasoningBuffer.Reset()
						reasoningBuffer.WriteString(currentReasoning)
						currentState = InContent
						continue
					}

				case InCharTransition:
					if token == "\n" {
						currentState = CancelPrefix
						continue
					}

					prefixBuffer.WriteString(token)
					currentPrefix := prefixBuffer.String()

					if util.HasSuffixCaseInsensitive(currentPrefix, CharTransitionSuffix) {
						characterIdStr := currentPrefix[len(CharTransitionPrefix) : len(currentPrefix)-len(CharTransitionSuffix)]
						characterId, err := strconv.Atoi(characterIdStr)
						if err != nil {
							logger.Warn("Invalid character prefix found in LLM response stream",
								zap.String("prefix", currentPrefix))
							return
						}

						if charTransitionSeen {
							// We have seen a transition before, this one should be a new message on the stack.
							addMessageToStack()
						} else {
							charTransitionSeen = true
						}

						currentMessage.CharacterID = &characterId
						currentState = InContent
						prefixBuffer.Reset()
					}

				case CancelPrefix:
					contentBuffer.WriteString(prefixBuffer.String() + token)
					prefixBuffer.Reset()
					currentState = InContent
				}
			}

			hasReasoning := reasoningBuffer.Len() > 0
			hasContent := contentBuffer.Len() > 0

			if hasReasoning {
				currentMessage.Reasoning += reasoningBuffer.String()
				reasoningBuffer.Reset()
			}

			if contentBuffer.Len() > 0 {
				if len(currentMessage.Content) == 0 {
					currentMessage.Content = strings.TrimLeftFunc(contentBuffer.String(), unicode.IsSpace)
				} else {
					currentMessage.Content += contentBuffer.String()
				}

				contentBuffer.Reset()
			}

			if hasReasoning || hasContent {
				if err := cs.UpdateChatMessage(session.ID, currentMessage.ID, currentMessage); err != nil {
					logger.Error("Failed to update response chat message",
						zap.Int("messageId", currentMessage.ID), zap.Error(err))
					return
				}
			}

		case <-ctx.Done():
			logger.Debug("Cancelled by context")
			return
		}
	}
}
