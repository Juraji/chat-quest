package chat_responses

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"juraji.nl/chat-quest/core/log"
	p "juraji.nl/chat-quest/core/providers"
	"juraji.nl/chat-quest/core/util"
	"juraji.nl/chat-quest/core/util/channels"
	c "juraji.nl/chat-quest/model/characters"
	cs "juraji.nl/chat-quest/model/chat-sessions"
	inst "juraji.nl/chat-quest/model/instructions"
	m "juraji.nl/chat-quest/model/memories"
	"juraji.nl/chat-quest/model/preferences"
	sc "juraji.nl/chat-quest/model/scenarios"
	w "juraji.nl/chat-quest/model/worlds"
	"strings"
	"sync"
)

const (
	CharIdPrefix     = "<ByCharacterId>"
	CharIdPrefixInit = "<"
	CharIdSuffix     = "</ByCharacterId>\n\n"
)

type characterTemplateVars struct {
	Character *c.Character
}

type instructionTemplateVars struct {
	MessageIndex         int
	Message              string
	IsTriggeredByMessage bool

	// Responding character info
	Character        *c.Character
	DialogueExamples []string

	// Session details
	IsSingleCharacter   bool
	OtherParticipants   []c.Character
	WorldDescription    string
	ScenarioDescription string

	// Memories
	Memories []m.Memory
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

	if contextCheckPoint(ctx, logger) {
		return
	}

	// Fetch Session
	session, err := cs.GetById(sessionId)
	if err != nil {
		logger.Error("Error fetching session", zap.Error(err))
		return
	}

	generateResponse(ctx, logger, session, responderId, nil)
}

func GenerateResponseByMessageCreated(ctx context.Context, triggerMessage *cs.ChatMessage) {
	if triggerMessage == nil || !triggerMessage.IsUser {
		// Ignore null and non-user
		return
	}

	sessionId := triggerMessage.ChatSessionID
	logger := log.Get().With(
		zap.String("source", "MessageCreated"),
		zap.Int("chatSessionId", sessionId))

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

	logger = logger.With(
		zap.Intp("responderId", responderId))

	generateResponse(ctx, logger, session, *responderId, triggerMessage)
}

func generateResponse(
	ctx context.Context,
	logger *zap.Logger,
	session *cs.ChatSession,
	responderId int,
	triggerMessage *cs.ChatMessage,
) {

	// Fetch chat history
	chatHistory, err := cs.GetUnarchivedChatMessages(session.ID)
	if err != nil {
		logger.Error("Error fetching chat history", zap.Error(err))
		return
	}
	if triggerMessage != nil && len(chatHistory) > 0 {
		chatHistory = chatHistory[:len(chatHistory)-1]
	}

	// Fetch preferences
	prefs, err := preferences.GetPreferences(true)
	if err != nil {
		logger.Error("Error getting preferences", zap.Error(err))
		return
	}

	if contextCheckPoint(ctx, logger) {
		return
	}

	// Build chat instructions
	instruction, ok := createChatInstruction(logger, session, responderId, prefs, chatHistory, triggerMessage)
	if !ok {
		logger.Error("Error creating chat instruction")
		return
	}

	// Build request messages
	requestMessages := createChatRequestMessages(chatHistory, instruction)

	if contextCheckPoint(ctx, logger) {
		return
	}

	// Get chat model instance
	chatModelInst, err := p.GetLlmModelInstanceById(*prefs.ChatModelId)
	if err != nil {
		logger.Error("Error fetching chat model instance", zap.Error(err))
		return
	}

	// Create target message
	responseMessage := cs.NewChatMessage(false, true, &responderId, "")
	if err := cs.CreateChatMessage(session.ID, responseMessage); err != nil {
		logger.Error("Failed to create response chat message", zap.Error(err))
		return
	}
	defer func() {
		responseMessage.IsGenerating = false
		if err := cs.UpdateChatMessage(session.ID, responseMessage.ID, responseMessage); err != nil {
			logger.Error("Failed to update response chat message upon finalization", zap.Error(err))
		}
	}()

	// Do LLM
	callLlmAndProcessResponse(ctx, logger, session.ID, chatModelInst, requestMessages, instruction, responseMessage)
}

