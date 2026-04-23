package processing

import (
	"context"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"juraji.nl/chat-quest/core/log"
	prov "juraji.nl/chat-quest/core/providers"
	"juraji.nl/chat-quest/core/util"
	cs "juraji.nl/chat-quest/model/chat-sessions"
	inst "juraji.nl/chat-quest/model/instructions"
	p "juraji.nl/chat-quest/model/preferences"
)

func GenerateResponseByMessageCreated(ctx context.Context, triggerMessage *cs.ChatMessage) error {
	if triggerMessage == nil || !triggerMessage.IsUser {
		// Ignore null and non-user
		return nil
	}

	sessionId := triggerMessage.ChatSessionID
	logger := log.Get().With(
		zap.String("source", "MessageCreated"),
		zap.Int("chatSessionId", sessionId),
		zap.Int("triggerMsgId", triggerMessage.ID))

	var err error

	// Cancellation
	ctx, cleanup := setupCancelBySystem(ctx, logger, "GenerateResponseByMessageCreated")
	defer cleanup()

	// Fetch Session
	session, err := cs.GetById(sessionId)
	if err != nil {
		return errors.Wrap(err, "error fetching session")
	}

	if session.PauseAutomaticResponses {
		return nil
	}

	if contextCheckPoint(ctx, logger) {
		return nil
	}

	// Select participant to respond with
	responderId, err := cs.RandomParticipantId(sessionId)
	if err != nil {
		logger.Error("Error getting random responder", zap.Error(err))
		return errors.Wrap(err, "error getting random responder")
	}
	if responderId == nil {
		logger.Warn("No participants to reply with, skipping generation")
		return nil
	}

	logger = logger.With(zap.Intp("responderId", responderId))

	return generateResponse(ctx, logger, session, triggerMessage, *responderId)
}

func GenerateResponseByParticipantTrigger(ctx context.Context, participant *cs.ChatParticipant) error {
	if participant == nil {
		// Ignore null
		return nil
	}

	sessionId := participant.ChatSessionID
	responderId := participant.CharacterID
	logger := log.Get().With(
		zap.String("source", "ParticipantTrigger"),
		zap.Int("chatSessionId", sessionId),
		zap.Int("responderId", responderId))

	// Cancellation
	ctx, cleanup := setupCancelBySystem(ctx, logger, "GenerateResponseByParticipantTrigger")
	defer cleanup()

	if contextCheckPoint(ctx, logger) {
		return nil
	}

	// Fetch Session
	session, err := cs.GetById(sessionId)
	if err != nil {
		logger.Error("Error fetching session", zap.Error(err))
		return errors.Wrap(err, "error fetching session")
	}

	return generateResponse(ctx, logger, session, nil, responderId)
}

