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
	"juraji.nl/chat-quest/core/util"
	cs "juraji.nl/chat-quest/model/chat-sessions"
	inst "juraji.nl/chat-quest/model/instructions"
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
	if instruction == nil {
		return make([]p.ChatRequestMessage, 0)
	}

	// Pre-allocate with capacity for history + possible system/world messages + user instruction
	messages := make([]p.ChatRequestMessage, 0, len(chatHistory)+3)

	// Add system and world setup messages
	if instruction.SystemPrompt != nil {
		messages = append(messages, p.ChatRequestMessage{
			Role:    p.RoleSystem,
			Content: *instruction.SystemPrompt,
		})
	}
	if instruction.WorldSetup != nil {
		messages = append(messages, p.ChatRequestMessage{
			Role:    p.RoleUser,
			Content: *instruction.WorldSetup,
		})
	}

	// Add chat history
	for _, msg := range chatHistory {
		if msg.IsUser {
			messages = append(messages, p.ChatRequestMessage{
				Role:    p.RoleUser,
				Content: msg.Content,
			})
		} else {
			var msgBuffer strings.Builder

			// Add reasoning if enabled and available
			if instruction.IncludeReasoning && len(msg.Reasoning) > 0 {
				msgBuffer.WriteString(instruction.ReasoningPrefix)
				msgBuffer.WriteString(msg.Reasoning)
				msgBuffer.WriteString(instruction.ReasoningSuffix)
				msgBuffer.WriteString("\n\n")
			}

			// Add character ID
			msgBuffer.WriteString(instruction.CharacterIdPrefix)
			msgBuffer.WriteString(fmt.Sprint(*msg.CharacterID))
			msgBuffer.WriteString(instruction.CharacterIdSuffix)
			msgBuffer.WriteString("\n\n")

			// Add the main content
			msgBuffer.WriteString(msg.Content)

			messages = append(messages, p.ChatRequestMessage{
				Role:    p.RoleAssistant,
				Content: msgBuffer.String(),
			})
		}
	}

	// Add user instruction message
	messages = append(messages, p.ChatRequestMessage{
		Role:    p.RoleUser,
		Content: instruction.Instruction,
	})

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
	const (
		msgPreviewLen = 100
		timeFormat    = "2006-01-02 15:04:05"
	)
	nowStr := time.Now().Format(timeFormat)

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
		timestamp := msg.CreatedAt.Format(timeFormat)

		msgPreviewBuffer.WriteString(fmt.Sprintf("%s %s %s\n", timestamp, rolePrefix, contentPreview))
	}

	tpl := strings.TrimSpace(`
Current Time: %s

Instruction ID: %d
Name: %s
Type: %s
Temperature: %f
Max Tokens: %d
TopP: %f
Presence Penalty: %f
Frequency Penalty: %f
Stream: %v
Stop Sequences: %s
Include Reasoning: %v
Reasoning Delimiters: %s%s
Character Delimiters: %s%s

## ——— System Prompt ————————————————————————————————————————— ##
%s

## ——— World Setup ——————————————————————————————————————————— ##
%s

## ——— %d Messages (Preview) ————————————————————————————————— ##
%s

## ——— Instruction ——————————————————————————————————————————— ##
%s`)
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
		util.StrPtrOrDefault(instruction.StopSequences, "<Not Set>"),
		instruction.IncludeReasoning,
		instruction.ReasoningPrefix,
		instruction.ReasoningSuffix,
		instruction.CharacterIdPrefix,
		instruction.CharacterIdSuffix,
		util.StrPtrOrDefault(instruction.SystemPrompt, "<Not Set>"),
		util.StrPtrOrDefault(instruction.WorldSetup, "<Not Set>"),
		len(includedMessages),
		msgPreviewBuffer.String(),
		instruction.Instruction)

	// Write contents to file
	path := core.Env().MkDataDir("instructions", fmt.Sprintf("last_%s_instruction.txt", instruction.Type))
	err := os.WriteFile(path, []byte(contents), 0644)
	if err != nil {
		logger.Error("Failed to write instructions to file",
			zap.String("path", path), zap.Error(err))
		return
	}
}
