package memory_generation

import (
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"juraji.nl/chat-quest/core/log"
	p "juraji.nl/chat-quest/core/providers"
	c "juraji.nl/chat-quest/model/characters"
	cs "juraji.nl/chat-quest/model/chat-sessions"
	i "juraji.nl/chat-quest/model/instructions"
	m "juraji.nl/chat-quest/model/memories"
	"juraji.nl/chat-quest/model/preferences"
	"strings"
	"sync"
)

type instructionTemplateVars struct {
	Participants []c.Character
}

var generationMutex sync.Mutex

func GenerateMemories(
	ctx context.Context,
	message *cs.ChatMessage,
) {
	if message.IsGenerating {
		// Skip messages that are still being generated
		return
	}

	// Lock while processing to avoid multiple messages invoking simultaneous generation
	// for the same message window.
	generationMutex.Lock()
	defer generationMutex.Unlock()

	sessionID := message.ChatSessionID
	logger := log.Get().With(zap.Int("chatSessionId", sessionID))
	if ctx.Err() != nil {
		logger.Debug("Cancelled by context")
		return
	}

	session, err := cs.GetById(sessionID)
	if err != nil {
		logger.Error("Error getting session", zap.Error(err))
		return
	}
	if !session.EnableMemories {
		// Memories are disabled for this session
		return
	}

	prefs, err := preferences.GetPreferences(true)
	if err != nil {
		logger.Error("Error getting preferences", zap.Error(err))
		return
	}

	memorizableMessages, ok := getMemorizableMessages(logger, prefs, sessionID)
	if !ok {
		return
	}

	memories, ok := generateAndExtractMemories(logger, ctx, sessionID, prefs, memorizableMessages)
	if !ok {
		return
	}
	if ctx.Err() != nil {
		logger.Debug("Cancelled by context")
		return
	}

	// We're done, save memories
	for _, memory := range memories {
		err = m.CreateMemory(session.WorldID, memory)
		if err != nil {
			logger.Error("Error creating memory", zap.Error(err))
			return
		}
	}

	// Update message processed states
	for _, chatMessage := range memorizableMessages {
		err = cs.SetMessageArchived(sessionID, chatMessage.ID)
		if err != nil {
			logger.Error("Error setting message archived bit", zap.Error(err))
		}
	}

	logger.Debug("Memory generation completed")
}

func generateAndExtractMemories(
	logger *zap.Logger,
	ctx context.Context,
	sessionID int,
	prefs *preferences.Preferences,
	messages []cs.ChatMessage,
) ([]*m.Memory, bool) {
	instruction, err := i.InstructionById(*prefs.MemoriesInstructionId)
	if err != nil {
		logger.Error("Could not fetch memory instruction", zap.Error(err))
		return nil, false
	}
	modelInstance, err := p.GetLlmModelInstanceById(*prefs.MemoriesModelId)
	if err != nil {
		logger.Error("Could not fetch memory model", zap.Error(err))
		return nil, false
	}

	participants, err := cs.GetParticipants(sessionID)
	if err != nil {
		logger.Error("Error fetching participants", zap.Error(err))
		return nil, false
	}

	templateVars := instructionTemplateVars{
		Participants: participants,
	}

	// Apply instruction template vars and generate memories
	instruction, err = i.ApplyInstructionTemplates(*instruction, templateVars)
	if err != nil {
		logger.Error("Error applying instruction templates", zap.Error(err))
		return nil, false
	}

	requestMessages := createChatRequestMessages(messages, instruction)
	memoryGenResponse, ok := callLlm(logger, ctx, modelInstance, requestMessages, instruction.Temperature)
	if !ok {
		return nil, false
	}

	// We expect a JSON markdown block, extract it
	markDownStartSeq := "```json"
	markDownEndSeq := "```"
	markdownStart := strings.Index(memoryGenResponse, markDownStartSeq)
	if markdownStart == -1 {
		logger.Error("Could not find markdown start in response")
		return nil, false
	}
	memoryGenResponse = memoryGenResponse[markdownStart+len(markDownStartSeq):]
	markdownEnd := strings.Index(memoryGenResponse, markDownEndSeq)
	if markdownEnd == -1 {
		logger.Error("Could not find markdown end in response")
		return nil, false
	}
	memoryGenResponse = memoryGenResponse[:markdownEnd]

	// Unmarshal memories
	var memories []*m.Memory
	err = json.Unmarshal([]byte(memoryGenResponse), &memories)
	if err != nil {
		logger.Error("Could not unmarshal memory response", zap.Error(err))
		return nil, false
	}

	logger.Debug("Memories generated successfully", zap.Int("memoryCount", len(memories)))
	return memories, true
}

func createChatRequestMessages(
	chatMessages []cs.ChatMessage,
	instruction *i.InstructionTemplate,
) []p.ChatRequestMessage {
	// Pre-allocate messages with history len + max number of messages added here
	messages := make([]p.ChatRequestMessage, 0, len(chatMessages)+3)

	// Add system and world setup messages
	messages = append(messages,
		p.ChatRequestMessage{Role: p.RoleSystem, Content: instruction.SystemPrompt},
		p.ChatRequestMessage{Role: p.RoleUser, Content: instruction.WorldSetup},
	)

	// Add chat history
	for _, msg := range chatMessages {
		if msg.IsUser {
			messages = append(messages, p.ChatRequestMessage{Role: p.RoleUser, Content: msg.Content})
		} else {
			content := fmt.Sprintf("<ByCharacterId>%v</ByCharacterId>\n\n%s", *msg.CharacterID, msg.Content)
			messages = append(messages, p.ChatRequestMessage{Role: p.RoleAssistant, Content: content})
		}
	}

	// Add user instruction message
	messages = append(messages, p.ChatRequestMessage{Role: p.RoleUser, Content: instruction.Instruction})
	return messages
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

func getMemorizableMessages(
	logger *zap.Logger,
	prefs *preferences.Preferences,
	sessionID int,
) ([]cs.ChatMessage, bool) {
	messages, err := cs.GetUnarchivedChatMessages(sessionID)
	if err != nil {
		logger.Error("Failed to get messages for session",
			zap.Int("sessionId", sessionID), zap.Error(err))
		return nil, false
	}

	triggerAfter := prefs.MemoryTriggerAfter
	requiredWindowSize := prefs.MemoryWindowSize

	windowSize := len(messages) - triggerAfter

	// Only proceed if we have enough messages to create a valid window
	if windowSize < requiredWindowSize {
		logger.Debug("Not enough messages to memorize",
			zap.Int("triggerAfter", triggerAfter),
			zap.Int("availableMessages", len(messages)),
			zap.Int("minWindowSize", requiredWindowSize),
			zap.Int("inWindow", windowSize))
		return nil, false
	}

	messageWindow := messages[:windowSize]

	logger.Debug("Found messages to memorize",
		zap.Int("triggerAfter", triggerAfter),
		zap.Int("availableMessages", len(messages)),
		zap.Int("minWindowSize", requiredWindowSize),
		zap.Int("inWindow", windowSize))

	return messageWindow, true
}
