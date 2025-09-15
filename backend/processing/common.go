package processing

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
	p "juraji.nl/chat-quest/core/providers"
	"juraji.nl/chat-quest/core/system"
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

// setupCancelBySystem returns a new context that will be cancelled upon emit of system.StopCurrentGeneration.
// It returns the new context and a cleanup function, to be called when the cancellation is not longer needed.
func setupCancelBySystem(ctx context.Context, logger *zap.Logger, name string) (context.Context, func()) {
	newCtx, ctxCancelFunc := context.WithCancel(ctx)
	ctxCancelKey := fmt.Sprintf("%s::%s", name, uuid.New())
	system.StopCurrentGeneration.AddListener(ctxCancelKey, func(_ context.Context, _ any) {
		logger.Info("Canceled by system", zap.String("cancelKey", ctxCancelKey))
		ctxCancelFunc()
	})
	cleanup := func() {
		system.StopCurrentGeneration.RemoveListener(ctxCancelKey)
	}

	return newCtx, cleanup
}
