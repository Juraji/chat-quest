package processing

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	p "juraji.nl/chat-quest/core/providers"
	cs "juraji.nl/chat-quest/model/chat-sessions"
	inst "juraji.nl/chat-quest/model/instructions"
)

const (
	CharIdTagPrefix     = "<ByCharacterId>"
	CharIdTagPrefixInit = "<"
	CharIdTagSuffix     = "</ByCharacterId>"
)

func contextCheckPoint(ctx context.Context, logger *zap.Logger) bool {
	if ctx.Err() != nil {
		logger.Error("Cancelled by context")
		return true
	}

	return false
}

func createChatRequestMessages(
	chatHistory []cs.ChatMessage,
	instruction *inst.Instruction,
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
			content := fmt.Sprintf("%s%v%s\n\n%s", CharIdTagPrefix, *msg.CharacterID, CharIdTagSuffix, msg.Content)
			messages = append(messages, p.ChatRequestMessage{Role: p.RoleAssistant, Content: content})
		}
	}

	// Add user instruction message
	messages = append(messages, p.ChatRequestMessage{Role: p.RoleUser, Content: instruction.Instruction})
	return messages
}
