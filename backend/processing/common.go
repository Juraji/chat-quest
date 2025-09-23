package processing

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"juraji.nl/chat-quest/core"
	p "juraji.nl/chat-quest/core/providers"
	"juraji.nl/chat-quest/core/system"
	cs "juraji.nl/chat-quest/model/chat-sessions"
	inst "juraji.nl/chat-quest/model/instructions"
)

const (
	PrefixInit           = "<"
	CharTransitionPrefix = "<ByCharacterId>"
	CharTransitionSuffix = "</ByCharacterId>"
	ReasoningPrefix      = "<think>"
	ReasoningSuffix      = "</think>"
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
			content := fmt.Sprintf("%s%v%s\n\n%s", CharTransitionPrefix, *msg.CharacterID, CharTransitionSuffix, msg.Content)
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

// logInstructionsToFile writes instruction details to a text file (per type) in the data directory.
// It formats the instruction information including ID, name, parameters, and content sections.
// If the log file already exists, it will be overwritten.
func logInstructionsToFile(logger *zap.Logger, instruction *inst.Instruction, includedMessages []cs.ChatMessage) {
	const msgPreviewLen = 100

	path := core.Env().MkDataDir("instructions", fmt.Sprintf("last_%s_instruction.txt", instruction.Type))

	nowStr := time.Now().Format("2006-01-02 15:04:05")

	var msgPreviewBuffer strings.Builder
	for _, msg := range includedMessages {
		contentPreview := strings.ReplaceAll(msg.Content, "\n", " ")
		if len(contentPreview) > msgPreviewLen {
			contentPreview = contentPreview[:msgPreviewLen] + "..."
		}

		rolePrefix := "User:"
		if msg.CharacterID != nil {
			rolePrefix = fmt.Sprintf("Character (%d):", *msg.CharacterID)
		}
		timestamp := msg.CreatedAt.Format("2006-01-02 15:04:05")

		msgPreviewBuffer.WriteString(fmt.Sprintf("%s %s %s\n", timestamp, rolePrefix, contentPreview))
	}

	tpl := "Current Time: %s\n\nID: %v\nName: %v\nType: %v\nTemperature: %v\nMaxTokens: %v\nTopP: %v\nPresencePenalty: %v\n" +
		"FrequencyPenalty: %v\nStream: %v\nStopSequences: %v\n\n" +
		"--- SystemPrompt ---\n%s\n\n--- WorldSetup ---\n%s\n\n--- %d Messages (Preview) ---\n%s\n--- Instruction ---\n%s\n"
	contents := fmt.Sprintf(
		tpl,
		nowStr,
		instruction.ID,
		instruction.Name,
		instruction.Type,
		instruction.Temperature,
		instruction.MaxTokens,
		instruction.TopP,
		instruction.PresencePenalty,
		instruction.FrequencyPenalty,
		instruction.Stream,
		instruction.StopSequences,
		instruction.SystemPrompt,
		instruction.WorldSetup,
		len(includedMessages),
		msgPreviewBuffer.String(),
		instruction.Instruction)

	err := os.WriteFile(path, []byte(contents), 0644)
	if err != nil {
		logger.Error("Failed to write instructions to file",
			zap.String("path", path), zap.Error(err))
		return
	}
}