func callLlmAndProcessResponse(ctx context.Context, logger *zap.Logger, sessionId int, chatModelInst *p.LlmModelInstance, requestMessages []p.ChatRequestMessage, instruction *inst.InstructionTemplate, responseMessage *cs.ChatMessage) {
	const (
		Idle = iota
		InPrefix
		InContent
	)

	var currentState = Idle
	var prefixBuffer strings.Builder
	var contentBuffer strings.Builder

	chatResponseChan := p.GenerateChatResponse(chatModelInst, requestMessages, instruction.Temperature)

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
				case Idle:
					// Check if this token starts a new prefix
					if token == CharIdPrefixInit {
						currentState = InPrefix
						prefixBuffer.WriteString(token)
					} else {
						// Output the token directly as it's not part of a prefix
						contentBuffer.WriteString(token)
					}

				case InPrefix:
					// Accumulate tokens until we complete the prefix
					prefixBuffer.WriteString(token)

					// Check if this completes a character ID prefix
					if strings.HasSuffix(prefixBuffer.String(), CharIdSuffix) {
						currentState = InContent
						prefixBuffer.Reset()
					}

				case InContent:
					// Output all tokens as they're part of the actual content
					contentBuffer.WriteString(token)
				}
			}

			if contentBuffer.Len() > 0 {
				responseMessage.Content += contentBuffer.String()
				if err := cs.UpdateChatMessage(sessionId, responseMessage.ID, responseMessage); err != nil {
					logger.Error("Failed to update response chat message", zap.Error(err))
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

func createChatRequestMessages(
	chatHistory []cs.ChatMessage,
	instruction *inst.InstructionTemplate,
) []p.ChatRequestMessage {
	// Pre-allocate messages with history len + max number of messages added here
	messages := make([]p.ChatRequestMessage, 0, len(chatHistory)+3)

	// Add system and world setup messages
	messages = append(messages,
		p.ChatRequestMessage{Role: p.RoleSystem, Content: instruction.SystemPrompt},
		p.ChatRequestMessage{Role: p.RoleUser, Content: instruction.WorldSetup},
	)

	// Add chat history
	for _, msg := range chatHistory {
		if msg.IsUser {
			messages = append(messages, p.ChatRequestMessage{Role: p.RoleUser, Content: msg.Content})
		} else {
			content := fmt.Sprintf("%s%v%s%s", CharIdPrefix, *msg.CharacterID, CharIdSuffix, msg.Content)
			messages = append(messages, p.ChatRequestMessage{Role: p.RoleAssistant, Content: content})
		}
	}

	// Add user instruction message
	messages = append(messages, p.ChatRequestMessage{Role: p.RoleUser, Content: instruction.Instruction})
	return messages
}

func createChatInstruction(
	logger *zap.Logger,
	session *cs.ChatSession,
	responderId int,
	prefs *preferences.Preferences,
	history []cs.ChatMessage,
	triggerMessage *cs.ChatMessage,
) (*inst.InstructionTemplate, bool) {
	instruction, err := inst.InstructionById(*prefs.ChatInstructionId)
	if err != nil {
		logger.Error("Error fetching chat instruction", zap.Error(err))
		return nil, false
	}

	// Asynchronously fetch stuff for template
	worldDescriptionChan := getWorldDescription(session)
	scenarioDescriptionChan := getScenarioDescription(session)
	participantsChan := getTemplatedParticipants(session, responderId)
	dialogueExamplesChan := getDialogueExamples(responderId)
	memoriesChan := getMemories(session, prefs, responderId, history, triggerMessage)

	// Unpack/handle everything
	responder, otherParticipants, err := (<-participantsChan).Unpack()
	if err != nil {
		logger.Error("Error fetching participants", zap.Error(err))
		return nil, false
	}
	worldDescription, err := (<-worldDescriptionChan).Unpack()
	if err != nil {
		logger.Error("Error fetching world description", zap.Error(err))
		return nil, false
	}
	scenarioDescription, err := (<-scenarioDescriptionChan).Unpack()
	if err != nil {
		logger.Error("Error fetching scenario description", zap.Error(err))
		return nil, false
	}
	memories, err := (<-memoriesChan).Unpack()
	if err != nil {
		logger.Error("Error fetching relevant memories", zap.Error(err))
		return nil, false
	}
	dialogueExamples, err := (<-dialogueExamplesChan).Unpack()
	if err != nil {
		logger.Error("Error fetching dialog examples", zap.Error(err))
		return nil, false
	}

	var messageContent string
	if triggerMessage != nil {
		messageContent = triggerMessage.Content
	} else {
		messageContent = ""
	}

	instructionVars := instructionTemplateVars{
		MessageIndex:         len(history),
		Message:              messageContent,
		IsTriggeredByMessage: triggerMessage != nil,
		Character:            responder,
		DialogueExamples:     dialogueExamples,
		IsSingleCharacter:    len(otherParticipants) == 0,
		OtherParticipants:    otherParticipants,
		WorldDescription:     worldDescription,
		ScenarioDescription:  scenarioDescription,
		Memories:             memories,
	}

	instruction.SystemPrompt, err = util.ParseAndApplyTextTemplate(instruction.SystemPrompt, instructionVars)
	if err != nil {
		logger.Error("Error parsing system prompt", zap.Error(err))
		return nil, false
	}
	instruction.WorldSetup, err = util.ParseAndApplyTextTemplate(instruction.WorldSetup, instructionVars)
	if err != nil {
		logger.Error("Error parsing world setup", zap.Error(err))
		return nil, false
	}
	instruction.Instruction, err = util.ParseAndApplyTextTemplate(instruction.Instruction, instructionVars)
	if err != nil {
		logger.Error("Error parsing instruction", zap.Error(err))
		return nil, false
	}

	return instruction, true
}

func getDialogueExamples(characterId int) chan *channels.Result[[]string] {
	resultChan := make(chan *channels.Result[[]string])

	go func() {
		defer close(resultChan)
		examples, err := c.DialogueExamplesByCharacterId(characterId)
		if err != nil {
			resultChan <- channels.NewErrResult[[]string](err)
			return
		}

		if len(examples) == 0 {
			// Shortcut, if there are no examples, we can skip this
			resultChan <- channels.NewResult(examples, nil)
			return
		}

		char, err := c.CharacterById(characterId)
		if err != nil {
			resultChan <- channels.NewErrResult[[]string](err)
			return
		}

		vars := &characterTemplateVars{
			Character: char,
		}

		for i, example := range examples {
			templated, err := util.ParseAndApplyTextTemplate(example, vars)
			if err != nil {
				err = errors.Wrapf(err, "Error parsing template for example '%s'", example)
				resultChan <- channels.NewErrResult[[]string](err)
				return
			}

			examples[i] = templated
		}

		resultChan <- channels.NewResult(examples, nil)
	}()

	return resultChan
}

func getWorldDescription(session *cs.ChatSession) chan *channels.Result[string] {
	resultChan := make(chan *channels.Result[string])

	go func(worldId int) {
		defer close(resultChan)
		world, err := w.WorldById(worldId)
		if err != nil {
			resultChan <- channels.NewResult[string]("", err)
			return
		}

		if world.Description != nil {
			resultChan <- channels.NewResult(*world.Description, nil)
		} else {
			resultChan <- channels.NewResult("", nil)
		}
	}(session.WorldID)

	return resultChan
}

func getScenarioDescription(session *cs.ChatSession) chan *channels.Result[string] {
	resultChan := make(chan *channels.Result[string])

	go func(scenarioId *int) {
		defer close(resultChan)
		if scenarioId == nil {
			resultChan <- channels.NewResult[string]("", nil)
			return
		}

		scenario, err := sc.ScenarioById(*scenarioId)
		if err != nil {
			resultChan <- channels.NewResult("", err)
		} else {
			resultChan <- channels.NewResult(scenario.Description, nil)
		}
	}(session.ScenarioID)

	return resultChan
}

func getMemories(
	session *cs.ChatSession,
	prefs *preferences.Preferences,
	responderId int,
	history []cs.ChatMessage,
	triggerMessage *cs.ChatMessage,
) chan *channels.Result[[]m.Memory] {
	memoriesChan := make(chan *channels.Result[[]m.Memory])

	go func() {
		defer close(memoriesChan)

		// Determine subject (what should be remember)
		var subject string
		if triggerMessage != nil {
			subject = triggerMessage.Content
		} else if len(history) > 0 {
			subject = history[len(history)-1].Content
		} else {
			// No topic message, skip memories
			memoriesChan <- channels.NewResult([]m.Memory{}, nil)
			return
		}

		memories, err := m.GetMemoriesByWorldAndCharacterIdWithEmbeddings(
			session.WorldID, responderId, *prefs.EmbeddingModelId)
		if err != nil {
			memoriesChan <- channels.NewResult([]m.Memory{}, err)
			return
		}
		if len(memories) == 0 {
			// No memories for character
			memoriesChan <- channels.NewResult([]m.Memory{}, nil)
			return
		}

		// Embed subject
		embeddingModelInst, err := p.GetLlmModelInstanceById(*prefs.EmbeddingModelId)
		if err != nil {
			memoriesChan <- channels.NewResult([]m.Memory{}, err)
			return
		}

		subjectEmbeddings, err := p.GenerateEmbeddings(embeddingModelInst, subject, true)
		if err != nil {
			memoriesChan <- channels.NewResult([]m.Memory{}, err)
			return
		}

		const batchSize = 10
		memoriesLen := len(memories)
		maxBatchCount := memoriesLen/batchSize + 1
		relevantMemoriesChan := make(chan []m.Memory, maxBatchCount)
		wg := sync.WaitGroup{}

		// Process memories in batches using goroutines
		for i := 0; i < memoriesLen; i += batchSize {
			end := i + batchSize
			if end > memoriesLen {
				end = memoriesLen
			}
			batch := memories[i:end]
			wg.Add(1)
			go func(batch []m.Memory, minP float32) {
				defer wg.Done()

				var relevantMemories []m.Memory
				for _, memory := range batch {
					if memory.AlwaysInclude {
						relevantMemories = append(relevantMemories, memory)
						continue
					}

					similarity, eErr := subjectEmbeddings.CosineSimilarity(memory.Embedding)
					if eErr != nil {
						panic(errors.Wrapf(eErr, "cosine similarity error for memory with id '%v'", memory.ID))
					}
					if similarity >= minP {
						relevantMemories = append(relevantMemories, memory)
					}
				}

				relevantMemoriesChan <- relevantMemories
			}(batch, prefs.MemoryMinP)
		}

		wg.Wait()
		close(relevantMemoriesChan)

		// Collect results from all batches
		var relevantMemories []m.Memory
		for rv := range relevantMemoriesChan {
			relevantMemories = append(relevantMemories, rv...)
		}

		memoriesChan <- channels.NewResult(relevantMemories, nil)
	}()

	return memoriesChan
}

func getTemplatedParticipants(session *cs.ChatSession, responderId int) chan *channels.PairResult[*c.Character, []c.Character] {
	resultChan := make(chan *channels.PairResult[*c.Character, []c.Character])

	go func(sessionId int) {
		defer close(resultChan)
		allParticipants, err := cs.GetParticipants(sessionId)
		if err != nil {
			resultChan <- channels.NewErrPairResult[*c.Character, []c.Character](err)
			return
		}

		otherParticipants := make([]c.Character, 0, len(allParticipants)-1)
		var responder *c.Character

		for _, participant := range allParticipants {
			err := applyCharacterTemplates(&participant)
			if err != nil {
				resultChan <- channels.NewErrPairResult[*c.Character, []c.Character](err)
				return
			}

			if participant.ID == responderId {
				responder = &participant
			} else {
				otherParticipants = append(otherParticipants, participant)
			}
		}

		if responder == nil {
			err := errors.New("Responder is not a participant")
			resultChan <- channels.NewErrPairResult[*c.Character, []c.Character](err)
			return
		}

		resultChan <- channels.NewPairResult(responder, otherParticipants, nil)
	}(session.ID)

	return resultChan
}

func applyCharacterTemplates(char *c.Character) error {
	fieldsToProcess := []*string{
		char.Appearance,
		char.Personality,
		char.History,
	}

	vars := &characterTemplateVars{
		Character: char,
	}

	for _, fieldPtr := range fieldsToProcess {
		if fieldPtr == nil {
			continue
		}

		templated, err := util.ParseAndApplyTextTemplate(*fieldPtr, vars)
		if err != nil {
			return errors.Wrap(err, "failed to apply template for character field")
		}

		*fieldPtr = templated
	}

	return nil
}

func contextCheckPoint(ctx context.Context, logger *zap.Logger) bool {
	if ctx.Err() != nil {
		logger.Error("Cancelled by context")
		return true
	}

	return false
}
