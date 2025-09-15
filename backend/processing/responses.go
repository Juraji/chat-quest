package processing

import (
	"context"
	"regexp"
	"strconv"
	"strings"

	"go.uber.org/zap"
	"juraji.nl/chat-quest/core/log"
	prov "juraji.nl/chat-quest/core/providers"
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
	// Fetch preferences
	prefs, err := p.GetPreferences(true)
	if err != nil {
		logger.Error("Error getting preferences", zap.Error(err))
		return
	}

	// Create instructions
	instructionVars := NewChatInstructionVars(chatHistory, session, prefs, triggerMessage, responderId)
	instruction, err := inst.InstructionById(*prefs.ChatInstructionId)
	if err != nil {
		logger.Error("Error fetching chat instruction", zap.Error(err))
		return
	}
	if err = instruction.ApplyTemplates(instructionVars); err != nil {
		logger.Error("Error applying instruction templates", zap.Error(err))
		return
	}

	// Build request messages
	requestMessages := createChatRequestMessages(chatHistory, instruction)

	if contextCheckPoint(ctx, logger) {
		return
	}

	// Get chat model instance
	chatModelInst, err := prov.GetLlmModelInstanceById(*prefs.ChatModelId)
	if err != nil {
		logger.Error("Error fetching chat model instance", zap.Error(err))
		return
	}

	// Create message stack
	const (
		Initial = iota
		InPrefix
		InContent
	)

	var currentState = Initial
	var prefixBuffer strings.Builder
	var contentBuffer strings.Builder
	var messageStack []*cs.ChatMessage
	var currentMessage *cs.ChatMessage
	var extractCharIdRegex = regexp.MustCompile(
		regexp.QuoteMeta(CharIdTagPrefix) + "(\\d+)" + regexp.QuoteMeta(CharIdTagSuffix))

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

	chatResponseChan := prov.GenerateChatResponse(
		chatModelInst,
		requestMessages,
		instruction.AsLlmParameters(),
	)

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
				case Initial:
					if token == CharIdTagPrefixInit {
						currentState = InPrefix
						prefixBuffer.WriteString(token)
					} else {
						// Output the token directly as it's not part of a prefix
						contentBuffer.WriteString(token)
					}

				case InPrefix:
					// Accumulate tokens until we complete the prefix
					prefixBuffer.WriteString(token)
					currentPrefix := strings.TrimSpace(prefixBuffer.String())

					// When we reach the prefix length, check whether this tag sequence is actually a char id tag.
					// If not, we write the current buffer to the content and continue as InContent,
					// as this tag was not meant for us.
					if len(currentPrefix) >= len(CharIdTagPrefix) &&
						!strings.HasPrefix(currentPrefix, CharIdTagPrefix) {
						contentBuffer.WriteString(currentPrefix)
						prefixBuffer.Reset()
						currentState = Initial
					}

					// Check if this completes a character ID prefix
					if strings.HasSuffix(currentPrefix, CharIdTagSuffix) {
						// We have the complete char id tag. Extract the ID within
						// and set it as the current message's character id.
						if number := extractCharIdRegex.FindStringSubmatch(currentPrefix); number != nil && len(number) > 1 {
							charId, err := strconv.Atoi(number[1])
							if err != nil {
								logger.Warn("Invalid character prefix found in LLM response stream",
									zap.String("prefix", currentPrefix))
								return
							}
							currentMessage.CharacterID = &charId
						} else {
							logger.Warn("Invalid character prefix found in LLM response stream",
								zap.String("prefix", currentPrefix))
						}

						currentState = InContent
						prefixBuffer.Reset()
					}
				case InContent:
					// Check if this token starts a new prefix
					if token == CharIdTagPrefixInit {
						currentState = InPrefix
						prefixBuffer.WriteString(token)

						addMessageToStack()
					} else {
						// Output the token directly as it's not part of a prefix
						contentBuffer.WriteString(token)
					}
				}
			}

			if contentBuffer.Len() > 0 {
				currentMessage.Content += contentBuffer.String()
				if err := cs.UpdateChatMessage(session.ID, currentMessage.ID, currentMessage); err != nil {
					logger.Error("Failed to update response chat message",
						zap.Int("messageId", currentMessage.ID), zap.Error(err))
					return
				}

				contentBuffer.Reset()
			}

		case <-ctx.Done():
			logger.Debug("Cancelled by context")
			return
		}
	}
}