func generateResponse(
	ctx context.Context,
	logger *zap.Logger,
	session *cs.ChatSession,
	triggerMessage *cs.ChatMessage,
	responderId int,
) error {
	if session.ChatModelId == nil {
		logger.Error("Chat model id is required on session")
		return errors.New("chat model id is required on session")
	}
	if session.ChatInstructionId == nil {
		logger.Error("Chat instruction id is required on session")
		return errors.New("chat instruction id is required on session")
	}

	// Fetch preferences
	prefs, err := p.GetPreferences(true)
	if err != nil {
		logger.Error("Error fetching preferences", zap.Error(err))
		return errors.Wrap(err, "error fetching preferences")
	}

	sessionMessageCount, err := cs.GetChatSessionMessageCount(session.ID)
	if err != nil {
		logger.Error("Error fetching chat session messages count", zap.Error(err))
		return errors.Wrap(err, "error fetching chat session messages count")
	}

	// Create instructions
	instructionVars := NewChatInstructionVars(session, prefs, triggerMessage, sessionMessageCount, responderId)
	instruction, err := inst.InstructionById(*session.ChatInstructionId)
	if err != nil {
		logger.Error("Error fetching chat instruction", zap.Error(err))
		return errors.Wrap(err, "error fetching chat instruction")
	}
	if err = instruction.ApplyTemplates(instructionVars); err != nil {
		logger.Error("Error applying instruction templates", zap.Error(err))
		return errors.Wrap(err, "error applying instruction templates")
	}

	messagesToFetch := prefs.MaxMessagesInContext
	if triggerMessage != nil {
		messagesToFetch++
	}
	includedHistory, err := cs.GetTailChatMessages(session.ID, messagesToFetch)
	if err != nil {
		logger.Error("Failed to fetch messages for context", zap.Error(err))
		return errors.Wrap(err, "failed to fetch messages for context")
	}

	// Log instruction contents
	logInstructionsToFile(logger, instruction, includedHistory)

	// Build request messages
	requestMessages := createChatRequestMessages(includedHistory, instruction)

	if contextCheckPoint(ctx, logger) {
		return nil
	}

	// Get chat model instance
	chatModelInst, err := prov.GetLlmModelInstanceById(*session.ChatModelId)
	if err != nil {
		logger.Error("Error fetching chat model instance", zap.Error(err))
		return errors.Wrap(err, "error fetching chat model instance")
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

	var characterIdPrefixInitial rune
	var characterIdPrefix string
	var characterIdSuffix string
	if instruction.CharacterIdPrefix != nil && instruction.CharacterIdSuffix != nil {
		characterIdPrefix = *instruction.CharacterIdPrefix
		characterIdSuffix = *instruction.CharacterIdSuffix
		characterIdPrefixInitial, _ = utf8.DecodeRuneInString(characterIdPrefix)
	}

	var reasoningPrefixInitial rune
	var reasoningPrefix string
	var reasoningSuffix string
	if instruction.ReasoningPrefix != nil && instruction.ReasoningSuffix != nil {
		reasoningPrefix = *instruction.ReasoningPrefix
		reasoningSuffix = *instruction.ReasoningSuffix
		reasoningPrefixInitial, _ = utf8.DecodeRuneInString(reasoningPrefix)
	}

	addMessageToStack := func() {
		newMessage := cs.NewChatMessage(false, true, &responderId, "")
		if err := cs.CreateChatMessage(session.ID, newMessage); err != nil {
			logger.Error("Failed to create response chat message", zap.Error(err))
		} else {
			currentMessage = newMessage
			messageStack = append(messageStack, currentMessage)
		}
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

	ctx, cancelCtx := context.WithCancel(ctx)
	chatResponseChan := prov.GenerateChatResponse(ctx, chatModelInst, requestMessages, instruction.AsLlmParameters())

	for {
		select {
		case response, hasNext := <-chatResponseChan:
			if response.Error != nil {
				cancelCtx()
				logger.Error("Error generating response", zap.Error(response.Error))
				return errors.Wrap(response.Error, "error generating response")
			}

			for _, token := range response.Content {
				switch currentState {
				case InContent:
					if token == characterIdPrefixInitial || token == reasoningPrefixInitial {
						currentState = PrefixDetected
						prefixBuffer.WriteRune(token)
					} else {
						// Output the token directly as it's not part of a prefix
						contentBuffer.WriteRune(token)
					}
				case PrefixDetected:
					// Accumulate prefix tokens.
					if token == '\n' {
						currentState = CancelPrefix
						continue
					}

					prefixBuffer.WriteRune(token)
					currentPrefix := prefixBuffer.String()

					// Figure out if we are in a known prefix (reasoning or Char transition)
					if len(currentPrefix) == len(reasoningPrefix) && strings.EqualFold(currentPrefix, reasoningPrefix) {
						currentState = InReasoning
						continue
					}
					if len(currentPrefix) == len(characterIdPrefix) && strings.EqualFold(currentPrefix, characterIdPrefix) {
						currentState = InCharTransition
						continue
					}

				case InReasoning:
					prefixBuffer.Reset()
					reasoningBuffer.WriteRune(token)
					currentReasoning := reasoningBuffer.String()

					if util.HasSuffixCaseInsensitive(currentReasoning, reasoningSuffix) {
						currentReasoning = strings.TrimPrefix(currentReasoning, reasoningPrefix)
						currentReasoning = strings.TrimSuffix(currentReasoning, reasoningSuffix)
						currentReasoning = strings.TrimSpace(currentReasoning)
						reasoningBuffer.Reset()
						reasoningBuffer.WriteString(currentReasoning)
						currentState = InContent
						continue
					}

				case InCharTransition:
					if token == '\n' {
						currentState = CancelPrefix
						continue
					}

					prefixBuffer.WriteRune(token)
					currentPrefix := prefixBuffer.String()

					if util.HasSuffixCaseInsensitive(currentPrefix, characterIdSuffix) {
						characterIdStr := currentPrefix[len(characterIdPrefix) : len(currentPrefix)-len(characterIdSuffix)]
						characterId, err := strconv.Atoi(characterIdStr)
						if err != nil {
							logger.Warn("Invalid character prefix found in LLM response stream, ignoring!",
								zap.String("prefix", currentPrefix))
							currentState = CancelPrefix
							continue
						}

						if charTransitionSeen {
							// We have seen a transition before, this one should be a new message on the stack.
							addMessageToStack()
						} else {
							charTransitionSeen = true
						}

						if charInSession, err := cs.CheckParticipantInSession(session.ID, characterId); !charInSession {
							logger.Warn("LLM tried to use character that is not in session, continuing as responder...",
								zap.Int("characterId", characterId),
								zap.Error(err))
							currentMessage.CharacterID = &responderId
						} else {
							currentMessage.CharacterID = &characterId
						}

						currentState = InContent
						prefixBuffer.Reset()
					}

				case CancelPrefix:
					contentBuffer.WriteString(prefixBuffer.String())
					contentBuffer.WriteRune(token)
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

			if hasContent {
				if len(currentMessage.Content) == 0 {
					currentMessage.Content = strings.TrimLeftFunc(contentBuffer.String(), unicode.IsSpace)
				} else {
					currentMessage.Content += contentBuffer.String()
				}

				contentBuffer.Reset()
			}

			if hasReasoning || hasContent {
				if err := cs.UpdateChatMessage(session.ID, currentMessage.ID, currentMessage); err != nil {
					cancelCtx()
					logger.Error("Failed to update response chat message",
						zap.Int("messageId", currentMessage.ID),
						zap.Error(err))
					return errors.Wrapf(err, "failed to update response chat message with id %d", currentMessage.ID)
				}
			}

			if response.TotalTokens != 0 || response.CompletionTokens != 0 {
				if err := cs.UpdateSessionStatistics(session.ID, response.TotalTokens, response.CompletionTokens); err != nil {
					logger.Warn("Failed to update response session statistics. (Does not break response processing!)", zap.Error(err))
				}
			}

			if !hasNext {
				cancelCtx()
				return nil
			}
		case <-ctx.Done():
			logger.Debug("Cancelled by context")
		}
	}
}
